package tasks_queues

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const KeywordTrendExplorerQueueName = "textify:keyword-trend-explorer"
const TypeOnKeywordTrendExplorerRequested = "on:keyword_trend_explorer_requested"

type KeywordTrendExplorerTaskPayload struct {
	KeywordTrendExplorerID uint `json:"keyword_trend_explorer_id"`
}

func NewOnKeywordTrendExplorerRequestedTask(p *KeywordTrendExplorerTaskPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(p)

	if err != nil {
		return nil, fmt.Errorf("[+] Failed to marshal task payload: %w", err)
	}

	return asynq.NewTask(
		TypeOnKeywordTrendExplorerRequested,
		payload,
		asynq.Queue(KeywordTrendExplorerQueueName),
		asynq.MaxRetry(5),
	), nil
}
