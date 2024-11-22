package repository

import (
	"context"
	"hs/repository/entity"
)

type CardRepo interface {
	QueryCards(ctx context.Context, param QueryCardsParam) ([]entity.Retinue, error)
}

type QueryCardsParam struct {
	Level int
}
