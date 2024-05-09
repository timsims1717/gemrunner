package load

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"github.com/gopxl/pixel"
)

func worldDialogShaders() {
	changeWorld := data.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if scr, okScr := e1.(*data.Scroll); okScr {
			for _, e2 := range scr.Elements {
				if ct, okCt := e2.(*data.Container); okCt {
					ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
					ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
					ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
					ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
					ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
					ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
					ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
					ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
					ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
					ct.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
				}
			}
		} else if ct, okCt := e1.(*data.Container); okCt {
			ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
			ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
			ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
			ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
			ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
			ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
			ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
			ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
			ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
			ct.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
		}
	}
}

func worldDialogNormalShaders() {
	changeWorld := data.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if scr, okScr := e1.(*data.Scroll); okScr {
			i := 0
			for _, e2 := range scr.Elements {
				if ct, okCt := e2.(*data.Container); okCt {
					pc := pixel.ToRGBA(constants.WorldPrimary[i])
					sc := pixel.ToRGBA(constants.WorldSecondary[i])
					dc := pixel.ToRGBA(constants.WorldDoodad[i])
					ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
					ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
					ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
					ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
					ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
					ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
					ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
					ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
					ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
					i++
				}
			}
		} else if ct, okCt := e1.(*data.Container); okCt {
			pc := pixel.ToRGBA(constants.WorldPrimary[data.SelectedWorldIndex])
			sc := pixel.ToRGBA(constants.WorldSecondary[data.SelectedWorldIndex])
			dc := pixel.ToRGBA(constants.WorldDoodad[data.SelectedWorldIndex])
			ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
			ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
			ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
			ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
			ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
			ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
			ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
			ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
			ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
		}
	}
}

func worldDialogCustomShaders() {
	changeWorld := data.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if scr, okScr := e1.(*data.Scroll); okScr {
			for _, e2 := range scr.Elements {
				if ct, okCt := e2.(*data.Container); okCt {
					ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
					ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
					ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
					ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
					ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
					ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
					ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
					ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
					ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
				}
			}
		} else if ct, okCt := e1.(*data.Container); okCt {
			ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
			ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
			ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
			ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
			ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
			ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
			ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
			ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
			ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
		}
	}
}

func worldDialogCustomShadersPrimary() {
	changeWorld := data.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if scr, okScr := e1.(*data.Scroll); okScr {
			for _, e2 := range scr.Elements {
				if ct, okCt := e2.(*data.Container); okCt {
					ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
					ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
					ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
				}
			}
		} else if ct, okCt := e1.(*data.Container); okCt {
			ct.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
			ct.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
			ct.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
		}
	}
}

func worldDialogCustomShadersSecondary() {
	changeWorld := data.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if scr, okScr := e1.(*data.Scroll); okScr {
			for _, e2 := range scr.Elements {
				if ct, okCt := e2.(*data.Container); okCt {
					ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
					ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
					ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
				}
			}
		} else if ct, okCt := e1.(*data.Container); okCt {
			ct.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
			ct.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
			ct.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
		}
	}
}

func worldDialogCustomShadersDoodad() {
	changeWorld := data.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if scr, okScr := e1.(*data.Scroll); okScr {
			for _, e2 := range scr.Elements {
				if ct, okCt := e2.(*data.Container); okCt {
					ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
					ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
					ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
				}
			}
		} else if ct, okCt := e1.(*data.Container); okCt {
			ct.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
			ct.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
			ct.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
		}
	}
}

func changeSelectedColor(key string) {
	switch key {
	case "red_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorPink)
	case "yellow_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorBurnt)
	case "red_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorPink)
	case "yellow_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorBurnt)
	case "red_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorPink)
	case "yellow_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorBurnt)
	}
}

func updateColorCheckbox(x *data.Checkbox) {
	switch x.Key {
	case "red_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_primary":
		data.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBurnt))
	case "red_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_secondary":
		data.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBurnt))
	case "red_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_doodad":
		data.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBurnt))
	}
}
