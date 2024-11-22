package cmd

import (
	"context"
	"net"

	"hs/domain"
	"hs/domain/connector"
	"hs/domain/dispatcher"
	"hs/domain/game"
	"hs/domain/matcher"
	"hs/domain/sender"

	"github.com/sirupsen/logrus"
)

var (
	g domain.IGame
	m domain.IMatcher
	d domain.IDispatcher
	s domain.ISender
)

func init() {
	s = sender.NewSender()
	g = game.NewGame(s)
	m = matcher.NewMatcher(g, s)
	d = dispatcher.NewDispatcher(g, m, s)
}

type Server struct {
	net.Listener

	port int
	ipv6 bool
}

func NewServer(port int, ipv6 bool) *Server {
	return &Server{
		port: port,
		ipv6: ipv6,
	}
}

func (s *Server) Start(ctx context.Context) error {
	ip := net.IPv4zero
	if s.ipv6 {
		ip = net.IPv6zero
	}

	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   ip,
		Port: s.port,
	})
	if err != nil {
		return err
	}
	defer l.Close()
	s.Listener = l
	logrus.Infof("server start: %s", l.Addr().String())

	go m.Start(ctx)

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) Shutdown() error {
	return s.Listener.Close()
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	logrus.Infof("new connection: %s", conn.RemoteAddr().String())

	c := connector.NewConnector(conn)
	if err := c.Listen(ctx, d); err != nil {
		logrus.Errorf("handle connection error: %v", err)
	}
}
