package domain

import (
	"context"
	"hs/pkg/protocol"
)

type ISender interface {
	Send(ctx context.Context, userId uint64, cmd protocol.Command) (bool, error)
}
