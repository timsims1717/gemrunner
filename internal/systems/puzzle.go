package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
)

func CreatePlayArea() *data.PlayArea {
	fp := &data.PlayArea{
		IMDraw: imdraw.New(nil),
	}
	InitBorder(fp)
	return fp
}

func SetPuzzle(playArea *data.PlayArea, puzzle *data.Puzzle) {
	if playArea == nil {
		panic("SetPuzzle: puzzle view is nil")
	} else if puzzle == nil {
		panic("SetPuzzle: puzzle is nil")
	}
	DisposePuzzle(playArea.Puzzle)
	VerifyPuzzle(puzzle)
	playArea.Puzzle = puzzle
}

func InitPuzzle(playArea *data.PlayArea) {
	if playArea == nil {
		panic("InitPuzzle: puzzle view is nil")
	} else if playArea.Puzzle == nil {
		panic("InitPuzzle: puzzle is nil")
	}
	for y, row := range playArea.Puzzle.Tiles.T {
		for x, tile := range row {
			tile.Coords = world.NewCoords(x, y)
			obj := object.New()
			obj.Pos = world.MapToWorld(tile.Coords)
			obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
			obj.Flip = tile.Metadata.Flipped
			obj.Layer = 2
			e := myecs.Manager.NewEntity().
				AddComponent(myecs.Object, obj).
				AddComponent(myecs.Tile, tile)
			tile.Object = obj
			tile.Entity = e
			if tile.TextData != nil {
				data.CreateFloatingText(tile, tile.TextData)
			}
		}
	}
	UpdateEditorOptions()
	InitPlayArea(playArea)
	UpdateEditorShaders(playArea.Puzzle)
	data.CurrPuzzleSet.CurrPuzzle.Update = true
}

func InitPlayArea(playArea *data.PlayArea) {
	InitPlayAreaView(playArea)
	//UpdateBackgroundShaders(playArea)
	//UpdatePuzzleShaders(playArea)
	ChangeFullWorldShader(playArea, playArea.Puzzle.Metadata.ShaderMode)
	UpdatePlayAreaView(playArea)
}

func InitPlayAreaView(fp *data.PlayArea) {
	if fp == nil {
		panic("InitPlayAreaView: puzzle view is nil")
	}
	w, h := float64(fp.Puzzle.Metadata.Width), float64(fp.Puzzle.Metadata.Height)
	InitBorder(fp)
	SetBorder(fp)
	wRatio := viewport.MainCamera.Rect.W() / (constants.PuzzleWidth * world.TileSize)
	hRatio := viewport.MainCamera.Rect.H() / (constants.PuzzleHeight * world.TileSize)
	maxRatio := wRatio
	if hRatio < wRatio {
		maxRatio = hRatio
	}
	maxRatio *= constants.ScreenRatioLimit

	constants.PickedRatio = 1.
	for constants.PickedRatio+1 < maxRatio {
		constants.PickedRatio += 1
	}

	if fp.WorldView == nil {
		fp.WorldView = viewport.New(nil)
		xWidth := w * world.TileSize * constants.PickedRatio
		yHeight := h * world.TileSize * constants.PickedRatio
		fp.WorldView.SetRect(pixel.R(0, 0, xWidth, yHeight))
		fp.WorldView.CamPos = pixel.ZV
		fp.WorldView.PortPos = viewport.MainCamera.CamPos
		ChangeFullWorldShader(fp, 0)
		fp.WorldView.Canvas.SetFragmentShader(data.WorldShader)
	}
	if fp.BackgroundView == nil {
		fp.BackgroundView = viewport.New(nil)
		fp.BackgroundView.ParentView = fp.WorldView
		fp.BackgroundView.SetRect(pixel.R(0, 0, world.TileSize*w, world.TileSize*h))
		fp.BackgroundView.CamPos = pixel.V(world.TileSize*0.5*w, world.TileSize*0.5*h)
		fp.BackgroundView.PortPos = pixel.ZV
		//UpdateBackgroundShaders(fp)
		//fp.BackgroundView.Canvas.SetFragmentShader(data.BossBlobShader)
	}
	if fp.PuzzleView == nil {
		fp.PuzzleView = viewport.New(nil)
		fp.PuzzleView.ParentView = fp.WorldView
		fp.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*w, world.TileSize*h))
		fp.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*w, world.TileSize*0.5*h)
		fp.PuzzleView.PortPos = pixel.ZV
		//UpdatePuzzleShaders(fp)
		//fp.PuzzleView.Canvas.SetFragmentShader(data.PuzzleShader)
	}
	if fp.PuzzleViewNoShader == nil {
		fp.PuzzleViewNoShader = viewport.New(nil)
		fp.PuzzleViewNoShader.ParentView = fp.WorldView
		fp.PuzzleViewNoShader.SetRect(pixel.R(0, 0, world.TileSize*w, world.TileSize*h))
		fp.PuzzleViewNoShader.CamPos = pixel.V(world.TileSize*0.5*w, world.TileSize*0.5*h)
		fp.PuzzleViewNoShader.PortPos = pixel.ZV
	}
	if fp.BorderView == nil {
		fp.BorderView = viewport.New(nil)
		fp.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(w+1), world.TileSize*(h+1)))
		fp.BorderView.CamPos = pixel.V(world.TileSize*0.5*w, world.TileSize*0.5*h)
	}
}

