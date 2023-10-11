package models

import "gorm.io/datatypes"

type Keyword struct {
	BaseModel
	KeywordTrendExplorerID uint                         `gorm:"column:keyword_trend_explorer_id;type:uint" json:"keyword_trend_explorer_id"`
	Keyword                string                       `gorm:"column:keyword;type:text" json:"keyword"`
	Source                 string                       `gorm:"column:source;type:text" json:"source"`
	Volume                 uint                         `gorm:"column:volume;type:uint" json:"volume"`
	CPC                    float64                      `gorm:"column:cpc;type:double" json:"cpc"`
	KD                     float64                      `gorm:"column:keyword_difficulty;type:double" json:"keyword_difficulty"`
	Serps                  []*Serp                      `gorm:"column:serps;type:text" json:"serps"`
	Trends                 datatypes.JSONType[[]*Trend] `gorm:"column:trends;type:text" json:"trends"`
	TotalSERPResult        uint                         `gorm:"column:total_serp_result;type:uint" json:"total_serp_result"`
}

type Trend struct {
	Time  uint64 `gorm:"column:time;type:bigint unsigned;" json:"time"`
	Value uint   `gorm:"column:value;type:uint" json:"value"`
}
