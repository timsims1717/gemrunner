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
	TileBGBatch = "tile_bg"
	TileFGBatch = "tile_fg"
	UIBatch     = "ui"
	UITile      = "ui_tile"
)

var (
	BlackColor = color.RGBA{
		R: 19,
		G: 19,
		B: 19,
		A: 255,
	}
)
