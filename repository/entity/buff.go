package entity

import "hs/pkg/list"

func init() {
	list.RegisterValue(Buff{})
}

type BuffEvent uint8

const (
	BuffEventNull BuffEvent = iota
	BuffOnPlace
	BuffOnSummon
	BuffOnBeforeAttack
	BuffOnAfterAttack
	BuffOnAttacked
	BuffOnHurt
	BuffOnDead
	BuffOnBuffAdded
	BuffOnBuyCard
	BuffOnSellCard
	BuffOnFightStart
	BuffOnFightEnd
	BuffOnRoundStart
	BuffOnRoundEnd
)

type BuffType uint8

const (
	BuffTypeAdd BuffType = iota + 1
	BuffTypeDamage
	BuffTypeAddBuff
	BuffTypeHalo
	BuffTypeSummon
)

type BuffTarget uint8

const (
	TargetSelf BuffTarget = iota + 1
	TargetNearby
	TargetAppoint
	TargetTrigger
	TargetAttacked
	TargetRandomAmicable
	TargetRandomEnemy
	TargetAllAmicable
	TargetAllEnemy
	TargetAll
)

type IBuff interface {
	Base() *BuffBase
	CaluateAttr(r *Retinue)
}

type BuffBase struct {
	Id     uint64     `json:"id,omitempty"`
	Name   string     `json:"name,omitempty"`
	Desc   string     `json:"desc,omitempty"`
	Type   BuffType   `json:"type,omitempty"`
	Event  BuffEvent  `json:"event,omitempty"`
	Target BuffTarget `json:"target,omitempty"`
}

type Buff struct {
	*BuffBase

	AddAttack int32  `json:"add_attack,omitempty"`
	AddHp     int32  `json:"add_hp,omitempty"`
	Damage    int32  `json:"damage,omitempty"`
	CardId    uint64 `json:"card_id,omitempty"`
	BuffId    uint64 `json:"buff_id,omitempty"`
}

func (b *Buff) Base() *BuffBase {
	return b.BuffBase
}
