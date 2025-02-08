package worker

import (
	"context"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/mail"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
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
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) Proccessor {

	server := asynq.NewServer(redisOpts, asynq.Config{
		Queues: map[string]int{
			CriticalQueue: 10,
			defaultQueue:  3,
			"low":         1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}
