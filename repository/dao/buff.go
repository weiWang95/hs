package dao

import "hs/repository"

var BuffRepo = NewBuffDao()

type buffDao struct {
}

func NewBuffDao() repository.BuffRepo {
	return &buffDao{}
}
