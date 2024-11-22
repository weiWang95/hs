package cmd

import (
	"context"
	"net"

	"hs/domain"
	"hs/domain/client"

	"github.com/sirupsen/logrus"
)

type Client struct {
	cli domain.IClient

	host string
	port int
}

func NewClient(host string, port int) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

func (c *Client) Run(ctx context.Context) error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(c.host), Port: c.port})
	if err != nil {
		return err
	}
	defer conn.Close()
	logrus.Infof("client start: %s", conn.RemoteAddr().String())

	c.cli = client.NewClient(conn)
	return c.cli.Run(ctx)
}

func (c *Client) Shutdown() error {
	return c.cli.Close()
}
