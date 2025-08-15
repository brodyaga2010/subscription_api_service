package service

import (
	"context"
)

type Service interface {
	Run(ctx context.Context) func() error
	Stop(ctx context.Context) func() error
}
