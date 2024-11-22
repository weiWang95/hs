package domain

import "context"

type IConnector interface {
	Listen(ctx context.Context, d IDispatcher) error
}
