package dao

import (
	"context"
	"hs/pkg/config"
	"hs/repository"
	"hs/repository/entity"

	"github.com/sirupsen/logrus"
)

var CardRepo = NewCardRepo()

type cardDao struct {
	m map[uint64]entity.Retinue
}

func NewCardRepo() repository.CardRepo {
	d := cardDao{m: make(map[uint64]entity.Retinue)}
	d.init()
	return &d
}

func (d *cardDao) init() {
	var data []entity.Retinue
	if err := config.LoadJsonConfig("card.json", &data); err != nil {
		logrus.Errorf("load cards data error: %v", err)
		return
	}

	for _, v := range data {
		v.Refresh()
		logrus.Debugf("load card: %v", v.Inspect())
		d.m[v.Id] = v
	}
}

func (d *cardDao) QueryCards(ctx context.Context, param repository.QueryCardsParam) ([]entity.Retinue, error) {
	r := make([]entity.Retinue, 0)
	for _, v := range d.m {
		if param.Level == 0 || v.Level == param.Level {
			r = append(r, v)
		}
	}

	return r, nil
}
