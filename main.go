package main

import (
	"github.com/mrinjamul/authenticator-desktop/config"
	"github.com/mrinjamul/authenticator-desktop/pages/pager"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	window := app.NewWindow("Authenticator Desktop")

	window.Resize(fyne.NewSize(400, 400))
	config := config.Initialize(window)
	pager := pager.Init(config)
	pager.Start()

	window.ShowAndRun()
}
