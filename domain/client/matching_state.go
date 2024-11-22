package client

import (
	"context"
	"fmt"
	"hs/pkg/protocol"
)

type matchingState struct {
	baseState
}

func NewMatchingState(g *Game) IState {
	return &matchingState{
		baseState: NewBaseState(g),
	}
}

func (n *matchingState) OnEnter(ctx context.Context) error {
	fmt.Println("matching")
	return n.g.JoinMatch(ctx)
}

func (n *matchingState) OnOperate(ctx context.Context, cmd string) (bool, error) {
	switch cmd {
	case "quit", "q":
		if err := n.g.QuitMatch(ctx); err != nil {
			return false, err
		}
		fmt.Println("quit match")

		return true, n.g.SwitchState(ctx, NewNormalState(n.g))
	}
	return false, nil
}

func (n *matchingState) OnRecv(ctx context.Context, cmd protocol.Command) error {
	if cmd.Type != protocol.TypeServerOperation {
		return n.baseState.OnRecv(ctx, cmd)
	}

	switch cmd.Action {
	case protocol.Matched:
		if err := n.g.GameMatched(ctx, cmd.Data); err != nil {
			return err
		}
		return n.g.SwitchState(ctx, NewGamingState(n.g))
	default:
		return n.baseState.OnRecv(ctx, cmd)
	}
}
