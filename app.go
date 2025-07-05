package main

import (
	"context"
	"time"
)

type ClipboardItem struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// App struct
type App struct {
	ctx              context.Context
	clipboardHistory []ClipboardItem
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		clipboardHistory: make([]ClipboardItem, 0),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetClipboardHistory returns the clipboard history
func (a *App) GetClipboardHistory() []ClipboardItem {
	return a.clipboardHistory
}
