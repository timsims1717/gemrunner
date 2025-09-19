package data

import "gemrunner/internal/constants"

type PlayerStats struct {
	Score  int `json:"score"`
	Deaths int `json:"deaths"`
	Gems   int `json:"gems"`
	Bombs  int `json:"bombs"`
	LScore int `json:"-"`
	LGems  int `json:"-"`
	LBombs int `json:"-"`
}

func NewStats() *PlayerStats {
	return &PlayerStats{
		Bombs: constants.SmallBombInv,
	}
}
