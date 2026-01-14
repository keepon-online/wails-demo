export namespace main {
	
	export class TraySettings {
	    minimizeToTray: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TraySettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.minimizeToTray = source["minimizeToTray"];
	    }
	}

}

export namespace updater {
	
	export class UpdateInfo {
	    available: boolean;
	    currentVersion: string;
	    latestVersion: string;
	    releaseNotes: string;
	    releaseUrl: string;
	    installerUrl: string;
	    debugInfo: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.releaseNotes = source["releaseNotes"];
	        this.releaseUrl = source["releaseUrl"];
	        this.installerUrl = source["installerUrl"];
	        this.debugInfo = source["debugInfo"];
	    }
	}
	export class UpdateProgress {
	    status: string;
	    progress: number;
	    message: string;
	    needRestart: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.message = source["message"];
	        this.needRestart = source["needRestart"];
	    }
	}

}

