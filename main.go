package main

import (
	"context"
	"embed"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application menu
	appMenu := menu.NewMenu()
	
	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Show Sona", keys.CmdOrCtrl("1"), func(_ *menu.CallbackData) {
		app.ShowWindow()
	})
	fileMenu.AddSeparator()
	
	// Platform-specific quit shortcut
	var quitKey *keys.Accelerator
	switch runtime.GOOS {
	case "windows":
		quitKey = keys.Key("alt+f4")
	case "darwin":
		quitKey = keys.Key("cmd+q")
	default: // Linux and others
		quitKey = keys.Key("ctrl+q")
	}
	
	fileMenu.AddText("Quit", quitKey, func(_ *menu.CallbackData) {
		app.Quit()
	})

	// Edit menu
	editMenu := appMenu.AddSubmenu("Edit")
	editMenu.AddText("Clear Clipboard History", keys.CmdOrCtrl("k"), func(_ *menu.CallbackData) {
		// Note: Menu bypasses confirmation dialog for simplicity
		// Only clears non-pinned items, preserves pinned ones
		app.ClearHistory()
	})

	// Help menu
	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("About Sona", nil, func(_ *menu.CallbackData) {
		// This will be handled by the frontend
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "Sona - Clipboard Manager",
		Width:            450,
		Height:           600,
		MinWidth:         400,
		MinHeight:        500,
		MaxWidth:         800,
		MaxHeight:        1000,
		DisableResize:    false,
		Fullscreen:       false,
		Frameless:        false,
		StartHidden:      false,
		HideWindowOnClose: false,
		Menu:             appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			// Hide window when user clicks X (like system tray behavior)
			app.HideWindow()
			return true
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
