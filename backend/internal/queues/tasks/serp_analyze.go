package tasks_queues

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const SerpExplorerQueueName = "textify:serp-Explorer"
const TypeOnSerpExplorerRequested = "on:serp_Explorer_requested"

type SerpExplorerTaskPayload struct {
	Query   string `json:"query,omitempty"`
	Country string `json:"country"`
	Period  string `json:"period"`
}

func NewOnSerpExplorerRequestedTask(p *SerpExplorerTaskPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(p)

	if err != nil {
		return nil, fmt.Errorf("[+] Failed to marshal task payload: %w", err)
	}

	return asynq.NewTask(
		TypeOnSerpExplorerRequested,
		payload,
		asynq.Queue(SerpExplorerQueueName),
		asynq.MaxRetry(5),
	), nil
}
