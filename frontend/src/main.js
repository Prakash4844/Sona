import './style.css';
import './app.css';

import { GetClipboardHistory } from '../wailsjs/go/main/App';

class ClipboardApp {
    constructor() {
        this.clipboardHistory = [];
        this.init();
    }

    async init() {
        await this.loadClipboardHistory();
        this.render();
    }

    async loadClipboardHistory() {
        try {
            this.clipboardHistory = await GetClipboardHistory();
        } catch (err) {
            console.error('Failed to load clipboard history:', err);
        }
    }

    render() {
        this.renderClipboardList();
    }

    renderClipboardList() {
        const clipboardList = document.getElementById('clipboard-list');
        const emptyState = document.getElementById('empty-state');

        if (this.clipboardHistory.length === 0) {
            emptyState.style.display = 'flex';
            return;
        }

        emptyState.style.display = 'none';

        // Clear existing items
        const existingItems = clipboardList.querySelectorAll('.clipboard-item');
        existingItems.forEach(item => item.remove());

        // Render items
        this.clipboardHistory.forEach(item => {
            const itemElement = this.createClipboardItemElement(item);
            clipboardList.appendChild(itemElement);
        });
    }

    createClipboardItemElement(item) {
        const div = document.createElement('div');
        div.className = 'clipboard-item';

        const content = item.content.length > 100 
            ? item.content.substring(0, 100) + '...' 
            : item.content;

        const timestamp = new Date(item.timestamp).toLocaleString();

        div.innerHTML = `
            <div class="item-content">${content}</div>
            <div class="item-timestamp">${timestamp}</div>
        `;

        return div;
    }
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.app = new ClipboardApp();
});
