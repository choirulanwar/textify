package pinterest_trends_api

type value struct {
	MomChange        momChange   `json:"mom_change"`
	NormalizedCount  int         `json:"normalizedCount"`
	WowChange        wowChange   `json:"wow_change"`
	Term             string      `json:"term"`
	Affinity         interface{} `json:"affinity"`
	YoyChange        yoyChange   `json:"yoy_change"`
	SeasonalityScore float64     `json:"seasonality_score"`
	SearchCount      int         `json:"searchCount"`
	ReverseRank      int         `json:"reverseRank"`
}

type momChange struct {
	Index int     `json:"index"`
	Value float64 `json:"value"`
}

type wowChange struct {
	Index int     `json:"index"`
	Value float64 `json:"value"`
}

type yoyChange struct {
	Index int     `json:"index"`
	Value float64 `json:"value"`
}

type PTRelatedResp struct {
	Term   string `json:"term"`
	Counts []int  `json:"counts"`
}

type PTExplorerResp struct {
	Values  []value `json:"values"`
	EndDate string  `json:"endDate"`
}

type metricCount struct {
	Count           int    `json:"count"`
	Date            string `json:"date"`
	NormalizedCount int    `json:"normalizedCount"`
}

type PTMetricResp struct {
	Term   string        `json:"term"`
	Counts []metricCount `json:"counts"`
}
