package sender

import (
	"context"
	"hs/domain"
	"hs/pkg/protocol"
	"hs/repository"
	"hs/repository/dao"
)

type Sender struct {
	connRepo repository.ConnectRepo
}

func NewSender() domain.ISender {
	return &Sender{
		connRepo: dao.ConnectRepo,
	}
}

func (s *Sender) Send(ctx context.Context, userId uint64, cmd protocol.Command) (bool, error) {
	c, err := s.connRepo.Find(ctx, userId)
	if err != nil {
		return false, err
	}
	if c == nil {
		return false, nil
	}

	_, err = c.Write(cmd.Encode())
	return true, err
}
