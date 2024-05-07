package constants

import "image/color"

// Tile Set Names
const (
	TileSetRock    = "rock"
	TileSetSlate   = "slate"
	TileSetBrick   = "brick"
	TileSetGravel  = "gravel"
	TileSetDirt    = "dirt"
	TileSetStone   = "stone"
	TileSetShell   = "shell"
	TileSetMetal   = "metal"
	TileSetIce     = "ice"
	TileSetTree    = "tree"
	TileSetCrystal = "crystal"
)

const (
	WorldMoss = iota
	WorldJungle
	WorldBrick
	WorldBeam
	WorldSlime
	WorldFungus
	WorldBeach
	WorldRedRock
	WorldIce
	WorldSnow
	WorldOtherPink
	WorldCrystal
	WorldAbyss
	WorldReef
	WorldSandstone
	WorldDunes
	WorldDark
	WorldDungeon
	WorldGravelPit
	WorldSpire
	WorldLava
	WorldBasalt
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
	ColorLightGray = color.RGBA{
		R: 210,
		G: 210,
		B: 210,
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
	ColorGold = color.RGBA{
		R: 197,
		G: 198,
		B: 93,
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
	ColorPink = color.RGBA{
		R: 230,
		G: 35,
		B: 165,
		A: 255,
	}

	Player1Color = ColorBlue
	Player2Color = ColorOrange
	Player3Color = ColorPurple
	Player4Color = ColorGold

	WorldNames = map[int]string{
		WorldMoss:      "Moss",
		WorldJungle:    "Jungle",
		WorldBrick:     "Brick",
		WorldBeam:      "Beams",
		WorldSlime:     "Shrooms",
		WorldFungus:    "Fungus",
		WorldBeach:     "Beach",
		WorldRedRock:   "Red Rock",
		WorldIce:       "Ice",
		WorldSnow:      "Snow",
		WorldOtherPink: "Other Pink",
		WorldCrystal:   "Crystal",
		WorldAbyss:     "Abyss",
		WorldReef:      "Reef",
		WorldSandstone: "Sandstone",
		WorldDunes:     "Dunes",
		WorldDark:      "Dark",
		WorldDungeon:   "Dungeon",
		WorldGravelPit: "Gravel Pit",
		WorldSpire:     "Spire",
		WorldLava:      "Lava",
		WorldBasalt:    "Basalt",
		WorldCustom:    "Custom",
	}

	WorldSprites = map[int]string{
		WorldMoss:      TileSetRock,
		WorldJungle:    TileSetTree,
		WorldBrick:     TileSetBrick,
		WorldBeam:      TileSetMetal,
		WorldSlime:     TileSetTree,
		WorldFungus:    TileSetDirt,
		WorldBeach:     TileSetShell,
		WorldRedRock:   TileSetStone,
		WorldIce:       TileSetIce,
		WorldSnow:      TileSetRock,
		WorldOtherPink: TileSetRock,
		WorldCrystal:   TileSetCrystal,
		WorldAbyss:     TileSetSlate,
		WorldReef:      TileSetShell,
		WorldSandstone: TileSetSlate,
		WorldDunes:     TileSetDirt,
		WorldDark:      TileSetGravel,
		WorldDungeon:   TileSetBrick,
		WorldGravelPit: TileSetGravel,
		WorldSpire:     TileSetMetal,
		WorldLava:      TileSetIce,
		WorldBasalt:    TileSetStone,
	}

	WorldPrimary = map[int]color.RGBA{
		WorldMoss:      ColorGray,
		WorldJungle:    ColorGreen,
		WorldBrick:     ColorRed,
		WorldBeam:      ColorPurple,
		WorldIce:       ColorCyan,
		WorldSnow:      ColorBlue,
		WorldOtherPink: ColorPink,
		WorldCrystal:   ColorPurple,
		WorldDark:      ColorBlue,
		WorldDungeon:   ColorGray,
		WorldBeach:     ColorYellow,
		WorldRedRock:   ColorTan,
		WorldSlime:     ColorBrown,
		WorldFungus:    ColorPurple,
		WorldSandstone: ColorTan,
		WorldDunes:     ColorGreen,
		WorldLava:      ColorOrange,
		WorldBasalt:    ColorGray,
		WorldAbyss:     ColorPurple,
		WorldReef:      ColorTan,
		WorldGravelPit: ColorGray,
		WorldSpire:     ColorGold,
	}

	WorldSecondary = map[int]color.RGBA{
		WorldMoss:      ColorGreen,
		WorldJungle:    ColorBrown,
		WorldBrick:     ColorBlue,
		WorldBeam:      ColorRed,
		WorldIce:       ColorLightGray,
		WorldSnow:      ColorCyan,
		WorldOtherPink: ColorRed,
		WorldCrystal:   ColorPink,
		WorldDark:      ColorGray,
		WorldDungeon:   ColorBlue,
		WorldBeach:     ColorBrown,
		WorldRedRock:   ColorGold,
		WorldSlime:     ColorTan,
		WorldFungus:    ColorBrown,
		WorldSandstone: ColorOrange,
		WorldDunes:     ColorTan,
		WorldLava:      ColorRed,
		WorldBasalt:    ColorOrange,
		WorldAbyss:     ColorBlue,
		WorldReef:      ColorPurple,
		WorldGravelPit: ColorRed,
		WorldSpire:     ColorGray,
	}

	WorldDoodad = map[int]color.RGBA{
		WorldMoss:      ColorGreen,
		WorldJungle:    ColorBrown,
		WorldBrick:     ColorBlue,
		WorldBeam:      ColorRed,
		WorldIce:       ColorLightGray,
		WorldSnow:      ColorLightGray,
		WorldOtherPink: ColorPink,
		WorldCrystal:   ColorBlue,
		WorldDark:      ColorBlue,
		WorldDungeon:   ColorBlue,
		WorldBeach:     ColorGreen,
		WorldRedRock:   ColorGreen,
		WorldSlime:     ColorTan,
		WorldFungus:    ColorRed,
		WorldSandstone: ColorGreen,
		WorldDunes:     ColorGreen,
		WorldLava:      ColorRed,
		WorldBasalt:    ColorOrange,
		WorldAbyss:     ColorBlue,
		WorldReef:      ColorPurple,
		WorldGravelPit: ColorRed,
		WorldSpire:     ColorGray,
	}

	WorldMusic = map[int]string{
		WorldMoss:      TrackJungle,
		WorldJungle:    TrackJungle,
		WorldBrick:     TrackUrban,
		WorldBeam:      TrackUrban,
		WorldIce:       TrackIce,
		WorldSnow:      TrackIce,
		WorldOtherPink: TrackJungle,
		WorldCrystal:   TrackJungle,
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

	WorldDoodads = map[int]string{
		WorldMoss:      DoodadReeds,
		WorldJungle:    DoodadNest,
		WorldBrick:     DoodadReeds,
		WorldBeam:      DoodadChain,
		WorldIce:       DoodadBoulder,
		WorldSnow:      DoodadSkull,
		WorldOtherPink: DoodadBoulder,
		WorldCrystal:   DoodadBoulder,
		WorldDark:      DoodadReeds,
		WorldDungeon:   DoodadSkull,
		WorldBeach:     DoodadReeds,
		WorldRedRock:   DoodadCactus2,
		WorldSlime:     DoodadMush1,
		WorldFungus:    DoodadMush1,
		WorldSandstone: DoodadNest,
		WorldDunes:     DoodadSkull,
		WorldLava:      DoodadBoulder,
		WorldBasalt:    DoodadReeds,
		WorldAbyss:     DoodadVine,
		WorldReef:      DoodadVine,
		WorldGravelPit: DoodadChain,
		WorldSpire:     DoodadGear,
	}
)
