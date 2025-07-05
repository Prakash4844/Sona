package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/atotto/clipboard"
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
	mutex            sync.RWMutex
	lastClipboard    string
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
	go a.startClipboardMonitoring()
}

func (a *App) startClipboardMonitoring() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		current, err := clipboard.ReadAll()
		if err != nil {
			continue
		}

		if current != "" && current != a.lastClipboard {
			a.addClipboardItem(current)
			a.lastClipboard = current
		}
	}
}

func (a *App) addClipboardItem(content string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if content already exists
	for _, item := range a.clipboardHistory {
		if item.Content == content {
			return
		}
	}

	newItem := ClipboardItem{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Content:   content,
		Timestamp: time.Now(),
	}

	// Add to beginning of slice
	a.clipboardHistory = append([]ClipboardItem{newItem}, a.clipboardHistory...)
}

// GetClipboardHistory returns the clipboard history
func (a *App) GetClipboardHistory() []ClipboardItem {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.clipboardHistory
}
