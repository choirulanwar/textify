export namespace file {
	
	export class FileInfo {
	    name: string;
	    size: number;
	    isDir: boolean;
	    modTime: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.isDir = source["isDir"];
	        this.modTime = source["modTime"];
	    }
	}

}

export namespace keyword_trend_explorer {
	
	export class KeywordTrendExplorerRequestPayload {
	    query?: string;
	    country: string;
	    language: string;
	    period: string;
	    include_serp?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new KeywordTrendExplorerRequestPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.country = source["country"];
	        this.language = source["language"];
	        this.period = source["period"];
	        this.include_serp = source["include_serp"];
	    }
	}

}

export namespace models {
	
	export class FeedBack {
	    title?: string;
	    body?: string;
	    labels?: string[];
	    assignee?: string;
	    state?: string;
	    state_reason?: string;
	    milestone?: number;
	    assignees?: string[];
	
	    static createFrom(source: any = {}) {
	        return new FeedBack(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.body = source["body"];
	        this.labels = source["labels"];
	        this.assignee = source["assignee"];
	        this.state = source["state"];
	        this.state_reason = source["state_reason"];
	        this.milestone = source["milestone"];
	        this.assignees = source["assignees"];
	    }
	}
	export class FeedbackReq {
	    issue_type?: number;
	    title?: string;
	    body?: string;
	    version?: string;
	
	    static createFrom(source: any = {}) {
	        return new FeedbackReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issue_type = source["issue_type"];
	        this.title = source["title"];
	        this.body = source["body"];
	        this.version = source["version"];
	    }
	}
	export class Setting {
	    id: number;
	    browser_path: string;
	    browser_visible: boolean;
	    session_google: string;
	    session_pinterest: string;
	    proxy_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Setting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.browser_path = source["browser_path"];
	        this.browser_visible = source["browser_visible"];
	        this.session_google = source["session_google"];
	        this.session_pinterest = source["session_pinterest"];
	        this.proxy_url = source["proxy_url"];
	    }
	}

}

export namespace pagination {
	
	export class Pagination {
	    limit?: number;
	    page?: number;
	    sort?: string;
	    total_rows?: number;
	    total_pages?: number;
	    rows?: any;
	
	    static createFrom(source: any = {}) {
	        return new Pagination(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.limit = source["limit"];
	        this.page = source["page"];
	        this.sort = source["sort"];
	        this.total_rows = source["total_rows"];
	        this.total_pages = source["total_pages"];
	        this.rows = source["rows"];
	    }
	}

}

export namespace resp {
	
	export class Response {
	    code: number;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}

}

