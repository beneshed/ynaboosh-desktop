package app

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type About struct {
	Author       string
	MajorVersion int
	MinorVersion int
	PatchVersion int
	License      string
}

func NewAbout() About {
	return About{
		Author:       "Ben Waters",
		MajorVersion: 0,
		MinorVersion: 1,
		PatchVersion: 0,
		License:      "MIT",
	}
}

func (a About) NewAboutScreen() *fyne.Container {
	authorLabel := widget.NewLabel("Author")
	author := widget.NewLabel(a.Author)
	versionLabel := widget.NewLabel("Version")
	version := widget.NewLabel(strings.Join([]string{strconv.Itoa(a.MajorVersion), strconv.Itoa(a.MinorVersion), strconv.Itoa(a.PatchVersion)}, "."))
	return fyne.NewContainer(container.New(layout.NewFormLayout(), authorLabel, author, versionLabel, version))
}
