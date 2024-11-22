package repository

import (
	"context"

	"hs/repository/entity"
)

type GameRepo interface {
	SaveGame(ctx context.Context, game *entity.Game) error
	GetGame(ctx context.Context, id uint64) (*entity.Game, error)
	DelGame(ctx context.Context, id uint64) error
}
