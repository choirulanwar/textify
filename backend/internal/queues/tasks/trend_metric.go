package tasks_queues

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const TrendMetricQueueName = "textify:trend-metric"
const TypeOnTrendMetricRequested = "on:trend_metric_requested"

type TrendMetricTaskPayload struct {
	KeywordTrendExplorerID uint `json:"keyword_trend_explorer_id"`
	KeywordID              uint `json:"keyword_id"`
}

func NewOnTrendMetricRequestedTask(p *TrendMetricTaskPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(p)

	if err != nil {
		return nil, fmt.Errorf("[+] Failed to marshal task payload: %w", err)
	}

	return asynq.NewTask(
		TypeOnTrendMetricRequested,
		payload,
		asynq.Queue(TrendMetricQueueName),
		asynq.MaxRetry(5),
	), nil
}
