package game

import (
	"context"
	"hs/pkg/list"
	"hs/pkg/utils"
	"hs/repository/dao/buff"
	"hs/repository/entity"

	"github.com/sirupsen/logrus"
)

type BuffSystem struct {
	runner *GameRunner
}

func NewBuffSystem(runner *GameRunner) *BuffSystem {
	return &BuffSystem{runner: runner}
}

func (s *BuffSystem) OnPlace(ctx context.Context, trigger *entity.Retinue, list *list.DoubleLinkList[entity.Retinue], placeIdx int, targetList *list.DoubleLinkList[entity.Retinue], targetIdx int) {
	s.OnEvent(ctx, entity.BuffOnPlace, trigger, list, targetList, placeIdx, targetIdx, 0)
}

func (s *BuffSystem) OnDead(ctx context.Context, trigger *entity.Retinue, list *list.DoubleLinkList[entity.Retinue], idx int, enemyList *list.DoubleLinkList[entity.Retinue], attackIdx int) {
	s.OnEvent(ctx, entity.BuffOnDead, trigger, list, enemyList, idx, attackIdx, 0)
}

func (s *BuffSystem) OnEvent(ctx context.Context, eventType entity.BuffEvent, trigger *entity.Retinue, triggerList, targetList *list.DoubleLinkList[entity.Retinue], triggerIdx, targetIdx int, stackDeep int) {
	if stackDeep > 10 {
		return
	}
	stackDeep++

	switch eventType {
	case entity.BuffOnPlace:
		trigger.EachBuff(entity.BuffOnPlace, s.onTriggerEvent(ctx, trigger, trigger, triggerList, targetList, triggerIdx, targetIdx, stackDeep, true))
	case entity.BuffOnBeforeAttack, entity.BuffOnAfterAttack, entity.BuffOnHurt:
		trigger.EachBuff(eventType, s.onTriggerEvent(ctx, trigger, trigger, triggerList, targetList, triggerIdx, targetIdx, stackDeep, true))
	case entity.BuffOnDead:
		trigger.EachBuff(entity.BuffOnDead, s.onTriggerEvent(ctx, trigger, trigger, triggerList, targetList, triggerIdx, targetIdx, stackDeep, true))
	case entity.BuffOnBuffAdded:
		trigger.EachBuff(entity.BuffOnBuffAdded, s.onTriggerEvent(ctx, trigger, trigger, triggerList, targetList, triggerIdx, targetIdx, stackDeep, false))
	case entity.BuffOnBuyCard:
	case entity.BuffOnSellCard:
	case entity.BuffOnFightStart:
	case entity.BuffOnFightEnd:
	case entity.BuffOnRoundStart:
	case entity.BuffOnRoundEnd:
	}

	triggerList.EachWithIdx(func(idx int, r *entity.Retinue) bool {
		return r.EachBuff(eventType, s.onHaloTrigger(ctx, r, trigger, triggerList, targetList, idx, targetIdx, stackDeep, true))
	})
}

func (s *BuffSystem) onTriggerEvent(ctx context.Context, self, trigger *entity.Retinue, triggerList, targetList *list.DoubleLinkList[entity.Retinue], triggerIdx, targetIdx int, stackDeep int, allowSelf bool) func(b entity.IBuff) bool {
	return func(b entity.IBuff) bool {
		logrus.Debugf("trigger event:%d type:%d => %d -> %d", b.Base().Event, b.Base().Type, triggerIdx, targetIdx)

		switch b.Base().Type {
		case entity.BuffTypeAdd:
			s.onBuffAdded(ctx, b, trigger, triggerList, targetList, triggerIdx, targetIdx, stackDeep, allowSelf)
		case entity.BuffTypeAddBuff:
		case entity.BuffTypeDamage:
			s.onDamageRecv(ctx, b, trigger, triggerList, targetList, triggerIdx, targetIdx, stackDeep)
		case entity.BuffTypeSummon:
			s.onSummon(ctx, b, trigger, triggerList, triggerIdx)
		}

		return true
	}
}

