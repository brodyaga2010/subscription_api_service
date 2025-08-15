package initial

import (
	"context"
	"eff_mobile/config"
	"eff_mobile/internal/db"
	"eff_mobile/internal/handler"
	"eff_mobile/internal/repository"
	"eff_mobile/internal/service"
	pkg "eff_mobile/pkg/service"
	"log"
	"os"
)

func Init(ctx context.Context, cfg *config.Config) (context.Context, []pkg.Service, error) {
	var services []pkg.Service
	logFile, err := os.OpenFile(cfg.Logger.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening the logs file: %v", err)
	}

	db := db.New(&cfg.Database)
	services = append(services, db)

	subscriptionRepo := repository.NewSubscriptionRepository(db.Pool)

	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	api := handler.New(subscriptionService, logFile, &cfg.Server)
	api.Init()

	services = append(services, api)

	return ctx, services, nil
}
