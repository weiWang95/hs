package repository

import (
	"context"

	"hs/repository/entity"
)

type ConnectRepo interface {
	Save(ctx context.Context, conn *entity.Connect) error
	Remove(ctx context.Context, userId uint64) error
	Find(ctx context.Context, userId uint64) (*entity.Connect, error)
	SetGameId(ctx context.Context, userId uint64, gameId uint64) error
}
