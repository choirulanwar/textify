package queue

import (
	"github.com/hibiken/asynq"
)

var Server *asynq.Server

func NewAsynqServer(
	redisOpt asynq.RedisClientOpt,
	concurrency int,
	queues map[string]int,
) *asynq.Server {

	Server = asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: concurrency,
			Queues:      queues,
		})

	return Server
}