func UpdateBackgroundShaders(fp *data.PlayArea) {
	fp.BackgroundView.Canvas.SetUniform("uTime", &data.ShaderTime)
	fp.BackgroundView.Canvas.SetUniform("uBallCount", int32(5))
	fp.BackgroundView.Canvas.SetUniform("uHeadPos", &data.BossPos)
	fp.BackgroundView.Canvas.SetUniform("uColorInner", util.RGBAToVec3(pixel.ToRGBA(constants.ColorGreen)))
	fp.BackgroundView.Canvas.SetUniform("uColorOuter", util.RGBAToVec3(pixel.ToRGBA(constants.ColorDarkGreen)))
}

func UpdatePuzzleShaders(fp *data.PlayArea) {
	// set puzzle shader uniforms
	// colors
	//setCanvasShaderColorsFromMD(fp.PuzzleView.Canvas, fp.Puzzle.Metadata)
	// world
	//fp.PuzzleView.Canvas.SetUniform("uMode", int32(fp.Puzzle.Metadata.ShaderMode))
	//fp.PuzzleView.Canvas.SetUniform("uTime", &data.ShaderTime)
	//fp.PuzzleView.Canvas.SetUniform("uSpeed", &fp.Puzzle.Metadata.ShaderSpeed)
	//fp.PuzzleView.Canvas.SetUniform("uXVar", &fp.Puzzle.Metadata.ShaderX)
	//fp.PuzzleView.Canvas.SetUniform("uYVar", &fp.Puzzle.Metadata.ShaderY)
	//fp.PuzzleView.Canvas.SetUniform("uCustom", &fp.Puzzle.Metadata.ShaderCustom)
	//fp.PuzzleView.Canvas.SetUniform("uParticle", int32(fp.Puzzle.Metadata.ParticleMode))
}

func ChangeFullWorldShader(fp *data.PlayArea, shaderMode int) {
	fp.Puzzle.Metadata.ShaderMode = shaderMode
	fp.Puzzle.Metadata.ShaderSpeed = constants.ShaderSpeeds[shaderMode]
	fp.Puzzle.Metadata.ShaderX = constants.ShaderXs[shaderMode]
	fp.Puzzle.Metadata.ShaderY = constants.ShaderYs[shaderMode]
	fp.Puzzle.Metadata.ShaderCustom = constants.ShaderCustom[shaderMode]
	UpdateWorldShaders(fp)
}

func UpdateWorldShaders(fp *data.PlayArea) {
	// colors
	setCanvasShaderColorsFromMD(fp.WorldView.Canvas, fp.Puzzle.Metadata)
	// world
	fp.WorldView.Canvas.SetUniform("uMode", int32(fp.Puzzle.Metadata.ShaderMode))
	fp.WorldView.Canvas.SetUniform("uTime", &data.ShaderTime)
	fp.WorldView.Canvas.SetUniform("uSpeed", &fp.Puzzle.Metadata.ShaderSpeed)
	fp.WorldView.Canvas.SetUniform("uXVar", &fp.Puzzle.Metadata.ShaderX)
	fp.WorldView.Canvas.SetUniform("uYVar", &fp.Puzzle.Metadata.ShaderY)
	fp.WorldView.Canvas.SetUniform("uCustom", &fp.Puzzle.Metadata.ShaderCustom)
	//fp.WorldView.Canvas.SetUniform("uParticle", int32(fp.Puzzle.Metadata.ParticleMode))
	// darkness
	darkness := int32(0)
	if fp.Puzzle.Metadata.Darkness && !data.EditorDraw {
		darkness = 1
	}
	fp.WorldView.Canvas.SetUniform("uDarkness", darkness)
	fp.WorldView.Canvas.SetUniform("uDarknessWidth", constants.DarknessWidth)
	fp.WorldView.Canvas.SetUniform("uDarknessDist", constants.DarknessDist)
	fp.WorldView.Canvas.SetUniform("uDarknessGrad", constants.DarknessGrad)
	for p := 0; p < constants.MaxPlayers; p++ {
		if fp.Level != nil {
			fp.WorldView.Canvas.SetUniform(fmt.Sprintf("uPlayer%dLoc", p+1), fp.Level.PLoc[p])
		} else {
			fp.WorldView.Canvas.SetUniform(fmt.Sprintf("uPlayer%dLoc", p+1), &mgl32.Vec2{})
		}
	}
}

