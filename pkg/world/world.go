package world

import (
	"github.com/gopxl/pixel"
)

var (
	TileSize float64
	Origin   = Coords{
		X: 0,
		Y: 0,
	}
	TileRect pixel.Rect
)

func SetTileSize(s float64) {
	TileSize = s
	TileRect = pixel.R(0, 0, s, s)
}

func MapToWorld(a Coords) pixel.Vec {
	return pixel.V(float64(a.X)*TileSize, float64(a.Y)*TileSize)
}

func WorldToMap(x, y float64) (int, int) {
	xi := int(x / TileSize)
	yi := int(y / TileSize)
	if x < 0 && x != float64(xi) {
		xi--
	}
	if y < 0 && y != float64(yi) {
		yi--
	}
	return xi, yi
}
