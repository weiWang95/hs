package client

import (
	"bufio"
	"context"
	"fmt"
	"hs/domain"
	"hs/pkg/protocol"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

type Client struct {
	conn net.Conn

	g  *Game
	ok chan struct{}
}

func NewClient(conn net.Conn) domain.IClient {
	c := &Client{conn: conn}
	c.g = NewGame(c)
	c.ok = make(chan struct{})
	return c
}

func (c *Client) Run(ctx context.Context) error {
	defer c.Close()

	done := make(chan struct{})
	closeDone := func() {
		if done != nil {
			close(done)
			done = nil
		}
	}
	go func() {
		defer closeDone()

		if err := c.handleScan(ctx); err != nil && err != ErrQuit {
			logrus.WithError(err).Error("handle scan error")
		}
	}()

	go func() {
		defer closeDone()

		if err := c.handleRecv(ctx); err != nil {
			logrus.WithError(err).Debug("handle recv error")
		}
	}()

	c.g.SwitchState(ctx, NewNormalState(c.g))

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		done = nil
	}

	return nil
}

func (c *Client) Close() error {
	if c.ok != nil {
		close(c.ok)
		c.ok = nil
	}

	return c.conn.Close()
}

func (c *Client) Send(ctx context.Context, cmd protocol.Command) error {
	logrus.Debugf("send: %+v", cmd)
	_, err := c.conn.Write(cmd.Encode())
	return err
}

func (c *Client) Ok() {
	c.ok <- struct{}{}
}

func (c *Client) handleRecv(ctx context.Context) error {
	buffer := make([]byte, 65535)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := c.conn.Read(buffer)
			if err != nil {
				return err
			}

			logrus.Debugf("recv: %s", string(buffer[:n]))
			cmd := protocol.Parse(buffer[:n])
			logrus.Debugf("cmd: %+v", cmd)

			if err := c.g.OnRecv(ctx, cmd); err != nil {
				return err
			}
		}
	}
}

func (c *Client) handleScan(ctx context.Context) error {
	scan := bufio.NewScanner(os.Stdin)
	var cmd string
	for {
		fmt.Println("please input command:")
		if !scan.Scan() {
			break
		}
		cmd = scan.Text()

		logrus.Debugf("scan: %s", cmd)

		ok, err := c.g.OnOperate(ctx, cmd)
		if err != nil {
			if err == ErrCannotOperate {
				continue
			}
			return err
		}
		if !ok {
			fmt.Println("unknown command")
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.ok:
			continue
		}
	}
	return nil
}
