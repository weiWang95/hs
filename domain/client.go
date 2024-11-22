package domain

import "context"

type IClient interface {
	Run(ctx context.Context) error
	Close() error
}
