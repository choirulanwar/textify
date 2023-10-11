package keyword_trend_explorer_repository

import (
	keyword_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/helper"
	"github.com/choirulanwar/textify/backend/pkg/pagination"
	"gorm.io/gorm"
)

type KeywordTrendExplorerRepository struct {
	db *gorm.DB
}

func NewKeywordTrendExplorerRepository(db *gorm.DB) *KeywordTrendExplorerRepository {
	return &KeywordTrendExplorerRepository{
		db: db,
	}
}

func (r *KeywordTrendExplorerRepository) Create(keywordTrendExplorer *models.KeywordTrendExplorer) error {
	return r.db.Create(keywordTrendExplorer).Error
}

func (r *KeywordTrendExplorerRepository) Update(keywordTrendExplorer *models.KeywordTrendExplorer) error {
	return r.db.Save(keywordTrendExplorer).Error
}

func (r *KeywordTrendExplorerRepository) Delete(keywordTrendExplorer *models.KeywordTrendExplorer) error {
	keywordRepository := keyword_repository.NewKeywordRepository(r.db)
	if err := keywordRepository.DeleteManyByKeywordTrendExplorerID(keywordTrendExplorer.ID); err != nil {
		return err
	}

	err := r.db.Delete(keywordTrendExplorer).Error

	return err
}

func (r *KeywordTrendExplorerRepository) FindByID(id uint) (*models.KeywordTrendExplorer, error) {
	var keywordTrendExplorer models.KeywordTrendExplorer
	err := r.db.First(&keywordTrendExplorer, id).Error
	if err != nil {
		return nil, err
	}
	return &keywordTrendExplorer, nil
}

func (r *KeywordTrendExplorerRepository) FindAll() ([]*models.KeywordTrendExplorer, error) {
	var keywordTrendExplorers []*models.KeywordTrendExplorer
	err := r.db.Find(&keywordTrendExplorers).Error
	if err != nil {
		return nil, err
	}
	return keywordTrendExplorers, nil
}

func (r *KeywordTrendExplorerRepository) FindAllPaginated(pagination *pagination.Pagination) (*pagination.Pagination, error) {
	var keywordTrendExplorers []*models.KeywordTrendExplorer

	err := r.db.Scopes(helper.Paginate(keywordTrendExplorers, pagination, r.db)).Find(&keywordTrendExplorers).Error
	if err != nil {
		return nil, err
	}

	pagination.Rows = keywordTrendExplorers

	return pagination, nil
}
