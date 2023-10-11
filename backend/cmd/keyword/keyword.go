package keyword

import (
	keyword_repository "github.com/choirulanwar/textify/backend/internal/repositories/keyword"
	"github.com/choirulanwar/textify/backend/pkg/resp"
	"github.com/choirulanwar/textify/backend/service"
	"github.com/hibiken/asynq"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

func (s *Service) GetKeywordsByExplorers(keywordTrendExplorerID uint) *resp.Response {
	keywordRepository := keyword_repository.NewKeywordRepository(s.App.DB)

	if result, err := keywordRepository.FindAllByExplorers(keywordTrendExplorerID); err != nil {
		return resp.Fail(err.Error())
	} else {
		return resp.Success(result)
	}
}

func (s *Service) DeleteKeyword(keywordID uint) *resp.Response {
	keywordRepository := keyword_repository.NewKeywordRepository(s.App.DB)
	keyword, err := keywordRepository.FindByID(keywordID)
	if err != nil {
		return resp.Fail(err.Error())
	}

	if err := keywordRepository.Delete(keyword); err != nil {
		return resp.Fail(err.Error())
	} else {
		runtime.EventsEmit(s.App.Ctx, "keyword:update", nil)
		return resp.Success(keyword)
	}
}
