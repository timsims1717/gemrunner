package typeface

import "github.com/gopxl/pixel"

func (item *Text) Len() int {
	return len(item.dotPosArray)
}

func (item *Text) GetStartPos() pixel.Vec {
	if len(item.dotPosArray) > 0 {
		return item.dotPosArray[0]
	}
	return item.Text.Orig
}

func (item *Text) GetEndPos() pixel.Vec {
	if len(item.dotPosArray) > 0 {
		return item.dotPosArray[len(item.dotPosArray)-1]
	}
	return item.Text.Orig
}

func (item *Text) GetDotPos(i int) pixel.Vec {
	if len(item.dotPosArray) > 0 && i < len(item.dotPosArray) {
		return item.dotPosArray[i]
	}
	return item.Text.Orig
}
