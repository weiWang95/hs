package dao

import (
	"context"
	"hs/pkg/errors"
	"hs/repository"
	"hs/repository/entity"
	"sync"
)

type gameDao struct {
	m sync.Map
}

func NewGameRepo() repository.GameRepo {
	return &gameDao{}
}

func (d *gameDao) SaveGame(ctx context.Context, game *entity.Game) error {
	d.m.Store(game.Id, game)
	return nil
}

func (d *gameDao) GetGame(ctx context.Context, id uint64) (*entity.Game, error) {
	v, ok := d.m.Load(id)
	if !ok {
		return nil, errors.NewError(errors.CodeNotFound, "game not found")
	}

	return v.(*entity.Game), nil
}

func (d *gameDao) DelGame(ctx context.Context, id uint64) error {
	d.m.Delete(id)
	return nil
}
