export interface IPinterestTrends {
	values: Value[];
	endDate: string;
}

export interface Value {
	mom_change: MomChange;
	normalizedCount: number;
	wow_change: WowChange;
	term: string;
	affinity: any;
	yoy_change: YoyChange;
	seasonality_score: number;
	searchCount: number;
	reverseRank: number;
}

export interface MomChange {
	index: number;
	value: number;
}

export interface WowChange {
	index: number;
	value: number;
}

export interface YoyChange {
	index: number;
	value: number;
}

export interface IPinterestRelated {
	term: string;
	counts: number[];
}

interface MetricCount {
	count: number;
	date: string;
	normalizedCount: number;
}

export interface IPinterestMetric {
	counts: MetricCount[];
	term: string;
}
