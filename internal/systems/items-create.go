package systems

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func CreateItem(block data.Block, pos pixel.Vec, key string, metadata data.TileMetadata, coords world.Coords) *data.BasicItem {
	var item *data.BasicItem
	switch block {
	case data.BlockJumpBoots:
		item = CreateJumpBoots(pos, key, metadata, coords)
	case data.BlockBox:
		item = CreateBox(pos, key, metadata, coords)
	case data.BlockKey:
		item = CreateKey(pos, key, metadata, coords)
	case data.BlockBigBomb:
		item = CreateBomb(pos, key, metadata, coords, "big", true)
	case data.BlockSmallBomb:
		item = CreateBomb(pos, key, metadata, coords, "small", false)
	case data.BlockJetpack:
		item = CreateJetpack(pos, key, metadata, coords)
	case data.BlockDisguise:
		item = CreateDisguise(pos, key, metadata, coords)
	case data.BlockDrill:
		item = CreateDrill(pos, key, metadata, coords)
	case data.BlockFlamethrower:
		item = CreateFlamethrower(pos, key, metadata, coords)
	case data.BlockGoopBucket:
		item = CreateGoopBucket(pos, key, metadata, coords)
	case data.BlockAirCannon:
		item = CreateAirCannon(pos, key, metadata, coords)
	case data.BlockTransporter:
		item = CreateTransporter(pos, key, metadata, coords)
	case data.BlockSnare:
		item = CreateSnare(pos, key, metadata, coords)
	}
	if item != nil {
		item.Block = block
	}
	return item
}
