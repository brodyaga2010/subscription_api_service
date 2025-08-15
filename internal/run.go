package internal

import (
	"context"
	"eff_mobile/config"
	"eff_mobile/internal/initial"
	"eff_mobile/pkg/service"
	"fmt"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	ctx, services, err := initial.Init(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error initializing services: %v", err)
	}

	if err := service.Run(ctx, services); err != nil {
		return fmt.Errorf("error running services: %v", err)
	}

	return nil
}
