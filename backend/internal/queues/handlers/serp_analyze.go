package handlers_queues

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	tasks_queues "github.com/choirulanwar/textify/backend/internal/queues/tasks"
	"github.com/choirulanwar/textify/backend/pkg/helper"
	"github.com/chromedp/chromedp"
	"github.com/hibiken/asynq"
)

func (h *HandlerQueues) HandOnSerpExplorerRequestedTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks_queues.SerpExplorerTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("[+] Failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	settingInfo, err := helper.GetSettingInfo(h.db)
	if err != nil {
		return fmt.Errorf("[+] Failed to get setting database: %w", asynq.SkipRetry)
	}

	dir, err := os.MkdirTemp("", "chromedp-example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
		chromedp.Flag("headless", settingInfo.BrowserVisible),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// ctx, cancel := chromedp.NewContext(
	// 	taskCtx,
	// )
	// defer cancel()

	var buf []byte
	if err := chromedp.Run(taskCtx, elementScreenshot(`https://pkg.go.dev/`, `img.Homepage-logo`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(taskCtx, fullScreenshot(`https://brank.as/`, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Print("[======] Called from HandOnSerpExplorerRequestedTask")

	return nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}
