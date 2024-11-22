package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"hs/pkg/protocol"

	"github.com/sirupsen/logrus"
)

var ErrCannotOperate = errors.New("cannot operate")

type gamingState struct {
	baseState
}

func NewGamingState(g *Game) IState {
	return &gamingState{
		baseState: NewBaseState(g),
	}
}

func (n *gamingState) OnEnter(ctx context.Context) error {
	fmt.Println("Match success, Game Start!!!")
	return nil
}

func (n *gamingState) OnExit(ctx context.Context) error {
	fmt.Println("Game Over!!!")
	return nil
}

func (n *gamingState) OnOperate(ctx context.Context, cmd string) (bool, error) {
	data := strings.Split(cmd, " ")
	logrus.Debugf("[Gaming]OnOperate: %v", data)

	switch data[0] {
	case "buy", "b":
		if len(data) < 2 {
			return false, nil
		}
		idx, _ := strconv.ParseUint(data[1], 10, 8)
		return true, n.g.BuyCard(ctx, int(idx))
	case "use", "u":
		if len(data) < 3 {
			return false, nil
		}
		idx, _ := strconv.ParseUint(data[1], 10, 8)
		to, _ := strconv.ParseUint(data[2], 10, 8)
		var target uint64
		if len(data) > 3 {
			target, _ = strconv.ParseUint(data[3], 10, 8)
		}
		return true, n.g.UseCard(ctx, int(idx), int(to), int(target))
	case "drag", "d":
		if len(data) < 3 {
			return false, nil
		}
		idx, _ := strconv.ParseUint(data[1], 10, 8)
		to, _ := strconv.ParseUint(data[2], 10, 8)
		return true, n.g.DragRetinue(ctx, int(idx), int(to))
	case "sell", "s":
		if len(data) < 2 {
			return false, nil
		}
		idx, _ := strconv.ParseUint(data[1], 10, 8)
		return true, n.g.SellCard(ctx, int(idx))
	case "upgrade", "up":
		return true, n.g.UpgradeShop(ctx)
	case "refresh", "r":
		return true, n.g.RefreshShop(ctx)
	case "help", "h":
		n.g.DrawHelp(ctx)
		return false, nil
	}
	return false, nil
}

func (n *gamingState) OnRecv(ctx context.Context, cmd protocol.Command) error {
	if cmd.Type != protocol.TypeServerOperation {
		return n.baseState.OnRecv(ctx, cmd)
	}

	switch cmd.Action {
	case protocol.GameStateChanged:
		return n.onGameStateChanged(ctx, cmd.Data)
	case protocol.GameDataSync:
		return n.onGameDataSync(ctx, cmd.Data)
	case protocol.GameOver:
		return n.onGameOver(ctx, cmd.Data)
	}

	return n.baseState.OnRecv(ctx, cmd)
}

func (n *gamingState) onGameStateChanged(ctx context.Context, data []byte) error {
	return n.g.GameStateChanged(ctx, data)
}

func (n *gamingState) onGameDataSync(ctx context.Context, data []byte) error {
	return n.g.GameDataSync(ctx, data)
}

func (n *gamingState) onGameOver(ctx context.Context, data []byte) error {
	if err := n.g.GameOver(ctx, data); err != nil {
		logrus.WithError(err).Error("game over error")
	}

	return n.g.SwitchState(ctx, NewNormalState(n.g))
}
