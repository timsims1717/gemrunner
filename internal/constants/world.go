package constants

import "image/color"

// Tile Set Names
const (
	TileSetRock   = "rock"
	TileSetSlate  = "slate"
	TileSetBrick  = "brick"
	TileSetGravel = "gravel"
	TileSetDirt   = "dirt"
	TileSetStone  = "stone"
	TileSetShell  = "shell"
	TileSetMetal  = "metal"
	TileSetIce    = "ice"
	TileSetTree   = "tree"
)

const (
	WorldMoss = iota
	WorldJungle
	WorldBrick
	WorldBeam
	WorldIce
	WorldSnow
	WorldDark
	WorldDungeon
	WorldBeach
	WorldRedRock
	WorldSlime
	WorldFungus
	WorldSandstone
	WorldDunes
	WorldLava
	WorldBasalt
	WorldAbyss
	WorldReef
	WorldGravelPit
	WorldSpire
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
		WorldMoss:      TileSetRock,
		WorldJungle:    TileSetTree,
		WorldBrick:     TileSetBrick,
		WorldBeam:      TileSetMetal,
		WorldIce:       TileSetIce,
		WorldSnow:      TileSetStone,
		WorldDark:      TileSetGravel,
		WorldDungeon:   TileSetBrick,
		WorldBeach:     TileSetShell,
		WorldRedRock:   TileSetRock,
		WorldSlime:     TileSetTree,
		WorldFungus:    TileSetDirt,
		WorldSandstone: TileSetSlate,
		WorldDunes:     TileSetDirt,
		WorldLava:      TileSetIce,
		WorldBasalt:    TileSetStone,
		WorldAbyss:     TileSetSlate,
		WorldReef:      TileSetShell,
		WorldGravelPit: TileSetGravel,
		WorldSpire:     TileSetMetal,
	}

	WorldPrimary = map[int]color.RGBA{
		WorldMoss:      ColorGray,
		WorldJungle:    ColorGreen,
		WorldBrick:     ColorRed,
		WorldBeam:      ColorPurple,
		WorldIce:       ColorCyan,
		WorldSnow:      ColorBlue,
		WorldDark:      ColorBlue,
		WorldDungeon:   ColorWhite,
		WorldBeach:     ColorYellow,
		WorldRedRock:   ColorRed,
		WorldSlime:     ColorBrown,
		WorldFungus:    ColorPurple,
		WorldSandstone: ColorTan,
		WorldDunes:     ColorGreen,
		WorldLava:      ColorOrange,
		WorldBasalt:    ColorGray,
		WorldAbyss:     ColorPurple,
		WorldReef:      ColorWhite,
		WorldGravelPit: ColorGray,
		WorldSpire:     ColorWhite,
	}

	WorldSecondary = map[int]color.RGBA{
		WorldMoss:      ColorGreen,
		WorldJungle:    ColorBrown,
		WorldBrick:     ColorBlue,
		WorldBeam:      ColorRed,
		WorldIce:       ColorWhite,
		WorldSnow:      ColorCyan,
		WorldDark:      ColorGray,
		WorldDungeon:   ColorBlue,
		WorldBeach:     ColorBrown,
		WorldRedRock:   ColorYellow,
		WorldSlime:     ColorTan,
		WorldFungus:    ColorBrown,
		WorldSandstone: ColorOrange,
		WorldDunes:     ColorTan,
		WorldLava:      ColorRed,
		WorldBasalt:    ColorOrange,
		WorldAbyss:     ColorBlue,
		WorldReef:      ColorPurple,
		WorldGravelPit: ColorYellow,
		WorldSpire:     ColorGray,
	}

	WorldMusic = map[int]string{
		WorldMoss:      TrackJungle,
		WorldJungle:    TrackJungle,
		WorldBrick:     TrackUrban,
		WorldBeam:      TrackUrban,
		WorldIce:       TrackIce,
		WorldSnow:      TrackIce,
		WorldDark:      TrackDark,
		WorldDungeon:   TrackDark,
		WorldBeach:     TrackBeach,
		WorldRedRock:   TrackBeach,
		WorldSlime:     TrackFungus,
		WorldFungus:    TrackFungus,
		WorldSandstone: TrackDesert,
		WorldDunes:     TrackDesert,
		WorldLava:      TrackLava,
		WorldBasalt:    TrackLava,
		WorldAbyss:     TrackReef,
		WorldReef:      TrackReef,
		WorldGravelPit: TrackMech,
		WorldSpire:     TrackMech,
	}
)
