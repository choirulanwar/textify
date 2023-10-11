package handlers_queues

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	google_trends_api "github.com/choirulanwar/textify/backend/internal/api/google_trends"
	pinterest_trends_api "github.com/choirulanwar/textify/backend/internal/api/pinterest_trends"
	tasks_queues "github.com/choirulanwar/textify/backend/internal/queues/tasks"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/queue"
	"github.com/hibiken/asynq"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/datatypes"
)

func (h *HandlerQueues) HandOnKeywordTrendExplorerRequestedTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks_queues.KeywordTrendExplorerTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("[+] Failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	keywordTrendExplorerInfo, err := h.keywordTrendExplorerRepository.FindByID(payload.KeywordTrendExplorerID)
	if err != nil {
		return err
	}

	keywordTrendExplorerInfo.Status = "in_progress"
	keywordTrendExplorerInfo.StartedAt = uint64(time.Now().Unix())

	if err := h.keywordTrendExplorerRepository.Update(keywordTrendExplorerInfo); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		h.handlerGoogleTrend(keywordTrendExplorerInfo)
	}()

	go func() {
		defer wg.Done()
		h.handlerPinterestTrend(keywordTrendExplorerInfo)
	}()

	wg.Wait()

	// keywordTrendExplorerInfo.Status = "completed"
	// keywordTrendExplorerInfo.CompletedAt = uint64(time.Now().Unix())

	// if err := h.keywordTrendExplorerRepository.Update(keywordTrendExplorerInfo); err != nil {
	// 	return err
	// }

	return nil
}

func (h *HandlerQueues) handlerGoogleTrend(keywordTrendExplorer *models.KeywordTrendExplorer) error {
	var period string

	switch keywordTrendExplorer.Period {
	case "1":
		period = "today 1-m"
	case "2":
		period = "today 12-m"
	case "3":
		period = "today 3-m"
	case "4":
		period = "today 5-y"
	default:
		period = "now 1-d"
	}

	country := keywordTrendExplorer.Country
	switch country {
	case "IT+ES+PT+GR+MT":
		country = "PT"
	case "DE+AT+CH":
		country = "DE"
	case "GB+IE":
		country = "GB"
	case "SE+DK+FI+NO":
		country = "SE"
	case "NL+BE+LU":
		country = "NL"
	case "PL+RO+HU+SK+CZ":
		country = "PL"
	case "MX+AR+CO+CL":
		country = "MX"
	case "AU+NZ":
		country = "AU"
	}

	gtClient := google_trends_api.NewGoogleTrend("", "0", country)

	return h.googleTrendsProcessor(gtClient, keywordTrendExplorer, period)
}

