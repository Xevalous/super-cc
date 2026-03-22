export namespace commands {
	
	export class AppInfo {
	    installed: boolean;
	    path: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.installed = source["installed"];
	        this.path = source["path"];
	        this.version = source["version"];
	    }
	}
	export class CleanupResult {
	    lastCheck: string;
	    status: string;
	    filesFound: number;
	    filesCleaned: number;
	    residualPaths: string[];
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new CleanupResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lastCheck = source["lastCheck"];
	        this.status = source["status"];
	        this.filesFound = source["filesFound"];
	        this.filesCleaned = source["filesCleaned"];
	        this.residualPaths = source["residualPaths"];
	        this.errors = source["errors"];
	    }
	}
	export class CrackStatus {
	    mode: string;
	    vipEnabled: boolean;
	    watermarkRemoved: boolean;
	    dllFound: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CrackStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.vipEnabled = source["vipEnabled"];
	        this.watermarkRemoved = source["watermarkRemoved"];
	        this.dllFound = source["dllFound"];
	    }
	}
	export class InstallationInfo {
	    status: string;
	    version: string;
	    path: string;
	    size: string;
	    lastRun: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallationInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.version = source["version"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.lastRun = source["lastRun"];
	    }
	}
	export class Version {
	    label: string;
	    version: string;
	    url: string;
	    type: string;
	    tag: string;
	    risk: string;
	
	    static createFrom(source: any = {}) {
	        return new Version(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.version = source["version"];
	        this.url = source["url"];
	        this.type = source["type"];
	        this.tag = source["tag"];
	        this.risk = source["risk"];
	    }
	}

}

