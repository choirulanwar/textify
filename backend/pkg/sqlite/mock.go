package sqlite

import (
	"time"

	"github.com/choirulanwar/textify/backend/model"
	"gorm.io/gorm"
)

func MockSetting(db *gorm.DB) {
	db.Create(&model.Setting{
		BrowserPath:      "",
		BrowserVisible:   false,
		SessionGoogle:    "",
		SessionPinterest: "",
		BaseModel: model.BaseModel{
			CreatedAt: uint64(time.Now().Unix()),
			UpdatedAt: uint64(time.Now().Unix()),
		},
	})
}

func MockTagList(db *gorm.DB) {

}
