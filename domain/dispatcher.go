package domain

import (
	"context"
)

type IDispatcher interface {
	Dispatch(ctx context.Context, command []byte) error
}
