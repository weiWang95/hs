package repository

import "context"

type MatcherRepo interface {
	GetCount(ctx context.Context) (int64, error)
	Add(ctx context.Context, userId uint64) error
	Remove(ctx context.Context, userId uint64) error
	BulkPop(ctx context.Context, limit int) ([]uint64, error)
}
