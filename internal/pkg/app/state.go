package app

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.bmvs.io/ynab"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GlobalState struct {
	DB             *gorm.DB
	YNABClient     ynab.ClientServicer
	Homescreen     fyne.CanvasObject
	Window         fyne.Window
	Application    fyne.App
	CurrentScreen  fyne.CanvasObject
	PreviousScreen *fyne.CanvasObject
	NextScreen     *fyne.CanvasObject
	RootComponent  Component
	Data           map[string]interface{}
}

func NewGlobalState() *GlobalState {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("world", db)
	err = InitializeDB(db)
	a := app.NewWithID("dev.benwaters.app")
	w := a.NewWindow("YNAB Desktop Importer")
	return &GlobalState{
		Application:   a,
		Window:        w,
		DB:            db,
		RootComponent: *NewComponent(container.NewWithoutLayout()),
		Data:          make(map[string]interface{}),
	}
}

func (s GlobalState) GetWindow() fyne.Window {
	return s.Window
}

func (s *GlobalState) UpdateStateData(key string, data interface{}) {
	s.Data[key] = data
}

func (s *GlobalState) Run() {
	s.Window.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("foo", func() {
				log.Println("foo")
			}),
			fyne.NewMenuItem("bar", func() {
				log.Println("bar")
			})),
	))
	s.Window.SetContent(s.RootComponent.Object)
	//w.SetFullScreen(true)

	s.Window.ShowAndRun()
}

func NewApp() *GlobalState {
	state := NewGlobalState()
	// rootComponent with toolbar
	state.RootComponent.AddChild(
		NewComponent(container.NewVBox(widget.NewToolbar(
			widget.NewToolbarAction(
				theme.HomeIcon(),
				func() {
					log.Println("clicked home")
				},
			),
			widget.NewToolbarSpacer(),
		)),
		))
	state.RootComponent.AddChildDeep(
		NewComponent(
			NewHomeScreen(*state.NewUploadFlow(), *state.NewSettingsScreen()),
		),
	)
	return state
}
