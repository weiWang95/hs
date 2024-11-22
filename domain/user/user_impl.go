package user

import (
	"context"
	"hs/domain"
	"hs/repository"
	"hs/repository/entity"
)

type userImpl struct {
	connRepo repository.ConnectRepo
}

func NewUser() domain.IUser {
	return &userImpl{}
}

func (u *userImpl) Login(ctx context.Context, in domain.LoginParam) (*entity.User, error) {
	return nil, nil
}
