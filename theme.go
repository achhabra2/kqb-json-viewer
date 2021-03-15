package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var bglLightColor color.RGBA = color.RGBA{240, 240, 240, 255}
var bglMedColor color.RGBA = color.RGBA{37, 40, 39, 255}
var bglDarkColor color.RGBA = color.RGBA{8, 8, 8, 255}

var bglBlueLight color.RGBA = color.RGBA{51, 153, 255, 255}
var bglBlueMed color.RGBA = color.RGBA{0, 71, 255, 255}
var bglBlueDark color.RGBA = color.RGBA{45, 60, 179, 255}

var bglGoldLight color.RGBA = color.RGBA{255, 204, 102, 255}
var bglGoldMed color.RGBA = color.RGBA{253, 184, 51, 255}
var bglGoldDark color.RGBA = color.RGBA{179, 123, 45, 255}

var highlightColor color.RGBA = color.RGBA{50, 50, 50, 100}
var primaryColor color.RGBA = color.RGBA{160, 85, 251, 255}
var buttonColor color.RGBA = color.RGBA{90, 112, 148, 255}

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return bglMedColor
	case theme.ColorNameButton:
		return buttonColor
	case theme.ColorNamePrimary:
		return primaryColor
	case theme.ColorNameHover:
		return highlightColor
	case theme.ColorNameForeground:
		return bglLightColor
	case theme.ColorNameInputBackground:
		return buttonColor
	case theme.ColorNameFocus:
		return bglGoldDark
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// if name == theme.IconNameHome {
	// 	fyne.NewStaticResource("myHome", homeBytes)
	// }

	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	// defaultFont := fyne.NewStaticResource("Pirulenrg", font.Pirulenrg)
	// return defaultFont
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var _ fyne.Theme = (*myTheme)(nil)
