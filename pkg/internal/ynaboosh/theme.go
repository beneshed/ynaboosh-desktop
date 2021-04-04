package ynaboosh

import 
(
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

var (
	//go:embed assets/NotoSansHebrew-Regular.ttf
	hebrewFontRegular []byte
	//go:embed assets/NotoSansHebrew-Black.ttf
	hebrewFontBlack []byte
)

type hebrewFontTheme struct {}

func (t hebrewFontTheme) Font(_ fyne.TextStyle) fyne.Resource {
	font := fyne.NewStaticResource("font", hebrewFontBlack)
	return font
}

func (t hebrewFontTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t hebrewFontTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t hebrewFontTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var _ fyne.Theme = (*hebrewFontTheme)(nil)