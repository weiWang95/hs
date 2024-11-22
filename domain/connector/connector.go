package connector

import (
	"context"
	"fmt"
	"net"

	"hs/domain"
	"hs/pkg/session"
	"hs/repository/dao"
	"hs/repository/entity"

	"github.com/sirupsen/logrus"
)

type Connector struct {
	net.Conn

	user *entity.User
}

func NewConnector(conn net.Conn) domain.IConnector {
	return &Connector{Conn: conn}
}

func (c *Connector) Listen(ctx context.Context, d domain.IDispatcher) error {
	defer c.close(ctx)

	if err := c.init(ctx); err != nil {
		return err
	}

	ctx = session.Set(ctx, c.user.Id)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			command, err := c.readCommand(ctx)
			if err != nil {
				return err
			}

			if err := d.Dispatch(ctx, command); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (c *Connector) readCommand(ctx context.Context) ([]byte, error) {
	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("recv: %v", buffer[:8])

	return buffer[:n], nil
}

func (c *Connector) init(ctx context.Context) error {
	user, err := dao.UserRepo.CreateUser(ctx)
	if err != nil {
		return err
	}
	c.user = user

	conn := entity.Connect{
		Conn:   c.Conn,
		UserId: user.Id,
	}
	if err := dao.ConnectRepo.Save(ctx, &conn); err != nil {
		return err
	}

	return nil
}

func (c *Connector) close(ctx context.Context) error {
	if c.user != nil {
		if err := dao.ConnectRepo.Remove(ctx, c.user.Id); err != nil {
			logrus.Errorf("remove connect error: %v", err)
		}
	}

	return c.Conn.Close()
}
