package game

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"hs/domain"
	"hs/pkg/list"
	"hs/pkg/protocol"
	"hs/pkg/safe"
	"hs/pkg/utils"
	"hs/repository/dao"
	"hs/repository/entity"
	"hs/service"

	"github.com/sirupsen/logrus"
)

type GameRunner struct {
	*entity.Game

	sender domain.ISender

	cardPool *CardPool

	serverCfg entity.ServerConfig
	shopCfg   entity.ShopConfig

	alivePlayerIds []uint64
	deadPlayerIds  []uint64

	buffSys    *BuffSystem
	abattoirs  []*Abattoir
	fightOut   sync.Map
	fightAfter sync.Map
}

func NewGameRunner(sender domain.ISender, game *entity.Game) domain.IGameRunner {
	return &GameRunner{
		sender: sender,
		Game:   game,
	}
}

func (r *GameRunner) State() entity.GameState {
	return r.Game.State
}

func (r *GameRunner) switchState(state entity.GameState) {
	r.Game.State = state
}

func (r *GameRunner) Init(ctx context.Context) error {
	r.cardPool = NewCardPool()
	if err := r.cardPool.Init(ctx); err != nil {
		return err
	}

	r.shopCfg = dao.ConfigRepo.GetShopConfig()
	r.serverCfg = dao.ConfigRepo.GetServerConfig()

	r.abattoirs = make([]*Abattoir, 0, len(r.Players)/2)
	r.alivePlayerIds = make([]uint64, 0, len(r.Players))
	r.deadPlayerIds = make([]uint64, 0, len(r.Players))

	r.buffSys = NewBuffSystem(r)
	return nil
}

func (r *GameRunner) Run(ctx context.Context) error {
	logrus.Debug("game runner start")

	if err := r.Init(ctx); err != nil {
		return err
	}

	r.refreshPlayerStateIds(ctx)

	for r.Round < r.serverCfg.MaxRound {
		r.Round += 1
		logrus.WithFields(logrus.Fields{"round": r.Round}).Debug("round start")

		if err := r.roundStart(ctx); err != nil {
			return err
		}
		r.switchState(entity.GameWaiting)

		r.SyncPlayersGameData(ctx)

		if err := r.waitRoundEnd(ctx); err != nil {
			return err
		}
		if err := r.roundEnd(ctx); err != nil {
			return err
		}
		r.switchState(entity.GamePlaying)

		if err := r.fightMatch(ctx); err != nil {
			return err
		}

		if err := r.fightRoundStart(ctx); err != nil {
			return err
		}
		if err := r.fightRoundProgress(ctx); err != nil {
			return err
		}
		if err := r.fightRoundEnd(ctx); err != nil {
			return err
		}

		r.refreshPlayerStateIds(ctx)

		if len(r.alivePlayerIds) <= 1 {
			logrus.Infof("Game Over[%d]", r.Game.Id)
			r.GameOver(ctx)
			return nil
		}

		r.clearAbattoirs()
	}
	return nil
}

