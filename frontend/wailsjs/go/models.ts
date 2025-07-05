export namespace main {
	
	export class ClipboardItem {
	    id: string;
	    content: string;
	    // Go type: time
	    timestamp: any;
	    is_pinned: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ClipboardItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.is_pinned = source["is_pinned"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Settings {
	    max_history_length: number;
	    hotkey: string;
	    ui_theme: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.max_history_length = source["max_history_length"];
	        this.hotkey = source["hotkey"];
	        this.ui_theme = source["ui_theme"];
	    }
	}

}

