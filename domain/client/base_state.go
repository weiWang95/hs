package client

import (
	"context"

	"hs/pkg/protocol"
)

type baseState struct {
	g *Game
}

func NewBaseState(g *Game) baseState {
	return baseState{g: g}
}

func (n *baseState) OnEnter(ctx context.Context) error {
	return nil
}

func (n *baseState) OnExit(ctx context.Context) error {
	return nil
}

func (n *baseState) OnOperate(ctx context.Context, cmd string) (bool, error) {
	return false, nil
}

func (n *baseState) OnRecv(ctx context.Context, cmd protocol.Command) error {
	if cmd.Type != protocol.TypeServerOperation {
		return nil
	}

	switch cmd.Action {
	case protocol.UserData:
		return protocol.Default.Decode(cmd.Data, n.g.user)
	case protocol.Ok:
		n.g.Ok()
	}

	return nil
}
