package entity

import (
	"hs/pkg/list"
)

type Shop struct {
	Level       int                           `json:"level"`
	UpgradeCost int                           `json:"upgrade_cost"`
	RetinueCap  int                           `json:"retinue_cap"`
	Retinue     *list.DoubleLinkList[Retinue] `json:"retinue"`
}

type ShopConfig struct {
	Levels []ShopLevelConfig `json:"levels"`
}

type ShopLevelConfig struct {
	Level         int   `json:"level"`
	UpgradeCost   int   `json:"upgrade_cost"`
	RetinueCap    int   `json:"retinue_cap"`
	Probabilities []int `json:"probabilities"`
}
