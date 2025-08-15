package handler

import (
	"context"
	"eff_mobile/config"
	"eff_mobile/internal/service"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "eff_mobile/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	l "github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type SubscriptionApi struct {
	srvc *service.SubscriptionService
	log  echo.Logger
	cfg  *config.Server
	e    *echo.Echo
	file *os.File
}

func New(srvc *service.SubscriptionService, logFile *os.File, cfg *config.Server) *SubscriptionApi {
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	e := echo.New()
	e.Logger.SetOutput(multiWriter)
	log := e.Logger

	return &SubscriptionApi{
		srvc: srvc,
		log:  log,
		cfg:  cfg,
		e:    e,
		file: logFile,
	}
}

func (api *SubscriptionApi) Init() {
	api.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","level":"INFO","remote_ip":"${remote_ip}","host":"${host}","method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}","bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))
	api.e.Use(middleware.Recover())
	api.e.Logger.SetLevel(l.INFO)

	api.e.POST("/subscriptions", api.Create)
	api.e.GET("/subscriptions/:id", api.Get)
	api.e.PUT("/subscriptions/:id", api.Update)
	api.e.DELETE("/subscriptions/:id", api.Delete)
	api.e.GET("/subscriptions", api.List)
	api.e.GET("/subscriptions/sum", api.CalculateAmount)
	api.e.GET("/docs/*", echoSwagger.WrapHandler)
}

func (api *SubscriptionApi) Run(ctx context.Context) func() error {
	return func() error {
		go func() {
			if err := api.e.Start(api.cfg.Host + ":" + api.cfg.Port); err != nil && err != http.ErrServerClosed {
				api.log.Error("Failed to start server:", err)
			}
		}()

		api.log.Info("Server started on ", api.cfg.Host+":"+api.cfg.Port)
		return nil
	}
}

func (api *SubscriptionApi) Stop(ctx context.Context) func() error {
	return func() error {
		<-ctx.Done()
		defer api.file.Close()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := api.e.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		api.log.Info("Server stopped")
		return nil
	}
}
