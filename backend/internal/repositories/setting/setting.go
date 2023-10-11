package setting_repository

import (
	"github.com/choirulanwar/textify/backend/models"
	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{
		db: db,
	}
}

func (r *SettingRepository) Update(Setting *models.Setting) error {
	return r.db.Where("id = ?", 1).Save(Setting).Error
}

func (r *SettingRepository) Find() (*models.Setting, error) {
	var Setting models.Setting
	err := r.db.Where("id = ?", 1).First(&Setting).Error
	if err != nil {
		return nil, err
	}
	return &Setting, nil
}