func DisposePuzzle(puzzle *data.Puzzle) {
	if puzzle == nil {
		return
	}
	for _, row := range puzzle.Tiles.T {
		for _, tile := range row {
			if tile != nil {
				if tile.FloatingText != nil {
					myecs.Manager.DisposeEntity(tile.FloatingText.Entity)
					myecs.Manager.DisposeEntity(tile.FloatingText.ShEntity)
					data.RemoveFloatingText(tile)
				}
				if tile.Entity != nil {
					myecs.Manager.DisposeEntity(tile.Entity)
					tile.Entity = nil
				}
			}
		}
	}
}

func VerifyPuzzle(puzzle *data.Puzzle) {
	// fix dimensions
	if puzzle.Metadata.Width < 6 {
		puzzle.Metadata.Width = constants.PuzzleWidth
	} else if puzzle.Metadata.Width > constants.PuzzleMaxWidth {
		puzzle.Metadata.Width = constants.PuzzleMaxWidth
	}
	if puzzle.Metadata.Height < 6 {
		puzzle.Metadata.Height = constants.PuzzleHeight
	} else if puzzle.Metadata.Height > constants.PuzzleMaxHeight {
		puzzle.Metadata.Height = constants.PuzzleMaxHeight
	}
	num := puzzle.Metadata.WorldNumber
	// fix sprites
	if puzzle.Metadata.WorldSprite == "" {
		puzzle.Metadata.WorldSprite = constants.WorldSprites[num]
	}
	if puzzle.Metadata.WorldLiquid == "" {
		puzzle.Metadata.WorldLiquid = constants.WorldLiquids[num]
	}
	data.SelectedWorldIndex = num
	if puzzle.Metadata.WorldNumber == constants.WorldCustom {
		for n, w := range constants.WorldSprites {
			if puzzle.Metadata.WorldSprite == w {
				data.SelectedWorldIndex = n
			}
		}
	}
	// todo: check world colors (all black, I assume)
	// fix music
	if puzzle.Metadata.MusicTrack == "" {
		puzzle.Metadata.MusicTrack = constants.WorldMusic[num]
	}
	for y, row := range puzzle.Tiles.T {
		for x, tile := range row {
			tile.Coords = world.NewCoords(x, y)
			tile.Puzzle = puzzle
		}
	}
}

func InitLevel(playArea *data.PlayArea) {
	if playArea == nil {
		panic("InitLevel: puzzle view is nil")
	} else if playArea.Puzzle == nil {
		panic("InitLevel: puzzle is nil")
	}
	DisposePuzzle(playArea.Puzzle)
	if playArea.Level != nil {
		DisposeLevel(playArea.Level)
	}
	InitPlayAreaView(playArea)
	playArea.Level = &data.Level{
		Tiles:    playArea.Puzzle.CopyTiles(),
		Puzzle:   playArea.Puzzle,
		Metadata: playArea.Puzzle.Metadata,
	}
	data.CurrLevel = playArea.Level
	InitLevelTiles(playArea.Level)
	InitPlayers(playArea.Level)
	CreateFakePlayer(playArea.Level)
	InitLevelDialogs(playArea.Level)
	InitBoss(playArea.Level)
	//UpdateBackgroundShaders(playArea)
	//UpdatePuzzleShaders(playArea)
	ChangeFullWorldShader(playArea, playArea.Level.Metadata.ShaderMode)
	UpdatePlayAreaView(playArea)

	FloatingTextStartLevel()
	SetPuzzleTitle(playArea.Level.Metadata.Name, playArea.Level.Metadata.PrimaryColor)
	UpdatePuzzleTimer()
}

func UpdateLevelLayer(level *data.Level, offset int) {
	for _, row := range level.Tiles.T {
		for _, tile := range row {
			tile.Object.Layer += offset
		}
	}
}
