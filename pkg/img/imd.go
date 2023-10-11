package img

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type IMDrawer struct {
	Key        string
	Index      int
	imd        *imdraw.IMDraw
	AutoDraw   bool
	AutoClear  bool
}

func AddIMDrawer(key string, autoDraw, autoClear bool) {
	if _, ok := IMDrawers[key]; ok {
		IMDrawers[key].AutoDraw = autoDraw
		IMDrawers[key].AutoClear = autoClear
	} else {
		IMDrawers[key] = &IMDrawer{
			Key:       key,
			Index:     len(imdraws),
			imd:       imdraw.New(nil),
			AutoDraw:  autoDraw,
			AutoClear: autoClear,
		}
		imdraws = append(imdraws, IMDrawers[key])
	}
}

func (i *IMDrawer) Clear() {
	i.imd.Clear()
}

func (i *IMDrawer) IMD() *imdraw.IMDraw {
	return i.imd
}

func (i *IMDrawer) Draw(target pixel.Target) {
	i.imd.Draw(target)
}