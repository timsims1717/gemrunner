package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/util"
)

func MagicSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC &&
			ch.Player > -1 &&
			ch.Player < constants.MaxPlayers &&
			reanimator.FrameSwitch &&
			len(ch.StoredBlocks) > 0 {
			var regened []int
			for i, t := range ch.StoredBlocks {
				if t.Flags.Regen {
					regened = append(regened, i)
				}
			}
			for i := len(ch.StoredBlocks) - 1; i >= 0; i-- {
				if util.Contains(i, regened) {
					if len(ch.StoredBlocks) > 1 {
						ch.StoredBlocks = append(ch.StoredBlocks[:i], ch.StoredBlocks[i+1:]...)
					} else {
						ch.StoredBlocks = []*data.Tile{}
					}
				}
			}
		}
	}
}
