package db

import (
	"context"
	"eff_mobile/config"
	"eff_mobile/pkg/pdb"
	"eff_mobile/pkg/service"
	"fmt"
	"log"
	"time"
)

type Database struct {
	Pool *pdb.Pool
}

var _ service.Service = (*Database)(nil)

func New(cfg *config.Database) *Database {
	p := pdb.New(cfg.ConnectionString, time.Second*time.Duration(cfg.Timeout))

	return &Database{
		Pool: p,
	}
}

func (d *Database) Run(ctx context.Context) func() error {
	return func() error {
		p, err := d.Pool.GetConnection(ctx)
		if err != nil {
			return fmt.Errorf("db pool get connection: %w", err)
		}

		if err := p.Ping(ctx); err != nil {
			return fmt.Errorf("db ping: %w", err)
		}

		log.Print("db: success ping!")

		return nil
	}
}

func (d *Database) Stop(ctx context.Context) func() error {
	return func() error {
		<-ctx.Done()

		d.Pool.Close()

		log.Print("db: connection closed!")

		return nil
	}
}
