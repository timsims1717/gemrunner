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
	TileSetBubble  = "bubble"
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
	WorldBubblegum
	WorldCrystal
	WorldAbyss
	WorldReef
	WorldGravelPit
	WorldSpire
	WorldSandstone
	WorldDunes
	WorldDark
	WorldDungeon
	WorldGilded
	WorldIvy
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
	ColorBurnt = color.RGBA{
		R: 193,
		G: 68,
		B: 27,
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
		WorldBubblegum: "Bubblegum",
		WorldCrystal:   "Crystal",
		WorldAbyss:     "Abyss",
		WorldReef:      "Reef",
		WorldGravelPit: "Gravel Pit",
		WorldSpire:     "Spire",
		WorldSandstone: "Sandstone",
		WorldDunes:     "Dunes",
		WorldDark:      "Dark",
		WorldDungeon:   "Dungeon",
		WorldGilded:    "Gilded",
		WorldIvy:       "Ivy",
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
		WorldBubblegum: TileSetBubble,
		WorldCrystal:   TileSetCrystal,
		WorldAbyss:     TileSetSlate,
		WorldReef:      TileSetShell,
		WorldGravelPit: TileSetGravel,
		WorldSpire:     TileSetMetal,
		WorldSandstone: TileSetSlate,
		WorldDunes:     TileSetDirt,
		WorldDark:      TileSetGravel,
		WorldDungeon:   TileSetBrick,
		WorldGilded:    TileSetBubble,
		WorldIvy:       TileSetCrystal,
		WorldLava:      TileSetIce,
		WorldBasalt:    TileSetStone,
	}

	WorldPrimary = map[int]color.RGBA{
		WorldMoss:      ColorGray,
		WorldJungle:    ColorGreen,
		WorldBrick:     ColorRed,
		WorldBeam:      ColorPurple,
		WorldSlime:     ColorBrown,
		WorldFungus:    ColorPurple,
		WorldBeach:     ColorYellow,
		WorldRedRock:   ColorBurnt,
		WorldIce:       ColorCyan,
		WorldSnow:      ColorBlue,
		WorldBubblegum: ColorPink,
		WorldCrystal:   ColorPurple,
		WorldAbyss:     ColorPurple,
		WorldReef:      ColorTan,
		WorldGravelPit: ColorGray,
		WorldSpire:     ColorGold,
		WorldSandstone: ColorTan,
		WorldDunes:     ColorGreen,
		WorldDark:      ColorBlue,
		WorldDungeon:   ColorGray,
		WorldGilded:    ColorGold,
		WorldIvy:       ColorGreen,
		WorldLava:      ColorOrange,
		WorldBasalt:    ColorGray,
	}

	WorldSecondary = map[int]color.RGBA{
		WorldMoss:      ColorGreen,
		WorldJungle:    ColorBrown,
		WorldBrick:     ColorBlue,
		WorldBeam:      ColorRed,
		WorldSlime:     ColorTan,
		WorldFungus:    ColorBrown,
		WorldBeach:     ColorBrown,
		WorldRedRock:   ColorYellow,
		WorldIce:       ColorLightGray,
		WorldSnow:      ColorCyan,
		WorldBubblegum: ColorRed,
		WorldCrystal:   ColorPink,
		WorldAbyss:     ColorBlue,
		WorldReef:      ColorPurple,
		WorldGravelPit: ColorRed,
		WorldSpire:     ColorGray,
		WorldSandstone: ColorOrange,
		WorldDunes:     ColorTan,
		WorldDark:      ColorGray,
		WorldDungeon:   ColorBlue,
		WorldGilded:    ColorGray,
		WorldIvy:       ColorGold,
		WorldLava:      ColorRed,
		WorldBasalt:    ColorOrange,
	}

	WorldDoodad = map[int]color.RGBA{
		WorldMoss:      ColorGreen,
		WorldJungle:    ColorBrown,
		WorldBrick:     ColorBlue,
		WorldBeam:      ColorRed,
		WorldSlime:     ColorTan,
		WorldFungus:    ColorRed,
		WorldBeach:     ColorGreen,
		WorldRedRock:   ColorGreen,
		WorldIce:       ColorLightGray,
		WorldSnow:      ColorLightGray,
		WorldBubblegum: ColorPink,
		WorldCrystal:   ColorBlue,
		WorldAbyss:     ColorBlue,
		WorldReef:      ColorPurple,
		WorldGravelPit: ColorRed,
		WorldSpire:     ColorGray,
		WorldSandstone: ColorGreen,
		WorldDunes:     ColorGreen,
		WorldDark:      ColorBlue,
		WorldDungeon:   ColorBlue,
		WorldGilded:    ColorGold,
		WorldIvy:       ColorGold,
		WorldLava:      ColorRed,
		WorldBasalt:    ColorOrange,
	}

	WorldMusic = map[int]string{
		WorldMoss:      TrackJungle,
		WorldJungle:    TrackJungle,
		WorldBrick:     TrackUrban,
		WorldBeam:      TrackUrban,
		WorldSlime:     TrackFungus,
		WorldFungus:    TrackFungus,
		WorldBeach:     TrackBeach,
		WorldRedRock:   TrackBeach,
		WorldIce:       TrackIce,
		WorldSnow:      TrackIce,
		WorldBubblegum: TrackCrystal,
		WorldCrystal:   TrackCrystal,
		WorldAbyss:     TrackReef,
		WorldReef:      TrackReef,
		WorldGravelPit: TrackMech,
		WorldSpire:     TrackMech,
		WorldSandstone: TrackDesert,
		WorldDunes:     TrackDesert,
		WorldDark:      TrackDark,
		WorldDungeon:   TrackDark,
		WorldGilded:    TrackGilded,
		WorldIvy:       TrackGilded,
		WorldLava:      TrackLava,
		WorldBasalt:    TrackLava,
	}

	WorldDoodads = map[int]string{
		WorldMoss:      DoodadReeds,
		WorldJungle:    DoodadNest,
		WorldBrick:     DoodadReeds,
		WorldBeam:      DoodadChain,
		WorldSlime:     DoodadMush1,
		WorldFungus:    DoodadMush1,
		WorldBeach:     DoodadReeds,
		WorldRedRock:   DoodadCactus2,
		WorldIce:       DoodadBoulder,
		WorldSnow:      DoodadSkull,
		WorldBubblegum: DoodadBoulder,
		WorldCrystal:   DoodadBoulder,
		WorldAbyss:     DoodadVine,
		WorldReef:      DoodadVine,
		WorldGravelPit: DoodadChain,
		WorldSpire:     DoodadGear,
		WorldSandstone: DoodadNest,
		WorldDunes:     DoodadSkull,
		WorldDark:      DoodadReeds,
		WorldDungeon:   DoodadSkull,
		WorldGilded:    DoodadGear,
		WorldIvy:       DoodadReeds,
		WorldLava:      DoodadBoulder,
		WorldBasalt:    DoodadReeds,
	}
)
