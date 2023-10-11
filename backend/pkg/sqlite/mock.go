package sqlite

import (
	"github.com/choirulanwar/textify/backend/models"
	"gorm.io/gorm"
)

func MockSetting(db *gorm.DB) {
	db.Create(&models.Setting{
		BrowserPath:      "",
		BrowserVisible:   false,
		SessionGoogle:    "",
		SessionPinterest: "",
	})
}

func MockKeyword(db *gorm.DB) {

}

func MockTrend(db *gorm.DB) {

}

func MockTagList(db *gorm.DB) {

}
