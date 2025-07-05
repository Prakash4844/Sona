package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ClipboardItem struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	IsPinned  bool      `json:"is_pinned"`
}

type Settings struct {
	MaxHistoryLength int    `json:"max_history_length"`
	Hotkey           string `json:"hotkey"`
	UITheme          string `json:"ui_theme"`
}

// App struct
type App struct {
	ctx              context.Context
	clipboardHistory []ClipboardItem
	settings         Settings
	mutex            sync.RWMutex
	lastClipboard    string
	settingsFile     string
	dataFile         string
}

// NewApp creates a new App application struct
func NewApp() *App {
	homeDir, _ := os.UserHomeDir()
	appDir := filepath.Join(homeDir, ".sona")
	os.MkdirAll(appDir, 0755)

	app := &App{
		clipboardHistory: make([]ClipboardItem, 0),
		settings: Settings{
			MaxHistoryLength: 50,
			Hotkey:           "ctrl+shift+v",
			UITheme:          "blur",
		},
		settingsFile: filepath.Join(appDir, "settings.json"),
		dataFile:     filepath.Join(appDir, "clipboard_history.json"),
	}

	app.loadSettings()
	app.loadClipboardHistory()
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.startClipboardMonitoring()
}

func (a *App) startClipboardMonitoring() {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// Try original library first for better UTF-8 handling
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

// readClipboardNative uses platform-specific commands for better UTF-8 support
func (a *App) readClipboardNative() (string, error) {
	var cmd *exec.Cmd
	
	switch goruntime.GOOS {
	case "darwin":
		// macOS: use pbpaste
		cmd = exec.Command("pbpaste")
	case "linux":
		// Linux: try xclip first, fallback to xsel
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case "windows":
		// Windows: use powershell Get-Clipboard
		cmd = exec.Command("powershell", "-command", "Get-Clipboard")
	default:
		// Fallback to original library
		return clipboard.ReadAll()
	}
	
	output, err := cmd.Output()
	if err != nil {
		// Fallback to original library if native command fails
		return clipboard.ReadAll()
	}
	
	// Remove trailing newline that pbpaste adds
	content := strings.TrimSuffix(string(output), "\n")
	
	// Debug: Print raw bytes and string conversion
	fmt.Printf("DEBUG: Raw clipboard bytes: %v\n", output)
	fmt.Printf("DEBUG: String conversion: %q\n", content)
	fmt.Printf("DEBUG: UTF-8 valid after conversion: %v\n", utf8.ValidString(content))
	
	// Ensure valid UTF-8
	if !utf8.ValidString(content) {
		content = strings.ToValidUTF8(content, "�")
		fmt.Printf("DEBUG: Cleaned to valid UTF-8: %q\n", content)
	}
	
	return content, nil
}

// writeClipboardNative uses platform-specific commands for better UTF-8 support
func (a *App) writeClipboardNative(content string) error {
	var cmd *exec.Cmd
	
	switch goruntime.GOOS {
	case "darwin":
		// macOS: use pbcopy
		cmd = exec.Command("pbcopy")
	case "linux":
		// Linux: try xclip first
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		// Windows: use clip command
		cmd = exec.Command("clip")
	default:
		// Fallback to original library
		return clipboard.WriteAll(content)
	}
	
	cmd.Stdin = strings.NewReader(content)
	err := cmd.Run()
	if err != nil {
		// Fallback to original library if native command fails
		return clipboard.WriteAll(content)
	}
	
	return nil
}

// cleanUTF8String ensures the string is valid UTF-8 and handles any encoding issues
func (a *App) cleanUTF8String(s string) string {
	if s == "" {
		return s
	}

	// Check if string is valid UTF-8
	if utf8.ValidString(s) {
		return s
	}

	// Log when we encounter invalid UTF-8
	fmt.Printf("Warning: Invalid UTF-8 detected in clipboard content, cleaning...\n")
	
	// If not valid UTF-8, try to clean it
	// Convert invalid sequences to replacement characters
	cleaned := strings.ToValidUTF8(s, "�")
	
	fmt.Printf("Original length: %d, Cleaned length: %d\n", len(s), len(cleaned))
	return cleaned
}

func (a *App) addClipboardItem(content string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Debug: Log the raw content
	fmt.Printf("DEBUG: Adding clipboard item - Original content: %q\n", content)
	fmt.Printf("DEBUG: Content bytes: %v\n", []byte(content))
	fmt.Printf("DEBUG: UTF-8 valid: %v\n", utf8.ValidString(content))
	
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
		IsPinned:  false,
	}

	fmt.Printf("DEBUG: New item content: %q\n", newItem.Content)

	// Add to beginning of slice
	a.clipboardHistory = append([]ClipboardItem{newItem}, a.clipboardHistory...)

	// Remove oldest non-pinned items if exceeding max length
	if len(a.clipboardHistory) > a.settings.MaxHistoryLength {
		// Keep pinned items and remove oldest non-pinned
		var newHistory []ClipboardItem
		pinnedCount := 0
		
		for _, item := range a.clipboardHistory {
			if item.IsPinned {
				newHistory = append(newHistory, item)
				pinnedCount++
			} else if len(newHistory)-pinnedCount < a.settings.MaxHistoryLength-pinnedCount {
				newHistory = append(newHistory, item)
			}
		}
		a.clipboardHistory = newHistory
	}

	a.saveClipboardHistory()
	
	// Emit event to frontend
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "clipboard-updated", a.clipboardHistory)
	}
}


