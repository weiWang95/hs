package entity

type ServerConfig struct {
	MaxGameCap                  int `json:"max_game_cap"`                   // 游戏最大容量
	MinPlayers                  int `json:"min_players"`                    // 最小玩家数量
	MaxRound                    int `json:"max_round"`                      // 最大轮数
	MaxPlayerRetinue            int `json:"max_player_retinue"`             // 玩家最大随从数量
	MaxPlayerCard               int `json:"max_player_card"`                // 玩家最大卡牌数量
	MaxPlayerGold               int `json:"max_player_gold"`                // 玩家最大金币数量
	MaxShopLevel                int `json:"max_shop_level"`                 // 商店最大等级
	GameRoundBaseDuration       int `json:"game_round_base_duration"`       // 回合基础时间
	GameRoundAdditionalDuration int `json:"game_round_additional_duration"` // 回合附加时间
}
