package dispatcher

import (
	"context"

	"hs/domain"
	"hs/pkg/protocol"
	"hs/pkg/session"
	"hs/pkg/utils"
	"hs/repository/dao"

	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	g domain.IGame
	m domain.IMatcher
	s domain.ISender
}

func NewDispatcher(g domain.IGame, m domain.IMatcher, s domain.ISender) domain.IDispatcher {
	return &Dispatcher{
		g: g,
		m: m,
		s: s,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, command []byte) error {
	cmd := protocol.Parse(command)
	logrus.Debugf("recv cmd: %+v", cmd)
	switch cmd.Type {
	case protocol.TypeUIOperation:
		return d.handleUI(ctx, cmd)
	case protocol.TypeGameOperation:
		return d.handleGame(ctx, cmd)
	default:
		logrus.Errorf("unknown type: %+v", cmd)
		return nil
	}
}

func (d *Dispatcher) handleUI(ctx context.Context, cmd protocol.Command) (err error) {
	s := session.Get(ctx)

	switch cmd.Action {
	case protocol.JoinMatch:
		err = d.m.Join(ctx, s.UserId)
	case protocol.QuitMatch:
		err = d.m.Quit(ctx, s.UserId)
	default:
		logrus.Errorf("unknown ui action: %+v", cmd)
		return nil
	}
	if err != nil {
		return err
	}

	_, err = d.SendOk(ctx, s.UserId)
	return err
}

func (d *Dispatcher) handleGame(ctx context.Context, cmd protocol.Command) error {
	s := session.Get(ctx)

	conn, err := dao.ConnectRepo.Find(ctx, s.UserId)
	if err != nil {
		return err
	}
	if conn == nil {
		logrus.Infof("user connect not found: %v", s.UserId)
		return nil
	}

	r, err := d.g.GetRunner(ctx, conn.GameId)
	if err != nil {
		return err
	}
	if r == nil {
		logrus.Infof("user game not found: %v:%v", s.UserId, conn.GameId)
		return nil
	}
	player := r.FindPlayer(ctx, s.UserId)

	switch cmd.Action {
	case protocol.BuyCard:
		err = r.BuyCard(ctx, player, int(utils.BytesToUint8(cmd.Data)))
	case protocol.UseCard:
		err = r.PlaceRetinue(ctx, player, int(cmd.Data[0]), int(cmd.Data[1]), int(cmd.Data[2]))
	case protocol.SellCard:
		err = r.SellCard(ctx, player, int(cmd.Data[0]))
	case protocol.DragCard:
		err = r.DragRetinue(ctx, player, int(cmd.Data[0]), int(cmd.Data[1]))
	case protocol.UpgradeShop:
		err = r.UpgradeShop(ctx, player)
	case protocol.RefreshShop:
		err = r.RefreshShop(ctx, player)
	default:
		logrus.Errorf("unknown game action: %+v", cmd)
		return nil
	}
	if err != nil {
		return err
	}

	_, err = d.SendOk(ctx, s.UserId)
	return err
}

func (d *Dispatcher) SendOk(ctx context.Context, userId uint64) (bool, error) {
	return d.s.Send(ctx, userId, protocol.New().Ok())
}
