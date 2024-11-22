package protocol

import (
	"encoding/binary"
	"time"
)

const (
	HeadIndexType   = 6
	HeadIndexAction = 7
)

const (
	headServer = 0b10000000
	headClient = 0b00000000
)

var baseTime = time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)

/*
bit,
0 Server 1/Client 0,
1-31 timestamp from 2024-10-01,
31-47 body length,
47-55 action type,
55-63 action
*/
type Head []byte

func NewHead() Head {
	return make([]byte, 8)
}

func ParseHead(b []byte) Head {
	return b[:8]
}

func (h Head) Server() Head {
	h[0] = h[0]&0b01111111 | headServer
	return h
}

func (h Head) Client() Head {
	h[0] = h[0]&0b01111111 | headClient
	return h
}

func (h Head) SetTimestamp() Head {
	bit := h[0]
	t := time.Now().Unix() - baseTime.Unix()
	binary.BigEndian.PutUint32(h[0:4], uint32(t))
	h[0] = h[0] | bit
	return h
}

func (h Head) Timestamp() time.Time {
	t := binary.BigEndian.Uint32(h[0:4])
	return baseTime.Add(time.Duration(t<<1>>1) * time.Second)
}

func (h Head) SetBodyLength(l uint16) Head {
	binary.BigEndian.PutUint16(h[4:6], l)
	return h
}

func (h Head) BodyLength() uint16 {
	return binary.BigEndian.Uint16(h[4:6])
}

func (h Head) SetType(t Type) Head {
	h[HeadIndexType] = byte(t)
	return h
}

func (h Head) Type() Type {
	return Type(h[HeadIndexType])
}

func (h Head) SetAction(a Action) Head {
	h[HeadIndexAction] = byte(a)
	return h
}

func (h Head) Action() Action {
	return Action(h[HeadIndexAction])
}
