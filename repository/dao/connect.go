package dao

import (
	"context"
	"sync"

	"hs/repository"
	"hs/repository/entity"
)

var ConnectRepo = NewConnectDao()

type ConnectDao struct {
	m sync.Map
}

func NewConnectDao() repository.ConnectRepo {
	return &ConnectDao{}
}

func (d *ConnectDao) Save(ctx context.Context, conn *entity.Connect) error {
	d.m.Store(conn.UserId, conn)
	return nil
}

func (d *ConnectDao) Remove(ctx context.Context, userId uint64) error {
	d.m.Delete(userId)
	return nil
}

func (d *ConnectDao) Find(ctx context.Context, userId uint64) (*entity.Connect, error) {
	v, ok := d.m.Load(userId)
	if !ok {
		return nil, nil
	}
	return v.(*entity.Connect), nil
}

func (d *ConnectDao) SetGameId(ctx context.Context, userId uint64, gameId uint64) error {
	v, ok := d.m.Load(userId)
	if !ok {
		return nil
	}
	v.(*entity.Connect).GameId = gameId
	return nil
}
