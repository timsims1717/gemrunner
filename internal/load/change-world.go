package load

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
	"github.com/gopxl/pixel"
)

func worldDialogShaders() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
					e2.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
			e1.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
		}
	}
}

func worldDialogNormalShaders() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			i := 0
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					pc := pixel.ToRGBA(constants.WorldPrimary[i])
					sc := pixel.ToRGBA(constants.WorldSecondary[i])
					dc := pixel.ToRGBA(constants.WorldDoodad[i])
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
					i++
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			pc := pixel.ToRGBA(constants.WorldPrimary[data.SelectedWorldIndex])
			sc := pixel.ToRGBA(constants.WorldSecondary[data.SelectedWorldIndex])
			dc := pixel.ToRGBA(constants.WorldDoodad[data.SelectedWorldIndex])
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
		}
	}
}

func worldDialogCustomShaders() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
		}
	}
}

func worldDialogCustomShadersPrimary() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
		}
	}
}

func worldDialogCustomShadersSecondary() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
		}
	}
}

func worldDialogCustomShadersDoodad() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
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

func updateColorCheckbox(x *ui.Element) {
	switch x.Key {
	case "red_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBurnt))
	case "red_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBurnt))
	case "red_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBurnt))
	}
}
