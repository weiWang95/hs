package repository

import (
	"context"

	"hs/repository/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context) (*entity.User, error)
	GetUser(ctx context.Context, id uint64) (*entity.User, error)
	BulkGetUsers(ctx context.Context, ids []uint64) ([]entity.User, error)
	SaveUser(ctx context.Context, user *entity.User) error
}
