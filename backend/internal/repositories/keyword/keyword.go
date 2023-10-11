package keyword_repository

import (
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/helper"
	"github.com/choirulanwar/textify/backend/pkg/pagination"
	"gorm.io/gorm"
)

type KeywordRepository struct {
	db *gorm.DB
}

func NewKeywordRepository(db *gorm.DB) *KeywordRepository {
	return &KeywordRepository{
		db: db,
	}
}

func (r *KeywordRepository) Create(keyword *models.Keyword) error {
	err := r.db.Create(keyword).Error
	if err != nil {
		return err
	}
	return nil
	// err := r.db.Where(models.Keyword{KeywordTrendExplorerID: keyword.KeywordTrendExplorerID, Keyword: keyword.Keyword, Source: keyword.Source}).Assign(keyword).FirstOrCreate(keyword).Error
	// if err != nil {
	// 	return err
	// }
	// return nil
}

// func (r *KeywordRepository) Create(keyword *models.Keyword) error {
// 	existingKeyword := models.Keyword{}
// 	err := r.db.Where("keyword_trend_explorer_id = ? AND source = ? AND keyword = ?",
// 		keyword.KeywordTrendExplorerID, keyword.Source, keyword.Keyword).First(&existingKeyword).Error

// 	if err != nil {
// 		if err := r.db.Create(keyword).Error; err != nil {
// 			return err
// 		}
// 	} else {
// 		updatedKeyword := keyword
// 		if err := r.db.Model(&existingKeyword).Assign(updatedKeyword).Save(&existingKeyword).Error; err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (r *KeywordRepository) Update(keyword *models.Keyword) error {
	return r.db.Save(keyword).Error
}

func (r *KeywordRepository) Delete(keyword *models.Keyword) error {
	return r.db.Delete(keyword).Error
}

func (r *KeywordRepository) DeleteManyByKeywordTrendExplorerID(keywordTrendExplorerID uint) error {
	err := r.db.Where("keyword_trend_explorer_id = ?", keywordTrendExplorerID).Delete(&models.Keyword{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *KeywordRepository) FindByID(id uint) (*models.Keyword, error) {
	var keyword models.Keyword
	err := r.db.First(&keyword, id).Error
	if err != nil {
		return nil, err
	}
	return &keyword, nil
}

func (r *KeywordRepository) FindAll() ([]*models.Keyword, error) {
	var keywords []*models.Keyword
	err := r.db.Find(&keywords).Error
	if err != nil {
		return nil, err
	}
	return keywords, nil
}

func (r *KeywordRepository) FindAllByExplorers(keywordTrendExplorerID uint) ([]*models.Keyword, error) {
	var keywords []*models.Keyword
	err := r.db.Where("keyword_trend_explorer_id = ?", keywordTrendExplorerID).Find(&keywords).Error
	if err != nil {
		return nil, err
	}
	return keywords, nil
}

func (r *KeywordRepository) FindAllByExplorersPaginated(keywordTrendExplorerID uint, pagination *pagination.Pagination) (*pagination.Pagination, error) {
	var keywords []*models.Keyword

	err := r.db.Scopes(helper.Paginate(keywords, pagination, r.db, func(db *gorm.DB) *gorm.DB {
		return db.Where("keyword_trend_explorer_id = ?", keywordTrendExplorerID)
	})).Find(&keywords).Error
	if err != nil {
		return nil, err
	}

	pagination.Rows = keywords

	return pagination, nil
}
