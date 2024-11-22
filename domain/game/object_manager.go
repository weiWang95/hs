package game

import (
	"hs/pkg/list"
	"hs/repository/entity"
	"sync"
)

var OM = NewObjectManager()

type ObjectManager struct {
	gamePool     *sync.Pool
	runnerPool   *sync.Pool
	playerPool   *sync.Pool
	shopPool     *sync.Pool
	listPool     *sync.Pool
	abattoirPool *sync.Pool
	retinuePool  *sync.Pool
}

func NewObjectManager() *ObjectManager {
	m := ObjectManager{
		gamePool:     &sync.Pool{},
		runnerPool:   &sync.Pool{},
		playerPool:   &sync.Pool{},
		shopPool:     &sync.Pool{},
		listPool:     &sync.Pool{},
		abattoirPool: &sync.Pool{},
		retinuePool:  &sync.Pool{},
	}
	m.init()
	return &m
}

func (m *ObjectManager) init() {
	m.gamePool.New = m.gameNew
	m.runnerPool.New = m.runnerNew
	m.playerPool.New = m.playerNew
	m.shopPool.New = m.shopNew
	m.listPool.New = m.listNew
	m.abattoirPool.New = m.abattoirNew
	m.retinuePool.New = m.retinueNew
}

func (m *ObjectManager) gameNew() any {
	return &entity.Game{}
}

func (m *ObjectManager) NewGame() *entity.Game {
	g := m.gamePool.Get().(*entity.Game)
	return g
}

func (m *ObjectManager) PutGame(g *entity.Game) {
	m.clearGame(g)
	m.gamePool.Put(g)
}

func (m *ObjectManager) clearGame(g *entity.Game) {
	for _, v := range g.Players {
		m.playerPool.Put(v)
	}
	clear(g.Players)

	for key := range g.Shop {
		m.shopPool.Put(g.Shop[key])
	}
	clear(g.Shop)
}

func (m *ObjectManager) runnerNew() any {
	return &GameRunner{}
}

func (m *ObjectManager) NewRunner() *GameRunner {
	r := m.runnerPool.Get().(*GameRunner)
	return r
}

func (m *ObjectManager) PutRunner(r *GameRunner) {
	m.clearRunner(r)
	m.runnerPool.Put(r)
}

func (m *ObjectManager) clearRunner(r *GameRunner) {
	if r.Game != nil {
		m.clearGame(r.Game)
		r.Game = nil
	}
}

func (m *ObjectManager) shopNew() any {
	return &entity.Shop{}
}

func (m *ObjectManager) NewShop() *entity.Shop {
	s := m.shopPool.Get().(*entity.Shop)
	return s
}

func (m *ObjectManager) PutShop(s *entity.Shop) {
	if s.Retinue != nil {
		m.listPool.Put(s.Retinue)
		s.Retinue = nil
	}

	m.shopPool.Put(s)
}

func (m *ObjectManager) playerNew() any {
	return &entity.Player{}
}

func (m *ObjectManager) NewPlayer() *entity.Player {
	p := m.playerPool.Get().(*entity.Player)
	return p
}

func (m *ObjectManager) PutPlayer(p *entity.Player) {
	m.clearPlayer(p)
	m.playerPool.Put(p)
}

func (m *ObjectManager) clearPlayer(p *entity.Player) {
	if p.CardList != nil {
		m.listPool.Put(p.CardList)
		p.CardList = nil
	}
	if p.RetinueList != nil {
		m.listPool.Put(p.RetinueList)
		p.RetinueList = nil
	}
}

func (m *ObjectManager) listNew() any {
	return list.NewDoubleLinkList[entity.Retinue]()
}

func (m *ObjectManager) NewList() *list.DoubleLinkList[entity.Retinue] {
	l := m.listPool.Get().(*list.DoubleLinkList[entity.Retinue])
	return l
}

func (m *ObjectManager) PutList(l *list.DoubleLinkList[entity.Retinue]) {
	l.Clear()
	m.listPool.Put(l)
}

func (m *ObjectManager) abattoirNew() any {
	return &Abattoir{}
}

func (m *ObjectManager) NewAbattoir() *Abattoir {
	a := m.abattoirPool.Get().(*Abattoir)
	return a
}

func (m *ObjectManager) PutAbattoir(a *Abattoir) {
	m.clearAbattoir(a)
	m.abattoirPool.Put(a)
}

func (m *ObjectManager) clearAbattoir(a *Abattoir) {
	if a.attackList != nil {
		m.listPool.Put(a.attackList)
		a.attackList = nil
	}
	if a.recvList != nil {
		m.listPool.Put(a.recvList)
		a.recvList = nil
	}
}

func (m *ObjectManager) retinueNew() any {
	return &entity.Retinue{}
}

func (m *ObjectManager) NewRetinue() *entity.Retinue {
	r := m.retinuePool.Get().(*entity.Retinue)
	return r
}

func (m *ObjectManager) PutRetinue(r *entity.Retinue) {
	m.retinuePool.Put(r)
}
