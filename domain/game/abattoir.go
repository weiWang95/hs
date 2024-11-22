package game

import (
	"context"
	"hs/pkg/list"
	"hs/repository/entity"
	"math/rand/v2"

	"github.com/sirupsen/logrus"
)

type Abattoir struct {
	runner *GameRunner

	attacker uint64
	recv     uint64

	attackerIdx int
	recvIdx     int

	attackList *list.DoubleLinkList[entity.Retinue]
	recvList   *list.DoubleLinkList[entity.Retinue]
	buffSys    *BuffSystem
}

func NewAbattoir(runner *GameRunner, playerA, playerB *entity.Player) *Abattoir {
	a := OM.NewAbattoir()
	a.runner = runner

	a.attacker = playerA.Id
	a.recv = playerB.Id
	a.attackerIdx = 0
	a.recvIdx = 0
	a.attackList = a.cloneRetinueList(playerA.RetinueList)
	a.recvList = a.cloneRetinueList(playerB.RetinueList)
	a.buffSys = NewBuffSystem(a.runner)

	return a
}

func (a *Abattoir) Fight(ctx context.Context) map[uint64]*list.DoubleLinkList[entity.Retinue] {
	a.computeFirstHand()

	for a.attackList.Size() != 0 && a.recvList.Size() != 0 {
		logrus.Tracef("Start: attacker: %v, recv: %v", a.attacker, a.recv)
		// 攻击者
		attacker := a.attackList.Get(a.attackerIdx)

		// 随机目标
		targetIdx := rand.IntN(a.recvList.Size())
		recv := a.recvList.Get(targetIdx)

		a.buffSys.OnEvent(ctx, entity.BuffOnBeforeAttack, attacker, a.attackList, a.recvList, a.attackerIdx, targetIdx, 0)
		logrus.Tracef("Before: attacker: %v, recv: %v", attacker.Inspect(), recv.Inspect())

		a.Attack(attacker, recv)

		a.buffSys.OnEvent(ctx, entity.BuffOnAfterAttack, attacker, a.attackList, a.recvList, a.attackerIdx, targetIdx, 0)
		logrus.Tracef("After: attacker: %v, recv: %v", attacker.Inspect(), recv.Inspect())

		a.buffSys.OnEvent(ctx, entity.BuffOnHurt, attacker, a.attackList, a.recvList, a.attackerIdx, targetIdx, 0)
		a.buffSys.OnEvent(ctx, entity.BuffOnHurt, recv, a.recvList, a.attackList, targetIdx, a.attackerIdx, 0)

		if attacker.FinalHp <= 0 {
			OM.PutRetinue(a.attackList.Del(a.attackerIdx))
			a.buffSys.OnEvent(ctx, entity.BuffOnDead, attacker, a.attackList, a.recvList, a.attackerIdx, targetIdx, 0)
		}
		if recv.FinalHp <= 0 {
			OM.PutRetinue(a.recvList.Del(targetIdx))
			a.buffSys.OnEvent(ctx, entity.BuffOnDead, recv, a.recvList, a.attackList, targetIdx, a.attackerIdx, 0)
			if a.recvIdx > targetIdx {
				a.recvIdx--
			}
		}

		logrus.Tracef("End: attacker: %v, recv: %v", a.attacker, a.recv)

		logrus.Trace("-- Switch Attacker --")
		a.switchAttacker()
	}

	r := make(map[uint64]*list.DoubleLinkList[entity.Retinue], 2)
	r[a.attacker] = a.attackList
	r[a.recv] = a.recvList

	return r
}

func (a *Abattoir) computeFirstHand() {
	if a.attackList.Size() > a.recvList.Size() {
		return
	}

	if a.attackList.Size() < a.recvList.Size() || rand.IntN(2) == 1 {
		a.switchAttacker()
	}
}

func (a *Abattoir) Attack(r1, r2 *entity.Retinue) {
	logrus.Tracef("attacker: %v, recv: %v", r1.Inspect(), r2.Inspect())
	r1.FinalHp -= r2.FinalAttack
	r2.FinalHp -= r1.FinalAttack
}

func (a *Abattoir) switchAttacker() {
	a.attacker, a.recv = a.recv, a.attacker
	a.attackerIdx, a.recvIdx = a.recvIdx, a.attackerIdx
	a.attackList, a.recvList = a.recvList, a.attackList
}

func (a *Abattoir) cloneRetinueList(list *list.DoubleLinkList[entity.Retinue]) *list.DoubleLinkList[entity.Retinue] {
	r := OM.NewList()
	r.Clear()

	list.Each(func(data *entity.Retinue) bool {
		target := OM.NewRetinue()
		data.CloneTo(target)
		r.Add(target)
		return true
	})

	return r
}
