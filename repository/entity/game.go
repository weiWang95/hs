package entity

type GameState int

const (
	GamePending GameState = iota
	GameWaiting
	GamePlaying
	GameFinished
)

type Game struct {
	Id    uint64    `json:"id"`
	State GameState `json:"state"`
	Round int       `json:"round"`

	Players map[uint64]*Player `json:"players"`
	Shop    map[uint64]*Shop   `json:"shop"`
}

type GameData struct {
	Id    uint64    `json:"id"`
	State GameState `json:"state"`
	Round int       `json:"round"`

	Player *Player `json:"player"`
	Shop   *Shop   `json:"shop"`
}

type GameResult struct {
	Players []Player `json:"players"`
}
