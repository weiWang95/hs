package client

import (
	"context"
	"hs/pkg/protocol"
	"hs/pkg/utils"
	"hs/repository/dao"
	"hs/repository/entity"
	"hs/service"
	"reflect"

	"github.com/sirupsen/logrus"
)

type GameData struct {
	Id    uint64           `json:"id"`
	State entity.GameState `json:"state"`
	Round int              `json:"round"`

	Player *entity.Player `json:"player"`
	Shop   *entity.Shop   `json:"shop"`
}

type Game struct {
	cli   *Client
	state IState
	user  *entity.User

	ok chan struct{}

	data   *GameData
	cfg    entity.ServerConfig
	drawer GameDrawer
}

func NewGame(c *Client) *Game {
	g := new(Game)
	g.cli = c
	g.ok = make(chan struct{})
	g.cfg = dao.ConfigRepo.GetServerConfig()
	g.drawer = NewCmdDrawer()
	return g
}

func (g *Game) SwitchState(ctx context.Context, state IState) error {
	if g.state != nil && !reflect.ValueOf(g.state).IsNil() {
		if err := g.state.OnExit(ctx); err != nil {
			return err
		}
	}
	logrus.Debugf("switch state: %T -> %T", g.state, state)
	g.state = state

	if state != nil && !reflect.ValueOf(g.state).IsNil() {
		if err := g.state.OnEnter(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) OnRecv(ctx context.Context, cmd protocol.Command) error {
	return g.state.OnRecv(ctx, cmd)
}

func (g *Game) OnOperate(ctx context.Context, cmd string) (bool, error) {
	return g.state.OnOperate(ctx, cmd)
}

func (g *Game) Ok() {
	g.cli.Ok()
}

func (g *Game) JoinMatch(ctx context.Context) error {
	return g.cli.Send(ctx, protocol.New().JoinMatch())
}

func (g *Game) QuitMatch(ctx context.Context) error {
	return g.cli.Send(ctx, protocol.New().QuitMatch())
}

func (g *Game) BuyCard(ctx context.Context, idx int) error {
	if !service.CanBuyCard(ctx, g.cfg, g.data.State, g.data.Shop, g.data.Player, int(idx)) {
		return ErrCannotOperate
	}

	return g.cli.Send(ctx, protocol.New().BuyCard(uint8(idx)))
}

func (g *Game) UseCard(ctx context.Context, idx int, toIdx int, target int) error {
	if !service.CanPlaceRetinue(ctx, g.cfg, g.data.State, g.data.Shop, g.data.Player, int(idx), int(toIdx), int(target)) {
		return ErrCannotOperate
	}

	return g.cli.Send(ctx, protocol.New().UseCard(uint8(idx), uint8(toIdx), uint8(target)))
}

func (g *Game) SellCard(ctx context.Context, idx int) error {
	if !service.CanSellCard(ctx, g.cfg, g.data.State, g.data.Shop, g.data.Player, int(idx)) {
		return ErrCannotOperate
	}
	return g.cli.Send(ctx, protocol.New().SellCard(uint8(idx)))
}

func (g *Game) DragRetinue(ctx context.Context, idx int, toIdx int) error {
	if !service.CanDragRetinue(ctx, g.cfg, g.data.State, g.data.Shop, g.data.Player, int(idx), int(toIdx)) {
		return ErrCannotOperate
	}

	return g.cli.Send(ctx, protocol.New().DragCard(uint8(idx), uint8(toIdx)))
}

func (g *Game) UpgradeShop(ctx context.Context) error {
	if !service.CanUpgradeShop(ctx, g.cfg, g.data.State, g.data.Shop, g.data.Player) {
		return ErrCannotOperate
	}
	return g.cli.Send(ctx, protocol.New().UpgradeShop())
}

func (g *Game) RefreshShop(ctx context.Context) error {
	if !service.CanRefreshShop(ctx, g.cfg, g.data.State, g.data.Shop, g.data.Player) {
		return ErrCannotOperate
	}
	return g.cli.Send(ctx, protocol.New().RefreshShop())
}

func (g *Game) GameMatched(ctx context.Context, data []byte) error {
	g.data = &GameData{Id: utils.BytesToUint64(data)}
	return nil
}

func (g *Game) GameStateChanged(ctx context.Context, data []byte) error {
	state := utils.BytesToUint8(data)
	g.data.State = entity.GameState(state)
	return nil
}

func (g *Game) GameDataSync(ctx context.Context, data []byte) error {
	defer g.Draw(ctx)
	return protocol.Default.Decode(data, g.data)
}

func (g *Game) GameOver(ctx context.Context, data []byte) error {
	var res entity.GameResult
	if err := protocol.Default.Decode(data, &res); err != nil {
		return err
	}

	g.drawer.ShowGameResult(ctx, &res)
	g.data = nil
	return nil
}

func (g *Game) Draw(ctx context.Context) {
	if err := g.drawer.Draw(ctx, g.data); err != nil {
		logrus.Errorf("draw ui error: %v", err)
	}
}

func (g *Game) DrawHelp(ctx context.Context) {
	g.drawer.ShowHelp(ctx)
}
