package matcher

import (
	"context"
	"time"

	"hs/domain"
	"hs/pkg/protocol"
	"hs/repository"
	"hs/repository/dao"

	"github.com/sirupsen/logrus"
)

type matcherImpl struct {
	game   domain.IGame
	sender domain.ISender

	repo repository.MatcherRepo
}

func NewMatcher(game domain.IGame, sender domain.ISender) domain.IMatcher {
	return &matcherImpl{
		repo:   dao.MatcherDao,
		game:   game,
		sender: sender,
	}
}

func (m *matcherImpl) Start(ctx context.Context) {
	cfg := dao.ConfigRepo.GetServerConfig()

	logrus.Debug("Matcher start!!!")
	for {
		time.Sleep(time.Second)

		select {
		case <-ctx.Done():
			logrus.Debugf("Matcher stop!!! %+v", ctx.Err())
			return
		default:
			if m.game.IsFull(ctx) {
				logrus.Debug("game cap is over")
				continue
			}

			count, err := m.repo.GetCount(ctx)
			if err != nil {
				logrus.Errorf("get count error: %v", err)
				continue
			}
			if count < int64(cfg.MinPlayers) {
				continue
			}

			userIds, err := m.repo.BulkPop(ctx, cfg.MinPlayers)
			if err != nil {
				logrus.Errorf("bulk pop error: %v", err)
				continue
			}

			gameId, err := m.game.Start(ctx, userIds)
			if err != nil {
				logrus.Errorf("start game error: %v", err)
				continue
			}

			for _, userId := range userIds {
				if _, err := m.sender.Send(ctx, userId, protocol.New().Matched(gameId)); err != nil {
					logrus.Errorf("send matched error: %v", err)
				}
			}

			logrus.Infof("match queue user count: %d", count-int64(cfg.MinPlayers))
		}
	}
}

func (m *matcherImpl) Join(ctx context.Context, userId uint64) error {
	logrus.Debugf("join match: %d", userId)
	if err := m.repo.Add(ctx, userId); err != nil {
		return err
	}

	return nil
}

func (m *matcherImpl) Quit(ctx context.Context, userId uint64) error {
	logrus.Debugf("quit match: %d", userId)
	if err := m.repo.Remove(ctx, userId); err != nil {
		return err
	}
	return nil
}
