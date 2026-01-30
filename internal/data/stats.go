package data

import "gemrunner/internal/constants"

type PlayerStats struct {
	Score     int `json:"score"`
	Deaths    int `json:"deaths"`
	Gems      int `json:"gems"`
	MaxBombs  int `json:"bombs"`
	CurrBombs int `json:"currBombs"`
	LScore    int `json:"-"`
	LGems     int `json:"-"`
	LBombs    int `json:"-"`

	Inventory *BasicItem `json:"inventory"`
}

func NewStats() *PlayerStats {
	return &PlayerStats{
		MaxBombs: constants.SmallBombInv,
	}
}
