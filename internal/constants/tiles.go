package constants

import "image/color"

// Editor Tiles
const (
	TileLadderBottom    = "ladder_bottom"
	TileLadderMiddle    = "ladder_middle"
	TileLadderTop       = "ladder_top"
	TileLadderCrackB    = "ladder_cracked_bottom"
	TileLadderCrackM    = "ladder_cracked_middle"
	TileExitLadderB     = "exit_ladder_bottom"
	TileExitLadderM     = "exit_ladder_middle"
	TileExitLadderT     = "exit_ladder_top"
	TileLadderCrackingM = "ladder_cracking_middle"
	TileLadderCrackingB = "ladder_cracking_bottom"

	TileLadderLedgeBottom    = "ladder_ledge_bottom"
	TileLadderLedgeMiddle    = "ladder_ledge_middle"
	TileLadderLedgeCrackB    = "ladder_ledge_cracked_bottom"
	TileLadderLedgeCrackM    = "ladder_ledge_cracked_middle"
	TileLadderLedgeCrackingM = "ladder_ledge_cracking_middle"
	TileLadderLedgeCrackingB = "ladder_ledge_cracking_bottom"

	TileTop       = "_top"
	TileBottom    = "_bottom"
	TileBottomTop = "_bottom_top"
	TileAlt       = "_alt"

	StrColorGreen  = "green"
	StrColorBrown  = "brown"
	StrColorGray   = "gray"
	StrColorCyan   = "cyan"
	StrColorBlue   = "blue"
	StrColorOrange = "orange"
	StrColorPurple = "purple"
	StrColorYellow = "yellow"
	TileDoor       = "door_"
	TileLock       = "lock_"
	TileUnlock     = "unlock_"
	TileOpen       = "_open"

	TileDoorYellow = "door_yellow"
	TileLockYellow = "lock_yellow"
	TileDoorOrange = "door_orange"
	TileLockOrange = "lock_orange"
	TileDoorGray   = "door_gray"
	TileLockGray   = "lock_gray"
	TileDoorCyan   = "door_cyan"
	TileLockCyan   = "lock_cyan"
	TileDoorBlue   = "door_blue"
	TileLockBlue   = "lock_blue"
	TileDoorGreen  = "door_green"
	TileLockGreen  = "lock_green"
	TileDoorPurple = "door_purple"
	TileLockPurple = "lock_purple"
	TileDoorBrown  = "door_brown"
	TileLockBrown  = "lock_brown"

	TileDemonRegen = "demon_regen_tile"

	TileFall        = "fall"
	TilePhase       = "phase"
	TileSpike       = "spike"
	TileCracked     = "cracked"
	TileCrackedShow = "cracked_show"
	TileCracking    = "cracking"
	TileTurf        = "turf"
	TileEmpty       = "empty"
)

// Characters
const (
	CharPlayer1 = "player1"
	CharPlayer2 = "player2"
	CharPlayer3 = "player3"
	CharPlayer4 = "player4"
	CharDemon   = "demon"
	CharFly     = "fly"
)

// Items
const (
	ItemKeyYellow = "key_yellow"
	ItemKeyOrange = "key_orange"
	ItemKeyGray   = "key_gray"
	ItemKeyCyan   = "key_cyan"
	ItemKeyBlue   = "key_blue"
	ItemKeyGreen  = "key_green"
	ItemKeyPurple = "key_purple"
	ItemKeyBrown  = "key_brown"
	ItemBox       = "box"
	ItemBoxPiece  = "box_piece"
	ItemGemYellow = "gem_yellow"
	ItemGemOrange = "gem_orange"
	ItemGemGray   = "gem_gray"
	ItemGemCyan   = "gem_cyan"
	ItemGemBlue   = "gem_blue"
	ItemGemGreen  = "gem_green"
	ItemGemPurple = "gem_purple"
	ItemGemBrown  = "gem_brown"
)

// Doodads
const (
	DoodadChain   = "chain"
	DoodadReeds   = "reeds"
	DoodadFlowers = "flowers"
)

// World Names
const (
	WorldRock = iota
	WorldSlate
	WorldBrick
	WorldGravel
	WorldDirt
	WorldStone
	WorldShell
	WorldMetal
	WorldIce
	WorldTree
	WorldCustom
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
	ColorCyan = color.RGBA{
		R: 26,
		G: 202,
		B: 202,
		A: 255,
	}
	ColorTan = color.RGBA{
		R: 255,
		G: 204,
		B: 149,
		A: 255,
	}

	Player1Color = ColorBlue
	Player2Color = ColorOrange
	Player3Color = ColorPurple
	Player4Color = ColorYellow

	WorldSprites = map[int]string{
		WorldRock:   "rock",
		WorldSlate:  "slate",
		WorldBrick:  "brick",
		WorldGravel: "gravel",
		WorldDirt:   "dirt",
		WorldStone:  "stone",
		WorldShell:  "shell",
		WorldMetal:  "metal",
		WorldIce:    "ice",
		WorldTree:   "tree",
		WorldCustom: "custom",
	}

	WorldPrimary = map[int]color.RGBA{
		WorldRock:   ColorGray,
		WorldSlate:  ColorOrange,
		WorldBrick:  ColorRed,
		WorldGravel: ColorBlue,
		WorldDirt:   ColorBrown,
		WorldStone:  ColorGreen,
		WorldShell:  ColorYellow,
		WorldMetal:  ColorPurple,
		WorldIce:    ColorCyan,
		WorldTree:   ColorGreen,
	}

	WorldSecondary = map[int]color.RGBA{
		WorldRock:   ColorGreen,
		WorldSlate:  ColorTan,
		WorldBrick:  ColorBlue,
		WorldGravel: ColorGray,
		WorldDirt:   ColorGray,
		WorldStone:  ColorOrange,
		WorldShell:  ColorBrown,
		WorldMetal:  ColorRed,
		WorldIce:    ColorWhite,
		WorldTree:   ColorBrown,
	}
)
