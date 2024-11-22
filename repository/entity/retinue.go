package entity

import (
	"fmt"
	"hs/pkg/list"

	"github.com/sirupsen/logrus"
)

func init() {
	list.RegisterValue(Retinue{})
}

type Grade int

const (
	GradeNormal Grade = iota
	GradeRare
	GradeEpic
	GradeLegend
)

type Retinue struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	Grade  Grade  `json:"grade"`
	Level  int    `json:"level"`
	Hp     int32  `json:"hp"`
	Attack int32  `json:"attack"`

	Buffs           []*BuffWithTimes `json:"buffs"`
	AdditionalBuffs []*BuffWithTimes `json:"additional_buffs"`

	BuffMap map[BuffEvent][]*BuffWithTimes `json:"buff_map"`

	FinalHp     int32 `json:"final_hp"`
	FinalAttack int32 `json:"final_attack"`
}

type BuffWithTimes struct {
	Times uint8 `json:"times"`
	Buff  IBuff `json:"buff"`
}

func (r *Retinue) Init() {
	if r.Buffs == nil {
		r.Buffs = make([]*BuffWithTimes, 0)
	}
	r.AdditionalBuffs = make([]*BuffWithTimes, 0)
	r.BuffMap = make(map[BuffEvent][]*BuffWithTimes)
	for i := range r.Buffs {
		event := r.Buffs[i].Buff.Base().Event
		if _, ok := r.BuffMap[event]; !ok {
			r.BuffMap[event] = make([]*BuffWithTimes, 0)
		}
		r.BuffMap[event] = append(r.BuffMap[event], r.Buffs[i])
	}
}

func (r *Retinue) AddBuff(b IBuff) {
	for _, v := range r.AdditionalBuffs {
		if v.Buff.Base().Id == b.Base().Id {
			logrus.Debugf("buff exist, append: %v", b.Base().Id)
			v.Times++

			b.CaluateAttr(r)
			return
		}
	}

	event := b.Base().Event
	bwt := BuffWithTimes{Times: 1, Buff: b}
	r.AdditionalBuffs = append(r.AdditionalBuffs, &bwt)
	if _, ok := r.BuffMap[event]; !ok {
		r.BuffMap[event] = make([]*BuffWithTimes, 0)
	}
	r.BuffMap[event] = append(r.BuffMap[event], &bwt)

	b.CaluateAttr(r)
}

func (r *Retinue) EachBuff(event BuffEvent, fn func(b IBuff) bool) bool {
	for i := range r.BuffMap[event] {
		for _ = range r.BuffMap[event][i].Times {
			if !fn(r.BuffMap[event][i].Buff) {
				return false
			}
		}
	}
	return true
}

func (r *Retinue) Refresh() {
	r.FinalAttack, r.FinalHp = r.Attack, r.Hp

	for _, v := range r.BuffMap[BuffEventNull] {
		for _ = range v.Times {
			v.Buff.CaluateAttr(r)
		}
	}
}

func (r *Retinue) CloneTo(p *Retinue) {
	p.Id = r.Id
	p.Name = r.Name
	p.Grade = r.Grade
	p.Level = r.Level
	p.Hp = r.Hp
	p.Attack = r.Attack
	p.FinalAttack = r.FinalAttack
	p.FinalHp = r.FinalHp

	p.Buffs = make([]*BuffWithTimes, 0, len(r.Buffs))
	for _, v := range r.Buffs {
		p.Buffs = append(p.Buffs, &BuffWithTimes{Times: v.Times, Buff: v.Buff})
	}

	p.AdditionalBuffs = make([]*BuffWithTimes, 0, len(r.AdditionalBuffs))
	for _, v := range r.AdditionalBuffs {
		p.AdditionalBuffs = append(p.AdditionalBuffs, &BuffWithTimes{Times: v.Times, Buff: v.Buff})
	}

	p.BuffMap = make(map[BuffEvent][]*BuffWithTimes, len(r.BuffMap))
	for k, v := range r.BuffMap {
		p.BuffMap[k] = make([]*BuffWithTimes, 0, len(v))
		for _, v2 := range v {
			p.BuffMap[k] = append(p.BuffMap[k], &BuffWithTimes{Times: v2.Times, Buff: v2.Buff})
		}
	}
}

func (r *Retinue) Inspect() string {
	return fmt.Sprintf(" [%d-%d (%d-%d)] ", r.FinalAttack, r.FinalHp, r.Attack, r.Hp)
}
