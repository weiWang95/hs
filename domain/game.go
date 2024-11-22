package domain

import (
	"context"
)

type IGame interface {
	Start(ctx context.Context, userIds []uint64) (uint64, error)
	GetRunner(ctx context.Context, id uint64) (IGameRunner, error)

	GetGameTotal(ctx context.Context) int64
	IsFull(ctx context.Context) bool
}
