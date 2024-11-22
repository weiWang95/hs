package domain

import (
	"context"
)

type IMatcher interface {
	Start(ctx context.Context)
	Join(ctx context.Context, userId uint64) error
	Quit(ctx context.Context, userId uint64) error
}
