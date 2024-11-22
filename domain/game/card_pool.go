package game

import (
	"context"
	"math/rand/v2"

	"hs/pkg/utils"
	"hs/repository"
	"hs/repository/dao"
	"hs/repository/entity"
)

type CardPool struct {
	m map[int][]entity.Retinue
}

func NewCardPool() *CardPool {
	return &CardPool{
		m: make(map[int][]entity.Retinue, 6),
	}
}

func (p *CardPool) Init(ctx context.Context) error {
	param := repository.QueryCardsParam{Level: 1}

	serverCfg := dao.ConfigRepo.GetServerConfig()

	for param.Level <= serverCfg.MaxShopLevel {
		r, err := dao.CardRepo.QueryCards(ctx, param)
		if err != nil {
			return err
		}
		p.m[param.Level] = r
		param.Level += 1
	}

	return nil
}

func (p *CardPool) Random(ctx context.Context, num int, probabilities []int) []entity.Retinue {
	r := make([]entity.Retinue, 0, num)

	for i := 0; i < num; i++ {
		level := utils.Random(probabilities) + 1
		idx := rand.IntN(len(p.m[level]))
		card := p.m[level][idx]
		card.Init()
		r = append(r, card)
	}

	return r
}

func (p *CardPool) GetCard(ctx context.Context, id uint64) *entity.Retinue {
	for _, arr := range p.m {
		for _, v := range arr {
			if v.Id == id {
				card := v
				card.Init()
				return &card
			}
		}
	}
	return nil
}