func (r *GameRunner) roundStart(ctx context.Context) error {
	logrus.Debugf("############## Round Start: %d ##############", r.Round)
	if err := r.eachPlayers(ctx, func(ctx context.Context, player *entity.Player) error {
		r.increasePlayerGlod(player)      // 水晶成长
		r.decreaseShopUpgradeCost(player) // 商店升级花费降低
		if err := r.refreshPlayerShop(ctx, player); err != nil {
			logrus.WithFields(logrus.Fields{"game_id": r.Id, "player_id": player.Id}).Errorf("refresh player shop error: %v", err)
		}

		logrus.Debugf("player state: %+v", player.Inspect())
		logrus.Debugf("player shop: %+v", r.inspectShop(r.Shop[player.Id]))

		// TODO: remove test
		// if r.CanUpgradeShop(ctx, player) {
		// 	r.UpgradeShop(ctx, player)
		// }

		// for r.CanBuyCard(ctx, player) {
		// 	r.BuyCard(ctx, player, 0)
		// }

		// for player.CardList.Size() > 0 && r.CanPlaceRetinue(ctx, player) {
		// 	r.PlaceRetinue(ctx, player, 0, 0)
		// }

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *GameRunner) waitRoundEnd(ctx context.Context) error {
	logrus.Debugf("############## Wait Round End: %d ##############", r.Round)
	select {
	case <-time.After(r.roundDuration()):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *GameRunner) roundEnd(ctx context.Context) error {
	logrus.Debugf("############## Round End: %d ##############", r.Round)
	return nil
}

func (r *GameRunner) fightMatch(ctx context.Context) error {
	logrus.Debugf("############## Fight Match: %d ##############", r.Round)
	logrus.Debugf("Alive Players:%+v", r.alivePlayerIds)

	utils.Shuffle(r.alivePlayerIds)   // 打乱
	if len(r.alivePlayerIds)%2 == 1 { // 幸运儿
		r.alivePlayerIds = append(r.alivePlayerIds, r.deadPlayerIds[rand.IntN(len(r.deadPlayerIds))])
	}

	for i := 0; i < len(r.alivePlayerIds); i += 2 {
		item := NewAbattoir(r, r.Players[r.alivePlayerIds[i]], r.Players[r.alivePlayerIds[i+1]])
		r.abattoirs = append(r.abattoirs, item)
	}

	return nil
}

func (r *GameRunner) fightRoundStart(ctx context.Context) error {
	logrus.Debugf("############## Fight Round Start: %d ##############", r.Round)
	for i, item := range r.abattoirs {
		logrus.Debugf("start fight round: no:%d player: %v VS %v", i+1, item.attacker, item.recv)
		logrus.Debugf("attacker: %s", r.inspectRetinues(item.attackList))
		logrus.Tracef("Origin attacker: %s", r.inspectRetinues(r.Players[item.attacker].RetinueList))
		logrus.Debugf("recv: %s", r.inspectRetinues(item.recvList))
		logrus.Tracef("Origin recv: %s", r.inspectRetinues(r.Players[item.recv].RetinueList))
	}
	return nil
}

func (r *GameRunner) fightRoundProgress(ctx context.Context) error {
	logrus.Debugf("############## Fight Round Progress: %d [%d] ##############", r.Round, len(r.abattoirs))
	r.fightOut.Clear()
	r.fightAfter.Clear()

	var wg sync.WaitGroup
	for key, _ := range r.abattoirs {
		func(idx int) {
			wg.Add(1)
			safe.Go(func() error {
				defer wg.Done()

				res := r.abattoirs[idx].Fight(ctx)

				var damage int
				var attacker uint64
				for playerId, list := range res {
					// 计算伤害
					if list.Size() != 0 {
						attacker = playerId

						damage += r.Shop[playerId].Level
						list.Each(func(a *entity.Retinue) bool {
							damage += a.Level
							return true
						})
					}
					r.fightAfter.Store(playerId, res[playerId])
				}

				// 记录应受到伤害
				for playerId, _ := range res {
					if playerId != attacker {
						logrus.Debugf("player %v -> %v, damage %v", attacker, playerId, damage)
						r.fightOut.Store(playerId, damage)
						break
					}
				}

				return nil
			})
		}(key)
	}
	wg.Wait()
	return nil
}

func (r *GameRunner) fightRoundEnd(ctx context.Context) error {
	logrus.Debugf("############## Fight Round End: %d ##############", r.Round)
	r.fightOut.Range(func(key, value any) bool {
		playerId := key.(uint64)
		damage := value.(int)

		logrus.Debugf("player %v receive damage %v", playerId, damage)
		r.Players[playerId].Hp -= int32(damage)
		logrus.Debugf("player %v hp %v", playerId, r.Players[playerId].Hp)
		if r.Players[playerId].Hp < 0 { // 死亡
			r.Players[playerId].State = entity.PlayerDead
			logrus.Debugf("player %v dead", playerId)
		}
		return true
	})

	return nil
}

func (r *GameRunner) eachPlayers(ctx context.Context, fn func(ctx context.Context, player *entity.Player) error) error {
	for i, _ := range r.Players {
		if err := fn(ctx, r.Players[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *GameRunner) refreshPlayerStateIds(ctx context.Context) {
	logrus.Debugf("refresh before: Alive:%v Dead:%v", r.alivePlayerIds, r.deadPlayerIds)

	clear(r.alivePlayerIds)
	clear(r.deadPlayerIds)
	r.alivePlayerIds = r.alivePlayerIds[0:0]
	r.deadPlayerIds = r.deadPlayerIds[0:0]

	for i := range r.Players {
		if r.checkPlayerConnState(ctx, r.Players[i].Id) && r.Players[i].State == entity.PlayerAlive {
			r.alivePlayerIds = append(r.alivePlayerIds, r.Players[i].Id)
		} else {
			r.deadPlayerIds = append(r.deadPlayerIds, r.Players[i].Id)
		}
	}

	logrus.Debugf("refresh after: Alive:%v Dead:%v", r.alivePlayerIds, r.deadPlayerIds)
}

func (r *GameRunner) increasePlayerGlod(player *entity.Player) {
	player.MaxGold = utils.Min(player.MaxGold+1, int32(r.serverCfg.MaxPlayerGold))
	player.Gold = player.MaxGold
}

func (r *GameRunner) decreaseShopUpgradeCost(player *entity.Player) {
	shop := r.Shop[player.Id]
	shop.UpgradeCost = utils.Max(shop.UpgradeCost-1, 0)
}

func (r *GameRunner) refreshPlayerShop(ctx context.Context, player *entity.Player) error {
	shop := r.Shop[player.Id]

	shop.Retinue.Clear()

	res := r.cardPool.Random(ctx, shop.RetinueCap, r.shopCfg.Levels[shop.Level-1].Probabilities)
	for i, _ := range res {
		shop.Retinue.Add(&res[i])
	}

	logrus.Debugf("player:%v refresh shop:%v", player.Id, r.inspectShop(shop))

	return nil
}

func (r *GameRunner) clearAbattoirs() {
	for i, _ := range r.abattoirs {
		OM.PutAbattoir(r.abattoirs[i])
	}
	clear(r.abattoirs)
	r.abattoirs = r.abattoirs[0:0]
}

func (r *GameRunner) checkPlayerConnState(ctx context.Context, playerId uint64) bool {
	conn, err := dao.ConnectRepo.Find(ctx, playerId)
	if err != nil || conn == nil {
		return false
	}

	return true
}

func (r *GameRunner) FindPlayer(ctx context.Context, id uint64) *entity.Player {
	return r.Players[id]
}

func (r *GameRunner) Send(ctx context.Context, playerId uint64, cmd protocol.Command) error {
	_, err := r.sender.Send(ctx, playerId, cmd)
	return err
}

func (r *GameRunner) SyncPlayersGameData(ctx context.Context) {
	for _, player := range r.Players {
		if err := r.SyncGameData(ctx, player.Id); err != nil {
			logrus.Errorf("SyncGameData player:%d err: %v", player.Id, err)
		}
	}
}

func (r *GameRunner) SyncGameData(ctx context.Context, playerId uint64) error {
	conn, err := dao.ConnectRepo.Find(ctx, playerId)
	if err != nil {
		return err
	}
	if conn == nil {
		return nil
	}

	data := entity.GameData{
		Id:     r.Id,
		State:  r.Game.State,
		Round:  r.Round,
		Shop:   r.Shop[playerId],
		Player: r.Players[playerId],
	}

	return r.Send(ctx, playerId, protocol.New().GameDataSync(data))
}

func (r *GameRunner) GameOver(ctx context.Context) {
	res := entity.GameResult{
		Players: make([]entity.Player, 0, len(r.Players)),
	}
	for _, player := range r.Players {
		res.Players = append(res.Players, *player)
	}
	for _, player := range r.Players {
		if err := r.Send(ctx, player.Id, protocol.New().GameOver(res)); err != nil {
			logrus.Errorf("GameOver player:%d err: %v", player.Id, err)
		}
	}
}

func (r *GameRunner) CanOperate(ctx context.Context, player *entity.Player) bool {
	return player.State == entity.PlayerAlive && r.Game.State == entity.GameWaiting
}

func (r *GameRunner) UpgradeShop(ctx context.Context, player *entity.Player) error {
	if !service.CanUpgradeShop(ctx, r.serverCfg, r.Game.State, r.Shop[player.Id], player) {
		return nil
	}

	if err := r.upgradeShop(ctx, player); err != nil {
		return err
	}

	return r.SyncGameData(ctx, player.Id)
}

func (r *GameRunner) upgradeShop(ctx context.Context, player *entity.Player) error {
	shop := r.Shop[player.Id]
	logrus.Debugf("player:%v upgrade shop: %v -> %v, cost:%v", player.Id, shop.Level, shop.Level+1, shop.UpgradeCost)

	shop.Level += 1
	player.Gold -= int32(shop.UpgradeCost)

	shop.UpgradeCost = r.shopCfg.Levels[shop.Level-1].UpgradeCost
	shop.RetinueCap = r.shopCfg.Levels[shop.Level-1].RetinueCap

	return nil
}

func (r *GameRunner) RefreshShop(ctx context.Context, player *entity.Player) error {
	if !service.CanRefreshShop(ctx, r.serverCfg, r.Game.State, r.Shop[player.Id], player) {
		return nil
	}

	if err := r.refreshShop(ctx, player); err != nil {
		return err
	}

	return r.SyncGameData(ctx, player.Id)
}

func (r *GameRunner) refreshShop(ctx context.Context, player *entity.Player) error {
	if err := r.refreshPlayerShop(ctx, player); err != nil {
		return err
	}
	player.Gold -= 1
	return nil
}

func (r *GameRunner) BuyCard(ctx context.Context, player *entity.Player, idx int) error {
	if !service.CanBuyCard(ctx, r.serverCfg, r.Game.State, r.Shop[player.Id], player, idx) {
		return nil
	}

	if err := r.buyCard(ctx, player, idx); err != nil {
		return err
	}

	return r.SyncGameData(ctx, player.Id)
}

func (r *GameRunner) buyCard(ctx context.Context, player *entity.Player, idx int) error {
	card := r.Shop[player.Id].Retinue.Del(idx)
	player.CardList.Add(card)
	player.Gold -= 3
	logrus.Debugf("player:%v buy card: %v", player.Id, card.Inspect())
	return nil
}

func (r *GameRunner) SellCard(ctx context.Context, player *entity.Player, idx int) error {
	if !service.CanSellCard(ctx, r.serverCfg, r.Game.State, r.Shop[player.Id], player, idx) {
		return nil
	}

	if player.RetinueList.Del(idx) != nil {
		player.Gold += 1
		logrus.Debugf("player:%v sell card: %v", player.Id, idx)
	}

	return r.SyncGameData(ctx, player.Id)
}

func (r *GameRunner) PlaceRetinue(ctx context.Context, player *entity.Player, cardIdx, retinueIdx, targetIdx int) error {
	if !service.CanPlaceRetinue(ctx, r.serverCfg, r.Game.State, r.Shop[player.Id], player, cardIdx, retinueIdx, targetIdx) {
		return nil
	}

	if err := r.placeRetinue(ctx, player, cardIdx, retinueIdx, targetIdx); err != nil {
		return err
	}

	return r.SyncGameData(ctx, player.Id)
}

func (r *GameRunner) placeRetinue(ctx context.Context, player *entity.Player, cardIdx, retinueIdx, targetIdx int) error {
	retinue := player.CardList.Del(cardIdx)
	logrus.Debugf("player:%v place retinue: %v at %v -> target: %v", player.Id, retinue.Inspect(), retinueIdx, targetIdx)

	player.RetinueList.AddAt(retinueIdx, retinue)
	if retinueIdx <= targetIdx {
		targetIdx += 1
	}

	r.buffSys.OnPlace(ctx, retinue, player.RetinueList, retinueIdx, player.RetinueList, targetIdx)

	logrus.Debugf("player:%v place retinue: %v", player.Id, retinue.Inspect())
	return nil
}

func (r *GameRunner) DragRetinue(ctx context.Context, player *entity.Player, from, to int) error {
	if !service.CanDragRetinue(ctx, r.serverCfg, r.Game.State, r.Shop[player.Id], player, from, to) {
		return nil
	}

	player.RetinueList.Move(from, to)
	logrus.Debugf("player:%v move retinue: %v -> %v", player.Id, from, to)

	return r.SyncGameData(ctx, player.Id)
}

func (r *GameRunner) roundDuration() time.Duration {
	return time.Duration(r.Round*r.serverCfg.GameRoundAdditionalDuration+r.serverCfg.GameRoundBaseDuration) * time.Second
}

func (r *GameRunner) inspectShop(shop *entity.Shop) string {
	return fmt.Sprintf("{%d %d [%s]}", shop.Level, shop.UpgradeCost, r.inspectRetinues(shop.Retinue))
}

func (r *GameRunner) inspectRetinues(l *list.DoubleLinkList[entity.Retinue]) string {
	var buf strings.Builder
	l.Each(func(r *entity.Retinue) bool {
		buf.WriteString(r.Inspect())
		return true
	})
	return buf.String()
}
