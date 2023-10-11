package queue

import (
	"github.com/hibiken/asynq"
)

var Client *asynq.Client

func NewAsynqClient(redisOpt asynq.RedisClientOpt) *asynq.Client {
	return asynq.NewClient(redisOpt)
}
