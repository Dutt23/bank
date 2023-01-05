package worker

import (
	"context"
	db "github/dutt23/bank/db/sqlc"

	"github.com/hibiken/asynq"
)

const (
	CriticalQueue = "critical"
	defaultQueue  = "default"
)

type Proccessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store) Proccessor {

	server := asynq.NewServer(redisOpts, asynq.Config{
		Queues: map[string]int{
			CriticalQueue: 10,
			defaultQueue:  3,
			"low":         1,
		},
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
