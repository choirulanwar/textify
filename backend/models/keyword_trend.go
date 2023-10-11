package models

type KeywordTrendExplorer struct {
	BaseModel
	Keyword  string `gorm:"column:keyword;type:text" json:"keyword"`
	Country  string `gorm:"column:country;type:text" json:"country"`
	Language string `gorm:"column:language;type:text" json:"language"`
	Period   string `gorm:"column:period;type:text" json:"period"`
	TaskModel
	Results []*Keyword `gorm:"foreignKey:KeywordTrendExplorerID" json:"results"`
}
