package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/choirulanwar/textify/backend/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	App *service.App
}

func New(a *service.App) *Service {
	return &Service{
		App: a,
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
				log.Println("Progress cancelled")
				return
			default:
				payload := ProgressPayload{Progress: progress}
				runtime.EventsEmit(s.App.Ctx, "progress:update", payload)
				log.Println("Progress:", progress)
				time.Sleep(100 * time.Millisecond)
			}
		}
		log.Println("Progress finished")
	}()

	runtime.EventsOnce(ctx, "progress:stop", func(optionData ...any) {
		fmt.Println("Stop progress from frontend")
		cancel()
	})
}
