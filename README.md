# Sona - Advanced Clipboard Manager

Sona is a powerful, feature-rich clipboard manager built with Wails and Go. It provides a beautiful, modern interface for managing your clipboard history with advanced features like pinning, editing, and search functionality.

## ✨ Features

### Core Functionality
- **Real-time Clipboard Monitoring**: Automatically captures everything you copy
- **Configurable History Length**: Set custom limits for your clipboard history
- **Search and Filter**: Quickly find specific clipboard items
- **Click to Copy**: Simply click any item to copy it back to your clipboard

### Advanced Features
- **📌 Pin Important Items**: Pin items to prevent them from being deleted
- **✏️ Edit Content**: Modify clipboard content before copying
- **🎨 Beautiful UI**: Modern blur glass interface with multiple theme options
- **⌨️ Keyboard Shortcuts**: Global hotkeys for quick access (Ctrl+Shift+V)
- **🎛️ Settings Management**: Customize appearance, history length, and hotkeys

### UI Features
- **Blur Glass Theme**: Elegant translucent interface
- **Dark/Light Themes**: Multiple solid color themes available
- **Responsive Design**: Works on various screen sizes
- **Toast Notifications**: Visual feedback for all actions
- **Modal Dialogs**: Clean interfaces for editing and settings

### Data Management
- **Persistent Storage**: Clipboard history saved between sessions
- **Clear History**: Option to clear non-pinned items
- **Smart Sorting**: Pinned items stay at the top
- **Duplicate Prevention**: Avoids storing duplicate entries

## 🚀 Getting Started

### Prerequisites
- Go 1.18 or later
- Node.js 16 or later
- Wails v2 CLI

### Development

To run in live development mode:
```bash
wails dev
```

This will start a Vite development server with hot reload for frontend changes.

### Building

To build a redistributable production package:
```bash
wails build
```

The built application will be available in the `build/bin/` directory.

### Installation

1. Download the latest release from the releases page
2. Extract the application to your Applications folder (macOS) or desired location
3. Run the application

## 🎮 Usage

### First Launch
1. Launch Sona
2. Start copying text - it will automatically appear in your clipboard history
3. Click any item to copy it back to your clipboard

### Keyboard Shortcuts
- `Ctrl+Shift+V` (default): Open/focus Sona window
- **Quit application**:
  - macOS: `Cmd+Q`
  - Windows: `Alt+F4`
  - Linux: `Ctrl+Q`
- `Cmd+K` (macOS) / `Ctrl+K` (other platforms): Clear clipboard history
- `Escape`: Close modal dialogs
- Search is instant as you type

### Managing Items
- **Copy**: Click any item or use the copy button
- **Edit**: Click the edit button to modify content
- **Pin**: Click the pin button to prevent deletion
- **Delete**: Click delete button (only available for non-pinned items)

### Settings
Access the settings panel to:
- Adjust maximum history length (10-1000 items)
- Change UI theme (Blur Glass, Dark, Light)
- Modify window opacity
- View current hotkey configuration

## 🛠️ Technical Details

### Architecture
- **Backend**: Go with Wails v2 framework
- **Frontend**: Vanilla JavaScript with modern ES6+ features
- **UI**: CSS3 with backdrop-filter for blur effects
- **Storage**: JSON files in user's home directory (`~/.sona/`)

### File Structure
```
~/.sona/
├── settings.json       # Application settings
└── clipboard_history.json  # Clipboard history data
```

### Dependencies
- `github.com/atotto/clipboard` - Cross-platform clipboard access
- `github.com/wailsapp/wails/v2` - Desktop app framework

## 🔧 Development

### Frontend Development
The frontend is built with vanilla JavaScript and modern CSS. Key files:
- `frontend/index.html` - Main HTML structure
- `frontend/src/main.js` - Application logic
- `frontend/src/style.css` - Styling and themes

### Backend Development
The backend is written in Go with the following key components:
- `app.go` - Main application logic and API
- `main.go` - Application entry point and configuration

### Building from Source
```bash
# Clone the repository
git clone <repository-url>
cd sona

# Install dependencies
go mod tidy

# Run in development mode
wails dev

# Build for production
wails build
```

## 🎨 Customization

### Themes
Sona supports multiple themes:
- **Blur Glass** (default): Modern translucent interface
- **Dark Solid**: Dark theme with solid backgrounds
- **Light Solid**: Light theme with solid backgrounds

### Settings
All settings are configurable through the UI:
- History length: 10-1000 items
- Window opacity: 50-100%
- UI theme selection
- Hotkey display (modification coming in future versions)

## 🚀 Roadmap

Future features planned:
- Global hotkey customization
- System tray integration
- Export/import functionality
- Rich text and image support
- Cloud synchronization
- Plugin system

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🐛 Issues

If you encounter any issues or have feature requests, please create an issue on GitHub.

---

**Sona** - Making clipboard management elegant and efficient.