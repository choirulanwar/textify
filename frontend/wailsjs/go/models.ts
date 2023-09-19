export namespace model {
	
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
	    created_at: number;
	    updated_at: number;
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
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.browser_path = source["browser_path"];
	        this.browser_visible = source["browser_visible"];
	        this.session_google = source["session_google"];
	        this.session_pinterest = source["session_pinterest"];
	        this.proxy_url = source["proxy_url"];
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

