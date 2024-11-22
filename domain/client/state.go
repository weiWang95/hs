package client

import (
	"context"

	"hs/pkg/protocol"
)

type IState interface {
	OnEnter(ctx context.Context) error
	OnExit(ctx context.Context) error

	OnOperate(ctx context.Context, cmd string) (bool, error)
	OnRecv(ctx context.Context, cmd protocol.Command) error
}
