package dao

import (
	"context"
	"sync"
	"sync/atomic"

	"hs/pkg/utils"
	"hs/repository"
	"hs/repository/entity"
)

var UserRepo = NewUserRepo()

func NewUserRepo() repository.UserRepo {
	return &userDao{}
}

type userDao struct {
	idSeq uint64
	m     sync.Map
}

func (d *userDao) CreateUser(ctx context.Context) (*entity.User, error) {
	user := &entity.User{
		Id:       atomic.AddUint64(&d.idSeq, 1),
		Nickname: utils.RandomString(8),
	}

	if err := d.SaveUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (d *userDao) GetUser(ctx context.Context, id uint64) (*entity.User, error) {
	v, ok := d.m.Load(id)
	if !ok {
		return nil, nil
	}
	return v.(*entity.User), nil
}

func (d *userDao) BulkGetUsers(ctx context.Context, ids []uint64) ([]entity.User, error) {
	r := make([]entity.User, 0)
	for _, id := range ids {
		v, ok := d.m.Load(id)
		if !ok {
			continue
		}
		r = append(r, *v.(*entity.User))
	}
	return r, nil
}

func (d *userDao) SaveUser(ctx context.Context, user *entity.User) error {
	d.m.Store(user.Id, user)
	return nil
}
