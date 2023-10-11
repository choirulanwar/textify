package keyword_trend_explorer

import (
	"time"

	tasks_queues "github.com/choirulanwar/textify/backend/internal/queues/tasks"
	keyword_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword"
	keyword_trend_explorer_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword_trend_explorer"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/pagination"
	"github.com/choirulanwar/textify/backend/pkg/resp"
	"github.com/choirulanwar/textify/backend/service"
	"github.com/hibiken/asynq"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type KeywordTrendExplorerRequestPayload struct {
	Query       string `json:"query,omitempty"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	Period      string `json:"period"`
	IncludeSerp bool   `json:"include_serp,omitempty"`
}

type Service struct {
	App         *service.App
	asynqClient *asynq.Client
}

func New(a *service.App) *Service {
	return &Service{
		App:         a,
		asynqClient: a.AsynqClient,
	}
}

func (s *Service) GetKeywordTrendExplorers(pagination pagination.Pagination) *resp.Response {
	keywordRearchRepository := keyword_trend_explorer_repository.NewKeywordTrendExplorerRepository(s.App.DB)
	if result, err := keywordRearchRepository.FindAllPaginated(&pagination); err != nil {
		return resp.Fail(err.Error())
	} else {
		return resp.Success(result)
	}
}

func (s *Service) GetKeywordTrendExplorer(ID uint) *resp.Response {
	keywordRearchRepository := keyword_trend_explorer_repository.NewKeywordTrendExplorerRepository(s.App.DB)
	if result, err := keywordRearchRepository.FindByID(ID); err != nil {
		return resp.Fail(err.Error())
	} else {
		return resp.Success(result)
	}
}

func (s *Service) GetKeywordsByExplorers(keywordTrendExplorerID uint, pagination pagination.Pagination) *resp.Response {
	keywordRepository := keyword_repository.NewKeywordRepository(s.App.DB)

	if result, err := keywordRepository.FindAllByExplorersPaginated(keywordTrendExplorerID, &pagination); err != nil {
		return resp.Fail(err.Error())
	} else {
		return resp.Success(result)
	}
}

func (s *Service) DeleteKeywordTrendExplorer(keywordTrendExplorerID uint) *resp.Response {
	keywordRearchRepository := keyword_trend_explorer_repository.NewKeywordTrendExplorerRepository(s.App.DB)
	keywordTrendExplorer, err := keywordRearchRepository.FindByID(keywordTrendExplorerID)
	if err != nil {
		return resp.Fail(err.Error())
	}

	if err := keywordRearchRepository.Delete(keywordTrendExplorer); err != nil {
		return resp.Fail(err.Error())
	} else {
		runtime.EventsEmit(s.App.Ctx, "keyword_trend_explorer:update", nil)
		return resp.Success(keywordTrendExplorer)
	}
}

func (s *Service) KeywordTrendExplorerRequest(p *KeywordTrendExplorerRequestPayload) *resp.Response {
	newData := &models.KeywordTrendExplorer{
		Keyword:   p.Query,
		Country:   p.Country,
		Language:  p.Language,
		Period:    p.Period,
		TaskModel: models.TaskModel{Status: "queue"},
	}

	keywordTrendExplorerRepository := keyword_trend_explorer_repository.NewKeywordTrendExplorerRepository(s.App.DB)

	err := keywordTrendExplorerRepository.Create(newData)
	if err != nil {
		return resp.Fail(err.Error())
	}

	task, err := tasks_queues.NewOnKeywordTrendExplorerRequestedTask(&tasks_queues.KeywordTrendExplorerTaskPayload{
		KeywordTrendExplorerID: newData.ID,
	})
	if err != nil {
		s.App.LogErrorf("[+] Failed to create keyword trend anlyzer task: %v", err)
		return resp.Fail(err.Error())
	}

	info, err := s.asynqClient.Enqueue(task, asynq.ProcessIn(5*time.Second))
	if err != nil {
		s.App.LogInfof(" [*] Failed to enqueue task. Type: %s. Queue: %s. MaxRetry: %d", info.Type, info.Queue, info.MaxRetry)
		return resp.Fail(err.Error())
	}
	s.App.LogInfof(" [*] Successfully enqueued task. Type: %s. Queue: %s. MaxRetry: %d", info.Type, info.Queue, info.MaxRetry)

	return resp.Success(newData)
}