func (h *HandlerQueues) googleTrendsProcessor(gtClient *google_trends_api.GoogleTrendsAPIService, keywordTrendExplorer *models.KeywordTrendExplorer, period string) error {
	var firstLevelParams []*google_trends_api.Query
	firstLevelParams = append(firstLevelParams, &google_trends_api.Query{
		Q:    strings.ToLower(keywordTrendExplorer.Keyword),
		Time: period,
	})

	var secondLevelParams []*google_trends_api.Query

	for _, firstLevel := range firstLevelParams {
		widgets, err := gtClient.GetWidgets(firstLevel)
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.gtClient.GetWidgets.gProcessor.firstLevel: %v", time.Now().Format(time.RFC3339), err)
			return err
		}

		time.Sleep(3000 * time.Millisecond)

		relateds, err := gtClient.GetRelated(widgets)
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.gtClient.GetRelated.gProcessor.firstLevel: %v", time.Now().Format(time.RFC3339), err)
			return err
		}

		for _, ranked := range relateds.Default.Ranked {
			for index, keyword := range ranked.Keywords {
				newKeyword := &models.Keyword{
					KeywordTrendExplorerID: keywordTrendExplorer.ID,
					Keyword:                strings.ToLower(keyword.Query),
					Source:                 "google",
				}

				err := h.keywordRepository.Create(newKeyword)
				if err != nil {
					h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.CreateKeyword.gProcessor.firstLevel: %v", time.Now().Format(time.RFC3339), err)
					return err
				}

				secondLevelParams = append(secondLevelParams, &google_trends_api.Query{
					Q:    strings.ToLower(keyword.Query),
					Time: period,
				})

				runtime.EventsEmit(h.ctx, fmt.Sprintf("keyword_trend_explorer:update:%d", keywordTrendExplorer.ID), nil)

				task, err := tasks_queues.NewOnTrendMetricRequestedTask(&tasks_queues.TrendMetricTaskPayload{
					KeywordTrendExplorerID: keywordTrendExplorer.ID,
					KeywordID:              newKeyword.ID,
				})
				if err != nil {
					h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.EnqueueTrendMetric.gProcessor.firstLevel: %v", time.Now().Format(time.RFC3339), err)
					return err
				}

				processTaskIn := time.Duration(3+index) * time.Second
				err = queue.Enqueue(h.asynqClient, h.logger, task, processTaskIn)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, secondLevel := range secondLevelParams {
		widgets, err := gtClient.GetWidgets(secondLevel)
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.gtClient.GetWidgets.gProcessor.secondLevel: %v", time.Now().Format(time.RFC3339), err)
			return err
		}

		time.Sleep(3000 * time.Millisecond)

		relateds, err := gtClient.GetRelated(widgets)
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.gtClient.GetRelated.gProcessor.secondLevel: %v", time.Now().Format(time.RFC3339), err)
			return err
		}

		for _, ranked := range relateds.Default.Ranked {
			for index, keyword := range ranked.Keywords {
				newKeyword := &models.Keyword{
					KeywordTrendExplorerID: keywordTrendExplorer.ID,
					Keyword:                strings.ToLower(keyword.Query),
					Source:                 "google",
				}

				err := h.keywordRepository.Create(newKeyword)
				if err != nil {
					h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.CreateKeyword.gProcessor.secondLevel: %v", time.Now().Format(time.RFC3339), err)
					return err
				}

				runtime.EventsEmit(h.ctx, fmt.Sprintf("keyword_trend_explorer:update:%d", keywordTrendExplorer.ID), nil)

				task, err := tasks_queues.NewOnTrendMetricRequestedTask(&tasks_queues.TrendMetricTaskPayload{
					KeywordTrendExplorerID: keywordTrendExplorer.ID,
					KeywordID:              newKeyword.ID,
				})
				if err != nil {
					h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.EnqueueTrendMetric.gProcessor.secondLevel: %v", time.Now().Format(time.RFC3339), err)
					return err
				}

				processTaskIn := time.Duration(3+index) * time.Second
				err = queue.Enqueue(h.asynqClient, h.logger, task, processTaskIn)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (h *HandlerQueues) handlerPinterestTrend(keywordTrendExplorer *models.KeywordTrendExplorer) error {
	period := keywordTrendExplorer.Period
	switch period {
	case "0":
		period = "1"
	}

	now := time.Now()
	backdate := now.AddDate(0, 0, -4)

	ptClient := pinterest_trends_api.NewPinterestTrend(period, keywordTrendExplorer.Country, backdate.Format("2006-01-02"), h.db)

	return h.pinterestProcessor(ptClient, keywordTrendExplorer, keywordTrendExplorer.Keyword, period)
}

func (h *HandlerQueues) pinterestProcessor(ptClient *pinterest_trends_api.PinterestTrendsAPIService, keywordTrendExplorer *models.KeywordTrendExplorer, seedKeyword string, period string) error {
	ptTrends, err := ptClient.GetTrends(url.PathEscape(seedKeyword))
	if err != nil {
		h.logger.Debugf("Error::ptTrends.ptClient.GetTrends::%v", err)
		return err
	}

	for index, trend := range ptTrends.Values {
		newKeyword := &models.Keyword{
			KeywordTrendExplorerID: keywordTrendExplorer.ID,
			Keyword:                strings.ToLower(trend.Term),
			Source:                 "pinterest",
		}

		err := h.keywordRepository.Create(newKeyword)
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.CreateKeyword.pProcessor: %v", time.Now().Format(time.RFC3339), err)
			return err
		}

		runtime.EventsEmit(h.ctx, fmt.Sprintf("keyword_trend_explorer:update:%d", keywordTrendExplorer.ID), nil)

		task, err := tasks_queues.NewOnTrendMetricRequestedTask(&tasks_queues.TrendMetricTaskPayload{
			KeywordTrendExplorerID: keywordTrendExplorer.ID,
			KeywordID:              newKeyword.ID,
		})
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.EnqueueTrendMetric.pProcessor: %v", time.Now().Format(time.RFC3339), err)
			return err
		}
		processTaskIn := time.Duration(3+index) * time.Second
		err = queue.Enqueue(h.asynqClient, h.logger, task, processTaskIn)
		if err != nil {
			return err
		}

		relateds, err := ptClient.GetRelated(url.PathEscape(trend.Term))
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.ptClient.GetRelated.pProcessor: %v", time.Now().Format(time.RFC3339), err)
			return err
		}

		for _, related := range relateds {
			var trends []*models.Trend

			newKeyword := &models.Keyword{
				KeywordTrendExplorerID: keywordTrendExplorer.ID,
				Keyword:                strings.ToLower(related.Term),
				Source:                 "pinterest",
			}

			err := h.keywordRepository.Create(newKeyword)
			if err != nil {
				h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.CreateKeyword.pProcessor: %v", time.Now().Format(time.RFC3339), err)
				return err
			}

			runtime.EventsEmit(h.ctx, fmt.Sprintf("keyword_trend_explorer:update:%d", keywordTrendExplorer.ID), nil)

			reversedMetrics := make([]int, len(related.Counts))
			for i := len(related.Counts) - 1; i >= 0; i-- {
				reversedMetrics[len(related.Counts)-1-i] = related.Counts[i]
			}

			currentTime := time.Now().AddDate(0, 0, -3)

			for i := 0; i < len(reversedMetrics); i++ {
				trend := &models.Trend{
					Time:  uint64(currentTime.Unix()),
					Value: uint(reversedMetrics[i]),
				}

				trends = append(trends, trend)

				currentTime = currentTime.AddDate(0, 0, -7)
			}

			reversedTrends := make([]*models.Trend, len(trends))
			for i := len(trends) - 1; i >= 0; i-- {
				reversedTrends[len(trends)-1-i] = trends[i]
			}

			newKeyword.Trends = datatypes.NewJSONType(reversedTrends)

			err = h.keywordRepository.Update(newKeyword)
			if err != nil {
				h.logger.Debugf("[%s] [ERROR] on queues.handlers.keyword_trend_explorer.InsertTrends.pProcessor: %v", time.Now().Format(time.RFC3339), err)
				return err
			}

			runtime.EventsEmit(h.ctx, fmt.Sprintf("keyword_trend_explorer:update:%d", keywordTrendExplorer.ID), nil)
		}

		time.Sleep(3000 * time.Millisecond)
	}

	return nil
}
