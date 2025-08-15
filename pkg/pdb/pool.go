package pdb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	conn     *pgxpool.Pool
	connStr  string
	duration time.Duration
}

func New(connStr string, duration time.Duration) *Pool {
	return &Pool{
		connStr:  connStr,
		duration: duration,
	}
}

func (p *Pool) GetConnection(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	if p.conn == nil {
		p.conn, err = createPool(ctx, p.connStr, p.duration)
		if err != nil {
			return nil, err
		}
	}

	if err := p.conn.Ping(ctx); err != nil {
		p.conn, err = createPool(ctx, p.connStr, p.duration)
		if err == nil {
			return p.conn, nil
		}

		ticker := time.NewTicker(p.duration)
		defer ticker.Stop()
		for range ticker.C {
			p.conn, err = createPool(ctx, p.connStr, p.duration)
			if err == nil {
				break
			}
		}
		if p.conn == nil {
			return nil, fmt.Errorf("failed to create db pool after retries: %w", err)
		}
	}

	return p.conn, nil
}

func (p *Pool) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
