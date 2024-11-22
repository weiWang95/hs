package buff

import (
	"hs/repository/entity"
)

type AdditionBuff struct {
	*entity.BuffBase
	AddAttack int32 `json:"attack,omitempty"`
	AddHp     int32 `json:"hp,omitempty"`
}

func (b *AdditionBuff) Base() *entity.BuffBase {
	return b.BuffBase
}

func (b *AdditionBuff) CaluateAttr(r *entity.Retinue) {
	r.FinalAttack += b.AddAttack
	r.FinalHp += b.AddHp
}
