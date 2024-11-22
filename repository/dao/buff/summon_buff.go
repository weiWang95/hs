package buff

import "hs/repository/entity"

type SummonBuff struct {
	*entity.BuffBase

	CardId uint64 `json:"card_id,omitempty"`
}

func (b *SummonBuff) Base() *entity.BuffBase {
	return b.BuffBase
}

func (b *SummonBuff) CaluateAttr(r *entity.Retinue) {
}
