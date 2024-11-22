package client

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrQuit = errors.New("quit")
)

type normalState struct {
	baseState
}

func NewNormalState(g *Game) IState {
	return &normalState{
		baseState: NewBaseState(g),
	}
}

func (n *normalState) OnEnter(ctx context.Context) error {
	fmt.Printf("input 'join' to join game match or 'quit' to quit\n")
	return nil
}

func (n *normalState) OnOperate(ctx context.Context, cmd string) (bool, error) {
	switch cmd {
	case "join", "j":
		return true, n.g.SwitchState(ctx, NewMatchingState(n.g))
	case "quit", "q":
		return true, ErrQuit
	}
	return false, nil
}
