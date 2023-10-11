package command

import (
	"context"
	"fmt"
	"time"

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

type ProgressPayload struct {
	Progress int `json:"progress"`
}

func (s *Service) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (s *Service) ProgressTracker() {
	ctx, cancel := context.WithCancel(s.App.Ctx)

	go func() {
		defer cancel()

		for progress := 0; progress <= 100; progress++ {
			select {
			case <-ctx.Done():
				return
			default:
				payload := ProgressPayload{Progress: progress}
				runtime.EventsEmit(s.App.Ctx, "progress:update", payload)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	runtime.EventsOnce(ctx, "progress:stop", func(optionData ...any) {
		cancel()
	})
}
