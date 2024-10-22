package typeface

import "github.com/gopxl/pixel"

func (item *Text) Len() int {
	return len(item.dotPosArray)
}

func (item *Text) GetStartPos() pixel.Vec {
	if len(item.dotPosArray) > 0 {
		return item.GetDotPos(0)
	}
	return item.Text.Orig
}

func (item *Text) GetEndPos() pixel.Vec {
	if len(item.dotPosArray) > 0 {
		return item.GetDotPos(len(item.dotPosArray) - 1)
	}
	return item.Text.Orig
}

func (item *Text) GetDotPos(i int) pixel.Vec {
	if len(item.dotPosArray) > 0 && i < len(item.dotPosArray) {
		return item.dotPosArray[i]
	}
	return item.Text.Orig
}

func (item *Text) GetStartOfLine(i int) (int, pixel.Vec) {
	count := 0
	for _, line := range item.rawLines {
		if count+len(line) > i {
			return count, item.GetDotPos(count)
		}
		count += len(line)
	}
	return 0, item.GetStartPos()
}

func (item *Text) GetEndOfLine(i int) (int, pixel.Vec) {
	count := 0
	for j, line := range item.rawLines {
		if j == len(item.rawLines)-1 {
			break
		}
		if count+len(line) > i {
			return count + len(line), item.GetDotPos(count + len(line))
		}
		count += len(line)
	}
	return len(item.dotPosArray) - 1, item.GetEndPos()
}