// Frontend API methods
func (a *App) GetClipboardHistory() []ClipboardItem {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	
	// Sort: pinned items first, then by timestamp (newest first)
	var pinned, unpinned []ClipboardItem
	for _, item := range a.clipboardHistory {
		if item.IsPinned {
			pinned = append(pinned, item)
		} else {
			unpinned = append(unpinned, item)
		}
	}
	
	result := append(pinned, unpinned...)
	return result
}

func (a *App) CopyToClipboard(content string) error {
	fmt.Printf("DEBUG: CopyToClipboard called with: %q\n", content)
	fmt.Printf("DEBUG: Content bytes: %v\n", []byte(content))
	
	a.lastClipboard = content
	err := clipboard.WriteAll(content)
	
	fmt.Printf("DEBUG: Write to clipboard error: %v\n", err)
	return err
}

func (a *App) EditClipboardItem(id, newContent string) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Clean the new content to ensure valid UTF-8
	cleanContent := a.cleanUTF8String(newContent)

	for i, item := range a.clipboardHistory {
		if item.ID == id {
			a.clipboardHistory[i].Content = cleanContent
			a.clipboardHistory[i].Timestamp = time.Now()
			a.saveClipboardHistory()
			runtime.EventsEmit(a.ctx, "clipboard-updated", a.clipboardHistory)
			return true
		}
	}
	return false
}

func (a *App) TogglePinItem(id string) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for i, item := range a.clipboardHistory {
		if item.ID == id {
			a.clipboardHistory[i].IsPinned = !a.clipboardHistory[i].IsPinned
			a.saveClipboardHistory()
			runtime.EventsEmit(a.ctx, "clipboard-updated", a.clipboardHistory)
			return a.clipboardHistory[i].IsPinned
		}
	}
	return false
}

func (a *App) DeleteClipboardItem(id string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for i, item := range a.clipboardHistory {
		if item.ID == id {
			if item.IsPinned {
				return fmt.Errorf("cannot delete pinned item")
			}
			a.clipboardHistory = append(a.clipboardHistory[:i], a.clipboardHistory[i+1:]...)
			a.saveClipboardHistory()
			runtime.EventsEmit(a.ctx, "clipboard-updated", a.clipboardHistory)
			return nil
		}
	}
	return fmt.Errorf("item not found")
}

func (a *App) ClearHistory() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Keep only pinned items
	var pinnedItems []ClipboardItem
	for _, item := range a.clipboardHistory {
		if item.IsPinned {
			pinnedItems = append(pinnedItems, item)
		}
	}
	
	a.clipboardHistory = pinnedItems
	a.saveClipboardHistory()
	runtime.EventsEmit(a.ctx, "clipboard-updated", a.clipboardHistory)
}

func (a *App) GetSettings() Settings {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.settings
}

func (a *App) UpdateSettings(newSettings Settings) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	a.settings = newSettings
	a.saveSettings()
}

func (a *App) loadSettings() {
	data, err := ioutil.ReadFile(a.settingsFile)
	if err != nil {
		return
	}
	
	json.Unmarshal(data, &a.settings)
}

func (a *App) saveSettings() {
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		return
	}
	
	// Ensure UTF-8 encoding when writing file
	err = ioutil.WriteFile(a.settingsFile, data, 0644)
	if err != nil {
		fmt.Printf("Error saving settings: %v\n", err)
	}
}

func (a *App) loadClipboardHistory() {
	data, err := ioutil.ReadFile(a.dataFile)
	if err != nil {
		return
	}
	
	// Validate UTF-8 before unmarshaling
	if !utf8.Valid(data) {
		fmt.Println("Warning: clipboard history file contains invalid UTF-8")
		return
	}
	
	err = json.Unmarshal(data, &a.clipboardHistory)
	if err != nil {
		fmt.Printf("Error loading clipboard history: %v\n", err)
	}
}

func (a *App) saveClipboardHistory() {
	data, err := json.MarshalIndent(a.clipboardHistory, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling clipboard history: %v\n", err)
		return
	}
	
	// Ensure UTF-8 encoding when writing file
	err = ioutil.WriteFile(a.dataFile, data, 0644)
	if err != nil {
		fmt.Printf("Error saving clipboard history: %v\n", err)
	}
}

func (a *App) HideWindow() {
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
	}
}

func (a *App) MinimizeWindow() {
	if a.ctx != nil {
		runtime.WindowMinimise(a.ctx)
	}
}

func (a *App) ShowWindow() {
	if a.ctx != nil {
		runtime.WindowShow(a.ctx)
	}
}

func (a *App) Quit() {
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

func (a *App) domReady(ctx context.Context) {
	// DOM is ready, can perform additional initialization here
}
