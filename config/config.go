package config

import (
	"fyne.io/fyne/v2"
)

type Config interface {
	GetWindow() fyne.Window
	Launch() func(key string)
	SetLauncher(launcher func(key string))
}

type config struct {
	window   fyne.Window
	launcher func(key string)
}

func (conf *config) GetWindow() fyne.Window {
	return conf.window
}

func (conf *config) Launch() func(key string) {
	return conf.launcher
}

func (conf *config) SetLauncher(launcher func(key string)) {
	conf.launcher = launcher
}

func Initialize(w fyne.Window) Config {
	return &config{
		window: w,
	}
}
