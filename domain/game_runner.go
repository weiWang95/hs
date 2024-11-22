package domain

import (
	"context"

	"hs/repository/entity"
)

type IGameRunner interface {
	Run(ctx context.Context) error
	State() entity.GameState

	FindPlayer(ctx context.Context, id uint64) *entity.Player

	UpgradeShop(ctx context.Context, player *entity.Player) error
	RefreshShop(ctx context.Context, player *entity.Player) error
	BuyCard(ctx context.Context, player *entity.Player, idx int) error
	SellCard(ctx context.Context, player *entity.Player, idx int) error
	PlaceRetinue(ctx context.Context, player *entity.Player, cardIdx, retinueIdx, targetIdx int) error
	DragRetinue(ctx context.Context, player *entity.Player, from, to int) error
}
