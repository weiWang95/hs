package protocol

import (
	"fmt"
	"hs/pkg/utils"

	"github.com/sirupsen/logrus"
)

type Type uint8

const (
	TypeGameOperation Type = iota + 1
	TypeUIOperation
	TypeServerOperation
)

type Action uint8

const (
	BuyCard Action = iota + 1
	UseCard
	SellCard
	DragCard
	UpgradeShop
	RefreshShop
)

const (
	JoinMatch Action = iota + 1
	QuitMatch
)

const (
	Ok Action = iota + 1
	Matched
	UserData
	GameStateChanged
	GameDataSync
	GameOver
)

type Protocol struct {
	Coder
}

func New() *Protocol {
	return &Protocol{Coder: Default}
}

type Command struct {
	Type   Type
	Action Action
	Data   []byte
}

func (c Command) Encode() []byte {
	head := NewHead()
	if c.Type == TypeServerOperation {
		head.Server()
	} else {
		head.Client()
	}

	head.SetTimestamp().
		SetType(c.Type).
		SetAction(c.Action).
		SetBodyLength(uint16(len(c.Data)))
	return append(head, c.Data...)
}

func (c Command) String() string {
	return fmt.Sprintf("{Type: %d, Acion: %d}", c.Type, c.Action)
}

func Parse(b []byte) Command {
	head := ParseHead(b)
	return Command{head.Type(), head.Action(), b[8:]}
}

func (p *Protocol) Ok() Command {
	return p.operation(TypeServerOperation, Ok, nil)
}

func (p *Protocol) Matched(gameId uint64) Command {
	return p.operation(TypeServerOperation, Matched, utils.Uint64ToBytes(gameId))
}

func (p *Protocol) UserData(data any) Command {
	bs, err := p.Encode(data)
	if err != nil {
		logrus.WithField("action", "UserData").Errorf("encode error: %v", err)
	}
	return p.operation(TypeServerOperation, UserData, bs)
}

func (p *Protocol) GameStateChanged(state uint8) Command {
	return p.operation(TypeServerOperation, GameStateChanged, utils.Uint8ToBytes(state))
}

func (p *Protocol) GameDataSync(data any) Command {
	bs, err := p.Encode(data)
	if err != nil {
		logrus.WithField("action", "GameDataSync").Errorf("encode error: %v", err)
	}
	return p.operation(TypeServerOperation, GameDataSync, bs)
}

func (p *Protocol) GameOver(data any) Command {
	bs, err := p.Encode(data)
	if err != nil {
		logrus.WithField("action", "GameOver").Errorf("encode error: %v", err)
	}
	return p.operation(TypeServerOperation, GameOver, bs)
}

func (p *Protocol) JoinMatch() Command {
	return p.operation(TypeUIOperation, JoinMatch, nil)
}

func (p *Protocol) QuitMatch() Command {
	return p.operation(TypeUIOperation, QuitMatch, nil)
}

func (p *Protocol) BuyCard(idx uint8) Command {
	return p.operation(TypeGameOperation, BuyCard, utils.Uint8ToBytes(idx))
}

func (p *Protocol) UseCard(idx uint8, toIdx uint8, target uint8) Command {
	data := append(append(utils.Uint8ToBytes(idx), utils.Uint8ToBytes(toIdx)...), utils.Uint8ToBytes(target)...)
	return p.operation(TypeGameOperation, UseCard, data)
}

func (p *Protocol) DragCard(idx uint8, toIdx uint8) Command {
	data := append(utils.Uint8ToBytes(idx), utils.Uint8ToBytes(toIdx)...)
	return p.operation(TypeGameOperation, DragCard, data)
}

func (p *Protocol) SellCard(idx uint8) Command {
	return p.operation(TypeGameOperation, SellCard, utils.Uint8ToBytes(idx))
}

func (p *Protocol) UpgradeShop() Command {
	return p.operation(TypeGameOperation, UpgradeShop, nil)
}

func (p *Protocol) RefreshShop() Command {
	return p.operation(TypeGameOperation, RefreshShop, nil)
}

func (p *Protocol) operation(t Type, a Action, data []byte) Command {
	return Command{t, a, data}
}
