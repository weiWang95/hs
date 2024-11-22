package buff

import "hs/repository/entity"

type DamageBuff struct {
	*entity.BuffBase

	Damage int32 `json:"damage,omitempty"`
}

func (b *DamageBuff) Base() *entity.BuffBase {
	return b.BuffBase
}

func (b *DamageBuff) CaluateAttr(r *entity.Retinue) {
}
