package google_trends_api

type comparison struct {
	Keyword                string `json:"keyword"`
	Geo                    string `json:"geo"`
	Backend                string `json:"backend,omitempty"`
	Property               string `json:"property,omitempty"`
	Timezone               int    `json:"tz,omitempty"`
	HL                     string `json:"hl,omitempty"`
	Category               string `json:"category,omitempty"`
	GranularTimeResolution bool   `json:"granularTimeResolution,omitempty"`
	Time                   string `json:"time,omitempty"`
}

type GTExplorerResp struct {
	Widgets []*explorerWidget `json:"widgets"`
}

type explorerWidget struct {
	Token   string          `json:"token"`
	Type    string          `json:"type" bson:"type"`
	Title   string          `json:"title" bson:"title"`
	ID      string          `json:"id" bson:"id"`
	Request *WidgetResponse `json:"request" bson:"request"`
}

type WidgetResponse struct {
	Geo                interface{}             `json:"geo,omitempty" bson:"geo"`
	Time               string                  `json:"time,omitempty" bson:"time"`
	Resolution         string                  `json:"resolution,omitempty" bson:"resolution"`
	Locale             string                  `json:"locale,omitempty" bson:"locale"`
	Restriction        widgetComparisonItem    `json:"restriction" bson:"restriction"`
	CompItem           []*widgetComparisonItem `json:"comparisonItem" bson:"comparison_item"`
	RequestOpt         requestOptions          `json:"requestOptions" bson:"request_option"`
	KeywordType        string                  `json:"keywordType" bson:"keyword_type"`
	Metric             []string                `json:"metric" bson:"metric"`
	Language           string                  `json:"language" bson:"language"`
	TrendinessSettings map[string]string       `json:"trendinessSettings" bson:"trendiness_settings"`
	DataMode           string                  `json:"dataMode,omitempty" bson:"data_mode"`
	UserConfig         map[string]string       `json:"userConfig,omitempty" bson:"user_config"`
	UserCountryCode    string                  `json:"userCountryCode,omitempty" bson:"user_country_code"`
}

type widgetComparisonItem struct {
	Geo                             map[string]string   `json:"geo,omitempty" bson:"geo"`
	Time                            string              `json:"time,omitempty" bson:"time"`
	ComplexKeywordsRestriction      keywordsRestriction `json:"complexKeywordsRestriction,omitempty" bson:"complex_keywords_restriction"`
	OriginalTimeRangeForExplorerURL string              `json:"originalTimeRangeForExplorerUrl,omitempty" bson:"original_time_range_for_explorer_url"`
}

type keywordsRestriction struct {
	Keyword []*keywordRestriction `json:"keyword" bson:"keyword"`
}

type keywordRestriction struct {
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}

type requestOptions struct {
	Property string `json:"property" bson:"property"`
	Backend  string `json:"backend" bson:"backend"`
	Category int    `json:"category" bson:"category"`
}

type GTRelatedResp struct {
	Default relatedList `json:"default" bson:"default"`
}

type relatedList struct {
	Ranked []*rankedList `json:"rankedList" bson:"ranked"`
}

type rankedList struct {
	Keywords []*rankedKeyword `json:"rankedKeyword" bson:"keywords"`
}

type rankedKeyword struct {
	Query          string       `json:"query,omitempty" bson:"query"`
	Topic          KeywordTopic `json:"topic,omitempty" bson:"topic"`
	Value          int          `json:"value" bson:"value"`
	FormattedValue string       `json:"formattedValue" bson:"formatted_value"`
	HasData        bool         `json:"hasData" bson:"has_data"`
	Link           string       `json:"link" bson:"link"`
}

type KeywordTopic struct {
	Mid   string `json:"mid" bson:"mid"`
	Title string `json:"title" bson:"title"`
	Type  string `json:"type" bson:"type"`
}

type GTInterestResp struct {
	Default multiline `json:"default" bson:"default"`
}

type multiline struct {
	TimelineData []*Timeline `json:"timelineData" bson:"timeline_data"`
}

type Timeline struct {
	Time              string   `json:"time" bson:"time"`
	FormattedTime     string   `json:"formattedTime" bson:"formatted_time"`
	FormattedAxisTime string   `json:"formattedAxisTime" bson:"formatted_axis_time"`
	Value             []int    `json:"value" bson:"value"`
	HasData           []bool   `json:"hasData" bson:"has_data"`
	FormattedValue    []string `json:"formattedValue" bson:"formatted_value"`
}
