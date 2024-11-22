package buff

import "hs/repository/entity"

type HaloBuff struct {
	*entity.BuffBase
}

func (b *HaloBuff) Base() *entity.BuffBase {
	return b.BuffBase
}

type AdditionHaloBuff struct {
	*HaloBuff

	AddAttack int32 `json:"attack,omitempty"`
	AddHp     int32 `json:"hp,omitempty"`
}
