package domain

import (
	"context"
	"hs/repository/entity"
)

type IUser interface {
	Login(ctx context.Context, in LoginParam) (*entity.User, error)
}

type LoginParam struct {
	Username string
	Password string
}
