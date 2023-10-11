package constants

import "image/color"

const (
	Title   = "Gem Runner"
	Release = 0
	Version = 1
	Build   = 20230125

	// World Constants
	TileSize     = 16.
	PuzzleWidth  = 28
	PuzzleHeight = 16

	// Editor Constants
	BlockSelectWidth = 6.

	// Batches
	BGBatch = "tile_bg"
	FGBatch = "tile_fg"
	UIBatch = "ui"
)

var (
	BlackColor = color.RGBA{
		R: 19,
		G: 19,
		B: 19,
		A: 255,
	}
	WhiteColor = color.RGBA{
		R: 245,
		G: 245,
		B: 245,
		A: 255,
	}
)
