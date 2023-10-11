package models

import "gorm.io/datatypes"

type Serp struct {
	Position               uint                                   `gorm:"column:position;type:uint" json:"position"`
	URL                    string                                 `gorm:"column:url;type:text" json:"url"`
	MetaTitle              string                                 `gorm:"column:meta_title;type:text" json:"meta_title"`
	MetaDescription        string                                 `gorm:"column:meta_description;type:text" json:"meta_description"`
	ContentLength          uint                                   `gorm:"column:content_length;type:uint" json:"content_length"`
	ContentScore           uint                                   `gorm:"column:content_score;type:uint" json:"content_score"`
	OutlineTalkingPoint    string                                 `gorm:"column:outline_talking_point;type:text" json:"outline_talking_point"`
	Intent                 string                                 `gorm:"column:intent;type:text" json:"intent"`
	Links                  datatypes.JSONType[[]*Link]            `gorm:"column:links;type:text" json:"links"`
	DomainAuthority        uint                                   `gorm:"column:domain_authority;type:uint" json:"domain_authority"`
	PageAuthority          uint                                   `gorm:"column:page_authority;type:uint" json:"page_authority"`
	TotalImage             uint                                   `gorm:"column:total_image;type:uint" json:"total_image"`
	EstimatedMonthlyVisits uint                                   `gorm:"column:estimated_monthly_visits;type:uint" json:"estimated_monthly_visits"`
	Terms                  datatypes.JSONType[[]*Term]            `gorm:"column:terms;type:text" json:"terms"`
	RelatedQuestions       datatypes.JSONType[[]*RelatedQuestion] `gorm:"column:related_questions;type:text" json:"related_questions"`
}

type Link struct {
	URL        string `gorm:"column:url;type:text" json:"url"`
	Ref        string `gorm:"column:ref;type:text" json:"ref"`
	AnchorText string `gorm:"column:anchor_text;type:text" json:"anchor_text"`
}

type Term struct {
	Source string `gorm:"column:source;type:text" json:"source"`
	Term   string `gorm:"column:term;type:text" json:"term"`
	Total  uint   `gorm:"column:total;type:uint" json:"total"`
}

type RelatedQuestion struct {
	Question string `gorm:"column:question;type:text" json:"question"`
	Answer   string `gorm:"column:answer;type:text" json:"answer"`
}
