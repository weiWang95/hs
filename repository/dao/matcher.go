package dao

import (
	"context"

	"hs/pkg/list"
	"hs/repository"
)

var MatcherDao = NewMatcherRepo()

type matcherDao struct {
	list *list.List[uint64]
}

func NewMatcherRepo() repository.MatcherRepo {
	return &matcherDao{list: list.NewList[uint64]()}
}

func (d *matcherDao) GetCount(ctx context.Context) (int64, error) {
	return int64(d.list.Size()), nil
}

func (d *matcherDao) Add(ctx context.Context, userId uint64) error {
	d.list.Add(userId)
	return nil
}

func (d *matcherDao) Remove(ctx context.Context, userId uint64) error {
	d.list.Remove(userId)
	return nil
}

func (d *matcherDao) BulkPop(ctx context.Context, limit int) ([]uint64, error) {
	r := make([]uint64, 0, limit)
	for i := 0; i < limit; i++ {
		r = append(r, d.list.RemoveAt(0))
	}
	return r, nil
}
