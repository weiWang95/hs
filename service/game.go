package service

import (
	"context"

	"hs/repository/entity"
)

func CanOperate(ctx context.Context, gameState entity.GameState, player *entity.Player) bool {
	return player.State == entity.PlayerAlive && gameState == entity.GameWaiting
}

func CanBuyCard(ctx context.Context, serverCfg entity.ServerConfig, gameState entity.GameState, shop *entity.Shop, player *entity.Player, idx int) bool {
	if !CanOperate(ctx, gameState, player) {
		return false
	}
	return player.Gold >= 3 && (player.CardList == nil || player.CardList.Size() < serverCfg.MaxPlayerCard) && (shop.Retinue != nil && shop.Retinue.Size() > idx)
}

func CanPlaceRetinue(ctx context.Context, serverCfg entity.ServerConfig, gameState entity.GameState, shop *entity.Shop, player *entity.Player, cardIdx, retinueIdx, targetIdx int) bool {
	if !CanOperate(ctx, gameState, player) {
		return false
	}

	if player.RetinueList != nil && player.RetinueList.Size() >= serverCfg.MaxPlayerRetinue {
		return false
	}

	if player.CardList == nil || cardIdx < 0 || cardIdx >= player.CardList.Size() || retinueIdx < 0 {
		return false
	}
	if targetIdx < 0 || (targetIdx != 0 && player.RetinueList != nil && targetIdx >= player.RetinueList.Size()) {
		return false
	}

	return true
}

func CanSellCard(ctx context.Context, serverCfg entity.ServerConfig, gameState entity.GameState, shop *entity.Shop, player *entity.Player, idx int) bool {
	return CanOperate(ctx, gameState, player)
}

func CanDragRetinue(ctx context.Context, serverCfg entity.ServerConfig, gameState entity.GameState, shop *entity.Shop, player *entity.Player, from, to int) bool {
	return CanOperate(ctx, gameState, player)
}

func CanUpgradeShop(ctx context.Context, serverCfg entity.ServerConfig, gameState entity.GameState, shop *entity.Shop, player *entity.Player) bool {
	if !CanOperate(ctx, gameState, player) {
		return false
	}
	return shop.Level < serverCfg.MaxShopLevel && player.Gold >= int32(shop.UpgradeCost)
}

func CanRefreshShop(ctx context.Context, serverCfg entity.ServerConfig, gameState entity.GameState, shop *entity.Shop, player *entity.Player) bool {
	if !CanOperate(ctx, gameState, player) {
		return false
	}
	return player.Gold > 0
}