func (s *BuffSystem) onHaloTrigger(ctx context.Context, self, trigger *entity.Retinue, triggerList, targetList *list.DoubleLinkList[entity.Retinue], triggerIdx, targetIdx int, stackDeep int, allowSelf bool) func(b entity.IBuff) bool {
	return func(b entity.IBuff) bool {
		if b.Base().Type != entity.BuffTypeHalo {
			return true
		}

		return true
	}
}

func (s *BuffSystem) onBuffAdded(ctx context.Context, b entity.IBuff, trigger *entity.Retinue, triggerList, targetList *list.DoubleLinkList[entity.Retinue], triggerIdx, targetIdx int, stackDeep int, allowSelf bool) {
	targets := s.getBuffTargets(ctx, b, triggerList, targetList, triggerIdx, targetIdx)
	for _, t := range targets {
		target := t.list.Get(t.idx)
		if target == nil {
			continue
		}
		logrus.Debugf("AddBuff %v to %v", b.Base().Name, target.Inspect())
		target.AddBuff(b)
		if b.Base().Type == entity.BuffTypeAdd && (allowSelf || (t.list == triggerList && t.idx != triggerIdx)) {
			s.OnEvent(ctx, entity.BuffOnBuffAdded, target, t.list, nil, t.idx, 0, stackDeep)
		}
	}
}

func (s *BuffSystem) onDamageRecv(ctx context.Context, b entity.IBuff, trigger *entity.Retinue, triggerList, targetList *list.DoubleLinkList[entity.Retinue], triggerIdx, targetIdx int, stackDeep int) {
	targets := s.getBuffTargets(ctx, b, triggerList, targetList, triggerIdx, targetIdx)
	for _, t := range targets {
		target := t.list.Get(t.idx)
		if target == nil {
			continue
		}

		v, ok := b.(*buff.DamageBuff)
		if !ok {
			continue
		}

		target.FinalHp -= v.Damage
		s.OnEvent(ctx, entity.BuffOnHurt, target, targetList, nil, targetIdx, 0, stackDeep)
	}
}

func (s *BuffSystem) onSummon(ctx context.Context, b entity.IBuff, trigger *entity.Retinue, triggerList *list.DoubleLinkList[entity.Retinue], triggerIdx int) {
	v, ok := b.(*buff.SummonBuff)
	if !ok {
		return
	}

	r := s.runner.cardPool.GetCard(ctx, v.CardId)
	if r == nil {
		return
	}

	triggerList.AddAt(triggerIdx+1, r)
}

type BuffTarget struct {
	list *list.DoubleLinkList[entity.Retinue]
	idx  int
}

func (s *BuffSystem) getBuffTargets(ctx context.Context, b entity.IBuff, triggerList, targetList *list.DoubleLinkList[entity.Retinue], triggerIdx, targetIdx int) []BuffTarget {
	switch b.Base().Target {
	case entity.TargetSelf:
		return []BuffTarget{{list: triggerList, idx: triggerIdx}}
	case entity.TargetNearby:
		return []BuffTarget{{list: triggerList, idx: triggerIdx - 1}, {list: triggerList, idx: triggerIdx + 1}}
	case entity.TargetTrigger:
		return []BuffTarget{{list: triggerList, idx: triggerIdx}}
	case entity.TargetAppoint:
		return []BuffTarget{{list: targetList, idx: targetIdx}}
	case entity.TargetAttacked:
		return []BuffTarget{{list: targetList, idx: targetIdx}}
	case entity.TargetRandomAmicable:
		return []BuffTarget{{list: triggerList, idx: utils.RandIntWithOut(triggerList.Size(), triggerIdx)}}
	case entity.TargetRandomEnemy:
		return []BuffTarget{{list: targetList, idx: utils.RandIntWithOut(targetList.Size(), -1)}}
	case entity.TargetAllAmicable:
		return s.allTargets(triggerList)
	case entity.TargetAllEnemy:
		return s.allTargets(targetList)
	case entity.TargetAll:
		return append(s.allTargets(triggerList), s.allTargets(targetList)...)
	}
	return []BuffTarget{}
}

func (s *BuffSystem) allTargets(list *list.DoubleLinkList[entity.Retinue]) []BuffTarget {
	res := make([]BuffTarget, 0, list.Size())
	for i := range list.Size() {
		res = append(res, BuffTarget{list: list, idx: i})
	}
	return res
}
