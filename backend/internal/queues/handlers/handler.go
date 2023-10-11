package handlers_queues

import (
	"context"

	keyword_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword"
	keyword_trend_explorer_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword_trend_explorer"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HandlerQueues struct {
	ctx                            context.Context
	logger                         *logrus.Logger
	db                             *gorm.DB
	asynqClient                    *asynq.Client
	keywordTrendExplorerRepository *keyword_trend_explorer_repository.KeywordTrendExplorerRepository
	keywordRepository              *keyword_repository.KeywordRepository
}

func NewHandlerQueues(ctx context.Context, logger *logrus.Logger, db *gorm.DB, asynqClient *asynq.Client, keywordTrendExplorerRepository *keyword_trend_explorer_repository.KeywordTrendExplorerRepository, keywordRepository *keyword_repository.KeywordRepository) *HandlerQueues {
	return &HandlerQueues{
		ctx:                            ctx,
		logger:                         logger,
		db:                             db,
		asynqClient:                    asynqClient,
		keywordTrendExplorerRepository: keywordTrendExplorerRepository,
		keywordRepository:              keywordRepository,
	}
}
