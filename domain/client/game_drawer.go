package client

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"hs/repository/entity"
)

type GameDrawer interface {
	Draw(ctx context.Context, data *GameData) error
	ShowHelp(ctx context.Context)
	ShowGameResult(ctx context.Context, data *entity.GameResult)
}

func NewCmdDrawer() GameDrawer {
	return &cmdDrawer{}
}

type cmdDrawer struct {
}

func (c *cmdDrawer) Draw(ctx context.Context, data *GameData) error {
	c.clear()

	var buf strings.Builder
	buf.WriteRune('\n')
	buf.WriteRune('\n')
	buf.WriteString(fmt.Sprintf("ID:%d State:%d Round:%d\n", data.Id, data.State, data.Round))
	buf.WriteString(fmt.Sprintf("Shop %v: UpgradeCost: %d\n", c.levelChar(data.Shop.Level), data.Shop.UpgradeCost))

	buf.WriteRune('\n')
	// var item entity.Retinue
	if data.Shop.Retinue != nil {
		data.Shop.Retinue.Each(func(item *entity.Retinue) bool {
			buf.WriteString(fmt.Sprintf(" [%s %d-%d (%d-%d)] ", c.levelChar(item.Level), item.FinalAttack, item.FinalHp, item.Attack, item.Hp))
			return true
		})
	}
	buf.WriteString("\n\n----------------------------------------------\n\n")

	if data.Player.RetinueList != nil {
		data.Player.RetinueList.Each(func(item *entity.Retinue) bool {
			buf.WriteString(fmt.Sprintf(" [%s %d-%d (%d-%d)] ", c.levelChar(item.Level), item.FinalAttack, item.FinalHp, item.Attack, item.Hp))
			return true
		})
	}
	buf.WriteRune('\n')
	buf.WriteRune('\n')

	buf.WriteString(fmt.Sprintf("Player[%d]: Hp:%d Gold:%d/%d\n", data.Player.Id, data.Player.Hp, data.Player.Gold, data.Player.MaxGold))
	if data.Player.CardList != nil {
		data.Player.CardList.Each(func(item *entity.Retinue) bool {
			buf.WriteString(fmt.Sprintf(" [%s %d-%d (%d-%d)] ", c.levelChar(item.Level), item.FinalAttack, item.FinalHp, item.Attack, item.Hp))
			return true
		})
	}
	buf.WriteRune('\n')

	fmt.Print(buf.String())
	return nil
}

func (c *cmdDrawer) ShowHelp(ctx context.Context) {
	var buf strings.Builder

	buf.WriteString("购买: buy|b {index}\n")
	buf.WriteString("卖出: sell|s {index}\n")
	buf.WriteString("使用: use|u {from} {to}\n")
	buf.WriteString("升级: upgrade|up\n")
	buf.WriteString("刷新: refresh|r\n")
	buf.WriteString("帮助: help|h\n")

	fmt.Print(buf.String())
}

func (c *cmdDrawer) ShowGameResult(ctx context.Context, data *entity.GameResult) {
	var buf strings.Builder
	buf.WriteString("Game Result:\n")

	for _, player := range data.Players {
		buf.WriteString(fmt.Sprintf("Player[%d-%s]: Hp:%d\n", player.Id, player.Nickname, player.Hp))
	}

	fmt.Print(buf.String())
}

func (c *cmdDrawer) levelChar(level int) string {
	switch level {
	case 1:
		return "①"
	case 2:
		return "②"
	case 3:
		return "③"
	case 4:
		return "④"
	case 5:
		return "⑤"
	case 6:
		return "⑥"
	}
	return "?"
}

func (c *cmdDrawer) clear() {
	// cmd := "clear"
	// if runtime.GOOS == "windows" {
	// 	cmd = "cls"
	// }

	// _ = exec.Command("bash", "-c", cmd).Run()
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
