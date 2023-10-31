package constants

import "image/color"

// Editor Tiles
const (
	TileLadderBottom = "ladder_bottom"
	TileLadderMiddle = "ladder_middle"
	TileLadderTop    = "ladder_top"
	TileTop          = "_top"
	TileBottom       = "_bottom"
	TileBottomTop    = "_bottom_top"
	TileAlt          = "_alt"
	TileDoorPink     = "door_pink"
	TileLockPink     = "lock_pink"
	TileDoorBlue     = "door_blue"
	TileLockBlue     = "lock_blue"
	TileFall         = "fall"
	TileTurf         = "turf"
	TileEmpty        = "empty"
)

// Characters
const (
	CharPlayer1 = "player_1"
	CharDevil   = "devil"
)

// Items
const (
	ItemKeyPink = "key_pink"
	ItemKeyBlue = "key_blue"
	ItemBox     = "box"
	ItemGem     = "gem"
)

// Doodads
const (
	DoodadChain   = "chain"
	DoodadReeds   = "reeds"
	DoodadFlowers = "flowers"
)

// World Names
const (
	WorldRock      = "rock"
	WorldSlate     = "slate"
	WorldBrick     = "brick"
	WorldGravel    = "gravel"
	WorldDirt      = "dirt"
	WorldSolidRock = "solid_stone"
	WorldShell     = "shell"
	WorldMetal     = "metal"
	WorldCustom    = "custom"
)

var (
	ColorBlack = color.RGBA{
		R: 19,
		G: 19,
		B: 19,
		A: 255,
	}
	ColorWhite = color.RGBA{
		R: 245,
		G: 245,
		B: 245,
		A: 255,
	}
	ColorRed = color.RGBA{
		R: 255,
		G: 77,
		B: 77,
		A: 255,
	}
	ColorBlue = color.RGBA{
		R: 75,
		G: 122,
		B: 255,
		A: 255,
	}
	ColorGreen = color.RGBA{
		R: 94,
		G: 143,
		B: 86,
		A: 255,
	}
	ColorOrange = color.RGBA{
		R: 255,
		G: 149,
		B: 75,
		A: 255,
	}
	ColorPurple = color.RGBA{
		R: 163,
		G: 73,
		B: 177,
		A: 255,
	}
	ColorYellow = color.RGBA{
		R: 255,
		G: 213,
		B: 9,
		A: 255,
	}
	ColorGray = color.RGBA{
		R: 91,
		G: 91,
		B: 91,
		A: 255,
	}
	ColorBrown = color.RGBA{
		R: 121,
		G: 95,
		B: 67,
		A: 255,
	}
)
