package worker

import (
	"context"
	db "github/dutt23/bank/db/sqlc"

	"github.com/hibiken/asynq"
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
	server := asynq.NewServer(redisOpts, asynq.Config{})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
