package data

type PlayerStats struct {
	Score int
}

func NewStats() *PlayerStats {
	return &PlayerStats{}
}
