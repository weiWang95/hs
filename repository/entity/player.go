package entity

import (
	"fmt"

	"hs/pkg/list"
)

type PlayerState int

const (
	PlayerAlive PlayerState = iota
	PlayerDead
)

type Player struct {
	Id       uint64      `json:"id"`
	Nickname string      `json:"nickname"`
	State    PlayerState `json:"state"`

	MaxHp  int32 `json:"max_hp"`
	Hp     int32 `json:"hp"`
	Shield int32 `json:"shield"`

	MaxGold int32 `json:"max_gold"`
	Gold    int32 `json:"gold"`

	CardList    *list.DoubleLinkList[Retinue] `json:"card_list"`
	RetinueList *list.DoubleLinkList[Retinue] `json:"retinue_list"`
}

func (p Player) Inspect() string {
	return fmt.Sprintf("{[%d] %d (%d/%d)}", p.Id, p.Hp, p.Gold, p.MaxGold)
}
