package service

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, s []Service) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	for _, svc := range s {
		g.Go(svc.Run(ctx))
		g.Go(svc.Stop(ctx))
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("service run error: %w", err)
	}

	return nil
}
