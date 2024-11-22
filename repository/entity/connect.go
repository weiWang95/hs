package entity

import "net"

type Connect struct {
	net.Conn

	UserId uint64
	GameId uint64
}
