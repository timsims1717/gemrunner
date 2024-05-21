package data

type PlayerStats struct {
	Score  int `json:"score"`
	Deaths int `json:"deaths"`
	Gems   int `json:"gems"`
	LScore int `json:"-"`
	LGems  int `json:"-"`
}

func NewStats() *PlayerStats {
	return &PlayerStats{}
}
