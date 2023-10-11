package handlers_queues

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	google_trends_api "github.com/choirulanwar/textify/backend/internal/api/google_trends"
	pinterest_trends_api "github.com/choirulanwar/textify/backend/internal/api/pinterest_trends"
	tasks_queues "github.com/choirulanwar/textify/backend/internal/queues/tasks"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/hibiken/asynq"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/datatypes"
)

func (h *HandlerQueues) HandOnTrendMetricRequestedTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks_queues.TrendMetricTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("[+] Failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	keywordTrendExplorerInfo, err := h.keywordTrendExplorerRepository.FindByID(payload.KeywordTrendExplorerID)
	if err != nil {
		return err
	}

	keywordInfo, err := h.keywordRepository.FindByID(payload.KeywordID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if keywordInfo.Source == "google" {
			h.handlerGoogleTrendsMetric(keywordTrendExplorerInfo, keywordInfo)
		}
	}()

	go func() {
		defer wg.Done()
		if keywordInfo.Source == "pinterest" {
			h.handlerPinterestTrendsMetric(keywordTrendExplorerInfo, keywordInfo)
		}
	}()

	wg.Wait()

	// keywordTrendExplorerInfo.Status = "completed"
	// keywordTrendExplorerInfo.CompletedAt = uint64(time.Now().Unix())

	// if err := h.keywordTrendExplorerRepository.Update(keywordTrendExplorerInfo); err != nil {
	// 	return err
	// }

	return nil
}

func (h *HandlerQueues) handlerGoogleTrendsMetric(keywordTrendExplorer *models.KeywordTrendExplorer, keyword *models.Keyword) error {
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

	return h.googleTrendsMetricProcessor(gtClient, keywordTrendExplorer, keyword, period)
}

func (h *HandlerQueues) googleTrendsMetricProcessor(gtClient *google_trends_api.GoogleTrendsAPIService, keywordTrendExplorer *models.KeywordTrendExplorer, keyword *models.Keyword, period string) error {
	params := &google_trends_api.Query{
		Q:    keyword.Keyword,
		Time: "today 12-m",
	}

	widgets, err := gtClient.GetWidgets(params)
	if err != nil {
		h.logger.Debugf("[%s] [ERROR] on queues.handlers.trend_metric.gtClient.GetWidgets.gTrendsMetricProcessor: %v", time.Now().Format(time.RFC3339), err)
		return err
	}

	time.Sleep(3000 * time.Millisecond)

	interests, err := gtClient.GetInterest(widgets)
	if err != nil {
		h.logger.Debugf("[%s] [ERROR] on queues.handlers.trend_metric.gtClient.GetInterest.gTrendsMetricProcessor: %v", time.Now().Format(time.RFC3339), err)
		return err
	}

	var trends []*models.Trend

	for _, interest := range interests.Default.TimelineData {
		t, err := strconv.ParseUint(interest.Time, 10, 64)
		if err != nil {
			h.logger.Debugf("[%s] [ERROR] on queues.handlers.trend_metric.parseInterestTime.gTrendsMetricProcessor: %v", time.Now().Format(time.RFC3339), err)
			return err
		}

		trend := &models.Trend{
			Time:  t,
			Value: uint(interest.Value[0]),
		}

		trends = append(trends, trend)
	}

	// var monthlyTrends []*models.Trend
	// groupByMonth := make(map[string]uint)
	// for _, trend := range trends {
	// 	timestamp := time.Unix(int64(trend.Time), 0)
	// 	month := timestamp.Format("January 2006")

	// 	groupByMonth[month] += trend.Value
	// }

	// for month, value := range groupByMonth {
	// 	timestamp, _ := time.Parse("January 2006", month)
	// 	unixTime := uint64(timestamp.Unix())

	// 	if value > 100 {
	// 		value = 100
	// 	}

	// 	monthlyTrends = append(monthlyTrends, &models.Trend{
	// 		Time:  unixTime,
	// 		Value: value,
	// 	})
	// }

	keyword.Trends = datatypes.NewJSONType(trends)
	err = h.keywordRepository.Update(keyword)
	if err != nil {
		h.logger.Debugf("[%s] [ERROR] on queues.handlers.trend_metric.InsertTrends.gTrendsMetricProcessor: %v", time.Now().Format(time.RFC3339), err)
		return err
	}

	time.Sleep(3000 * time.Millisecond)

	runtime.EventsEmit(h.ctx, fmt.Sprintf("keyword_trend_explorer:update:%d", keywordTrendExplorer.ID), nil)

	return nil
}

func (h *HandlerQueues) handlerPinterestTrendsMetric(keywordTrendExplorer *models.KeywordTrendExplorer, keyword *models.Keyword) error {
	var period string

	switch keywordTrendExplorer.Period {
	case "0":
		period = "1"
	default:
		period = keywordTrendExplorer.Period
	}

	now := time.Now()
	backdate := now.AddDate(0, 0, -4)

	ptClient := pinterest_trends_api.NewPinterestTrend(period, keywordTrendExplorer.Country, backdate.Format("2006-01-02"), h.db)

	return h.pinterestTrendsMetricProcessor(ptClient, keywordTrendExplorer, keyword, period)
}
func (h *HandlerQueues) pinterestTrendsMetricProcessor(ptClient *pinterest_trends_api.PinterestTrendsAPIService, keywordTrendExplorer *models.KeywordTrendExplorer, keyword *models.Keyword, period string) error {
	ptMetric, err := ptClient.GetMetric(url.PathEscape(keyword.Keyword))
	if err != nil {
		h.logger.Debugf("[%s] [ERROR] on queues.handlers.trend_metric.ptClient.GetMetric.pTrendsMetricProcessor: %v", time.Now().Format(time.RFC3339), err)
		return err
	}

	var trends []*models.Trend

	sourceDateFormat := "2006-01-02"

	for _, metrics := range ptMetric {
		for _, metric := range metrics.Counts {
			date, err := time.Parse(sourceDateFormat, metric.Date)
			if err != nil {
				h.logger.Debugf("[%s] [ERROR] on queues.handlers.trend_metric.parseMetricDate.pTrendsMetricProcessor: %v", time.Now().Format(time.RFC3339), err)
				return err
			}
			unixTime := date.Unix()

			trend := &models.Trend{
				Time:  uint64(unixTime),
				Value: uint(metric.Count),
			}

			trends = append(trends, trend)
		}
	}

	// var monthlyTrends []*models.Trend
	// groupByMonth := make(map[string]uint)
	// for _, trend := range trends {
	// 	timestamp := time.Unix(int64(trend.Time), 0)
	// 	month := timestamp.Format("January 2006")

	// 	groupByMonth[month] += trend.Value
	// }

	// for month, value := range groupByMonth {
	// 	timestamp, _ := time.Parse("January 2006", month)
	// 	unixTime := uint64(timestamp.Unix())

	// 	if value > 100 {
	// 		value = 100
	// 	}

	// 	monthlyTrends = append(monthlyTrends, &models.Trend{
	// 		Time:  unixTime,
	// 		Value: value,
	// 	})
	// }

	keyword.Trends = datatypes.NewJSONType(trends)
	err = h.keywordRepository.Update(keyword)
	if err != nil {
		h.logger.Debugf("[%s] [ERROR] on queues.handlers.trend_metric.InsertTrends.pTrendsMetricProcessor: %v", time.Now().Format(time.RFC3339), err)
		return err
	}

	time.Sleep(3000 * time.Millisecond)

	runtime.EventsEmit(h.ctx, fmt.Sprintf("keyword_trend_explorer:update:%d", keywordTrendExplorer.ID), nil)

	return nil
}
