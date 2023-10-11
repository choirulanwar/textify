package queue

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

func Enqueue(asynqClient *asynq.Client, logger *logrus.Logger, task *asynq.Task, processIn time.Duration) error {
	if processIn == 0 {
		processIn = 5 * time.Second
	}

	info, err := asynqClient.Enqueue(task, asynq.ProcessIn(processIn))
	if err != nil {
		logger.Infof(" [*] Failed enqueued task. Type: %s. Queue: %s. MaxRetry: %d", info.Type, info.Queue, info.MaxRetry)
		return err
	}
	logger.Infof(" [*] Successfully enqueued task. Type: %s. Queue: %s. MaxRetry: %d", info.Type, info.Queue, info.MaxRetry)

	return nil
}
