import './style.css';
import './app.css';

import { 
    GetClipboardHistory, 
    CopyToClipboard, 
    EditClipboardItem, 
    TogglePinItem, 
    DeleteClipboardItem, 
    ClearHistory, 
    GetSettings, 
    UpdateSettings,
    HideWindow,
    MinimizeWindow
} from '../wailsjs/go/main/App';

import { EventsOn } from '../wailsjs/runtime/runtime';

class ClipboardApp {
    constructor() {
        this.clipboardHistory = [];
        this.filteredHistory = [];
        this.settings = {};
        this.currentEditId = null;
        
        this.init();
    }

    async init() {
        await this.loadSettings();
        await this.loadClipboardHistory();
        this.setupWailsEvents();
        this.render();
        
        // Delay event listener setup to ensure DOM is fully ready
        setTimeout(() => {
            this.setupEventListeners();
        }, 100);
    }

    async loadSettings() {
        try {
            this.settings = await GetSettings();
            this.applyTheme();
        } catch (err) {
            console.error('Failed to load settings:', err);
        }
    }

    async loadClipboardHistory() {
        try {
            this.clipboardHistory = await GetClipboardHistory();
            this.filteredHistory = [...this.clipboardHistory];
        } catch (err) {
            console.error('Failed to load clipboard history:', err);
        }
    }


    setupWailsEvents() {
        EventsOn('clipboard-updated', (newHistory) => {
            this.clipboardHistory = newHistory;
            this.filterHistory();
            this.renderClipboardList();
        });
    }

    setupEventListeners() {
        // Search functionality
        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.addEventListener('input', (e) => {
                this.filterHistory(e.target.value);
                this.renderClipboardList();
            });
        }

        // Header actions
        const settingsBtn = document.getElementById('settings-btn');
        if (settingsBtn) {
            settingsBtn.addEventListener('click', () => {
                this.showSettingsModal();
            });
        }


        const hideBtn = document.getElementById('hide-btn');
        if (hideBtn) {
            hideBtn.addEventListener('click', () => {
                MinimizeWindow();
            });
        }


        // Settings modal
        const closeSettingsBtn = document.getElementById('close-settings');
        if (closeSettingsBtn) {
            closeSettingsBtn.addEventListener('click', () => {
                this.hideModal('settings-modal');
            });
        }

        const saveSettingsBtn = document.getElementById('save-settings');
        if (saveSettingsBtn) {
            saveSettingsBtn.addEventListener('click', () => {
                this.saveSettings();
            });
        }

        const aboutBtn = document.getElementById('about-btn');
        if (aboutBtn) {
            aboutBtn.addEventListener('click', () => {
                this.hideModal('settings-modal');
                this.showModal('about-modal');
            });
        }

        // About modal
        const closeAboutBtn = document.getElementById('close-about');
        if (closeAboutBtn) {
            closeAboutBtn.addEventListener('click', () => {
                this.hideModal('about-modal');
            });
        }

        const closeAboutBtnFooter = document.getElementById('close-about-btn');
        if (closeAboutBtnFooter) {
            closeAboutBtnFooter.addEventListener('click', () => {
                this.hideModal('about-modal');
            });
        }

        // Edit modal
        const closeEditBtn = document.getElementById('close-edit');
        if (closeEditBtn) {
            closeEditBtn.addEventListener('click', () => {
                this.hideModal('edit-modal');
            });
        }

        const cancelEditBtn = document.getElementById('cancel-edit');
        if (cancelEditBtn) {
            cancelEditBtn.addEventListener('click', () => {
                this.hideModal('edit-modal');
            });
        }

        const saveEditBtn = document.getElementById('save-edit');
        if (saveEditBtn) {
            saveEditBtn.addEventListener('click', () => {
                this.saveEdit();
            });
        }


        // Keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                this.hideAllModals();
            }
            
        });

        // Modal click outside to close
        document.querySelectorAll('.modal-overlay').forEach(overlay => {
            overlay.addEventListener('click', (e) => {
                if (e.target === overlay) {
                    this.hideModal(overlay.id);
                }
            });
        });
    }

    filterHistory(searchTerm = '') {
        if (!searchTerm.trim()) {
            this.filteredHistory = [...this.clipboardHistory];
        } else {
            const term = searchTerm.toLowerCase();
            this.filteredHistory = this.clipboardHistory.filter(item =>
                item.content.toLowerCase().includes(term)
            );
        }
    }

    render() {
        this.renderClipboardList();
    }

    renderClipboardList() {
        const clipboardList = document.getElementById('clipboard-list');
        const emptyState = document.getElementById('empty-state');

        if (this.filteredHistory.length === 0) {
            emptyState.style.display = 'flex';
            // Remove existing items
            const existingItems = clipboardList.querySelectorAll('.clipboard-item');
            existingItems.forEach(item => item.remove());
            return;
        }

        emptyState.style.display = 'none';

        // Clear existing items
        const existingItems = clipboardList.querySelectorAll('.clipboard-item');
        existingItems.forEach(item => item.remove());

        // Render items
        this.filteredHistory.forEach(item => {
            const itemElement = this.createClipboardItemElement(item);
            clipboardList.appendChild(itemElement);
        });
    }

    createClipboardItemElement(item) {
        const div = document.createElement('div');
        div.className = `clipboard-item ${item.is_pinned ? 'pinned' : ''}`;
        div.setAttribute('data-id', item.id);

        const content = item.content.length > 200 
            ? item.content.substring(0, 200) + '...' 
            : item.content;

        const timestamp = new Date(item.timestamp).toLocaleString();

        // Create content div and set text content directly to preserve UTF-8
        const contentDiv = document.createElement('div');
        contentDiv.className = 'item-content';
        contentDiv.textContent = content;

        div.innerHTML = `
            <div class="item-meta">
                <span class="item-timestamp">${timestamp}</span>
                <div class="item-actions">
                    <button class="action-btn copy-btn" title="Copy to clipboard">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                        </svg>
                    </button>
                    <button class="action-btn edit-btn" title="Edit">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                        </svg>
                    </button>
                    <button class="action-btn pin-btn ${item.is_pinned ? 'pinned' : ''}" title="${item.is_pinned ? 'Unpin' : 'Pin'}">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="12" y1="17" x2="12" y2="22"></line>
                            <path d="M5 17h14v-1.76a2 2 0 0 0-1.11-1.79l-1.78-.9A2 2 0 0 1 15 10.76V6h1a2 2 0 0 0 0-4H8a2 2 0 0 0 0 4h1v4.76a2 2 0 0 1-1.11 1.79l-1.78.9A2 2 0 0 0 5 15.24V17z"></path>
                        </svg>
                    </button>
                    ${!item.is_pinned ? `
                        <button class="action-btn delete-btn" title="Delete">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <polyline points="3,6 5,6 21,6"></polyline>
                                <path d="M19,6v14a2,2 0,0,1 -2,2H7a2,2 0,0,1 -2,-2V6m3,0V4a2,2 0,0,1 2,-2h4a2,2 0,0,1 2,2v2"></path>
                            </svg>
                        </button>
                    ` : ''}
                </div>
            </div>
        `;

        // Insert the content div at the beginning
        div.insertBefore(contentDiv, div.firstChild);

        // Add event listeners
        div.querySelector('.copy-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            this.copyToClipboard(item.content);
        });

        div.querySelector('.edit-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            this.editItem(item);
        });

        div.querySelector('.pin-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            this.togglePin(item.id);
        });

        const deleteBtn = div.querySelector('.delete-btn');
        if (deleteBtn) {
            deleteBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.deleteItem(item.id);
            });
        }

        // Click to copy
        div.addEventListener('click', () => {
            this.copyToClipboard(item.content);
        });

        return div;
    }

    async copyToClipboard(content) {
        try {
            await CopyToClipboard(content);
            this.showToast('Copied to clipboard!');
        } catch (err) {
            console.error('Failed to copy:', err);
            this.showToast('Failed to copy', 'error');
        }
    }

    editItem(item) {
        this.currentEditId = item.id;
        document.getElementById('edit-content').value = item.content;
        this.showModal('edit-modal');
    }

    async saveEdit() {
        if (!this.currentEditId) return;

        const newContent = document.getElementById('edit-content').value;
        if (!newContent.trim()) return;

        try {
            await EditClipboardItem(this.currentEditId, newContent);
            this.hideModal('edit-modal');
            this.currentEditId = null;
            this.showToast('Item updated!');
        } catch (err) {
            console.error('Failed to edit item:', err);
            this.showToast('Failed to update item', 'error');
        }
    }

    async togglePin(id) {
        try {
            await TogglePinItem(id);
        } catch (err) {
            console.error('Failed to toggle pin:', err);
            this.showToast('Failed to pin/unpin item', 'error');
        }
    }

    async deleteItem(id) {
        try {
            await DeleteClipboardItem(id);
            this.showToast('Item deleted!');
        } catch (err) {
            console.error('Failed to delete item:', err);
            this.showToast(err.toString().replace('Error: ', ''), 'error');
        }
    }

    async clearHistory() {
        if (confirm('Are you sure you want to clear the clipboard history? Pinned items will be preserved.')) {
            try {
                await ClearHistory();
                this.showToast('History cleared!');
            } catch (err) {
                console.error('Failed to clear history:', err);
                this.showToast('Failed to clear history', 'error');
            }
        }
    }

    showSettingsModal() {
        // Populate settings form
        document.getElementById('max-history').value = this.settings.max_history_length || 100;
        document.getElementById('hotkey').value = this.settings.hotkey || 'ctrl+shift+v';
        document.getElementById('theme').value = this.settings.ui_theme || 'blur';

        this.showModal('settings-modal');
    }

    async saveSettings() {
        const newSettings = {
            max_history_length: parseInt(document.getElementById('max-history').value),
            hotkey: document.getElementById('hotkey').value,
            ui_theme: document.getElementById('theme').value
        };

        try {
            await UpdateSettings(newSettings);
            this.settings = newSettings;
            this.applyTheme();
            this.hideModal('settings-modal');
            this.showToast('Settings saved!');
        } catch (err) {
            console.error('Failed to save settings:', err);
            this.showToast('Failed to save settings', 'error');
        }
    }

    applyTheme() {
        const body = document.body;
        body.className = ''; // Clear existing theme classes
        
        if (this.settings.ui_theme) {
            body.classList.add(`theme-${this.settings.ui_theme}`);
        }

    }

    showModal(modalId) {
        document.getElementById(modalId).style.display = 'flex';
    }

    hideModal(modalId) {
        document.getElementById(modalId).style.display = 'none';
    }

    hideAllModals() {
        document.querySelectorAll('.modal-overlay').forEach(modal => {
            modal.style.display = 'none';
        });
    }

    showToast(message, type = 'success') {
        // Create toast element
        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.textContent = message;
        toast.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 12px 20px;
            background: ${type === 'error' ? 'var(--danger-color)' : 'var(--success-color)'};
            color: white;
            border-radius: 6px;
            z-index: 10000;
            font-size: 14px;
            font-weight: 500;
            opacity: 0;
            transform: translateX(100%);
            transition: all 0.3s ease;
        `;

        document.body.appendChild(toast);

        // Animate in
        setTimeout(() => {
            toast.style.opacity = '1';
            toast.style.transform = 'translateX(0)';
        }, 10);

        // Remove after 3 seconds
        setTimeout(() => {
            toast.style.opacity = '0';
            toast.style.transform = 'translateX(100%)';
            setTimeout(() => {
                if (toast.parentNode) {
                    toast.parentNode.removeChild(toast);
                }
            }, 300);
        }, 3000);
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Ensure DOM is fully ready before initializing
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initApp);
    } else {
        initApp();
    }
});

function initApp() {
    window.app = new ClipboardApp();
}

