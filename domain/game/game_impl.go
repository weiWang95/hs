package game

import (
	"context"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"hs/domain"
	"hs/pkg/errors"
	"hs/pkg/safe"
	"hs/pkg/xid"
	"hs/repository"
	"hs/repository/dao"
	"hs/repository/entity"

	"github.com/sirupsen/logrus"
)

type gameImpl struct {
	sender domain.ISender

	userRepo repository.UserRepo
	gameRepo repository.GameRepo
	cfgRepo  repository.ConfigRepo

	om *ObjectManager

	cap       int64
	total     int64
	runnerMap sync.Map
}

func NewGame(sender domain.ISender) domain.IGame {
	cfg := dao.ConfigRepo.GetServerConfig()

	return &gameImpl{
		sender:   sender,
		userRepo: dao.UserRepo,
		gameRepo: dao.NewGameRepo(),
		cfgRepo:  dao.ConfigRepo,
		om:       OM,
		cap:      int64(cfg.MaxGameCap),
	}
}

func (g *gameImpl) GetGameTotal(ctx context.Context) int64 {
	return g.total
}

func (g *gameImpl) IsFull(ctx context.Context) bool {
	return g.total >= g.cap
}

func (g *gameImpl) GetRunner(ctx context.Context, id uint64) (domain.IGameRunner, error) {
	runner, ok := g.runnerMap.Load(id)
	if !ok {
		return nil, errors.NewError(errors.CodeNotFound, "game runner not found")
	}

	return runner.(*GameRunner), nil
}

func (g *gameImpl) Start(ctx context.Context, userIds []uint64) (uint64, error) {
	players, err := g.userRepo.BulkGetUsers(ctx, userIds)
	if err != nil {
		return 0, err
	}

	game := g.initGame(players)

	logrus.Infof("start game[%d]: %+v", game.Id, players)
	if err := g.gameRepo.SaveGame(ctx, game); err != nil {
		return 0, err
	}

	for id, _ := range game.Players {
		if err := dao.ConnectRepo.SetGameId(ctx, id, game.Id); err != nil {
			return 0, err
		}
	}

	runner := g.om.NewRunner()
	runner.Game = game
	runner.sender = g.sender
	g.runnerMap.Store(runner.Id, runner)
	atomic.AddInt64(&g.total, 1)

	safe.Go(func() error {
		ctx := context.TODO()
		defer g.ClearGame(ctx, runner)

		return g.RunGame(ctx, runner)
	})

	return game.Id, nil
}

func (g *gameImpl) initGame(players []entity.User) *entity.Game {
	game := g.om.NewGame()
	game.Id = xid.New()
	game.State = entity.GamePending
	game.Round = 0

	game.Players = make(map[uint64]*entity.Player, len(players))
	game.Shop = make(map[uint64]*entity.Shop, len(players))

	for _, player := range players {
		p := g.om.NewPlayer()
		p.Id = player.Id
		p.Nickname = player.Nickname
		p.State = entity.PlayerAlive
		p.MaxHp = 30
		p.Hp = 30
		p.Shield = 0
		p.MaxGold = 10 // 2
		p.Gold = 10    // 2
		p.CardList = g.om.NewList()
		p.RetinueList = g.om.NewList()
		game.Players[player.Id] = p

		shop := g.om.NewShop()
		shop.Level = 2 // 1
		cfg := g.cfgRepo.GetShopConfig().Levels[shop.Level-1]
		shop.UpgradeCost = cfg.UpgradeCost
		shop.RetinueCap = cfg.RetinueCap
		shop.Retinue = g.om.NewList()
		game.Shop[player.Id] = shop
	}

	return game
}

func (g *gameImpl) RunGame(ctx context.Context, runner *GameRunner) error {
	if err := runner.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (g *gameImpl) ClearGame(ctx context.Context, runner *GameRunner) {
	if e := recover(); e != nil {
		debug.PrintStack()
		logrus.Errorf("game runner panic: %v", e)
	}

	if err := g.gameRepo.DelGame(ctx, runner.Game.Id); err != nil {
		logrus.WithFields(logrus.Fields{"game_id": runner.Game.Id}).Errorf("clear game error: %v", err)
	}

	g.runnerMap.Delete(runner.Id)
	atomic.AddInt64(&g.total, -1)
	g.om.PutRunner(runner)
}
