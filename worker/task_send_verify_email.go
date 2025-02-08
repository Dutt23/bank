package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/util"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const taskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal task payload for sending email")
	}
	task := asynq.NewTask(taskSendVerifyEmail, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enqueue email task")
	}

	log.Info().Str("type", task.Type()).Bytes("payload", taskInfo.Payload).Str("queue", taskInfo.Queue).Int("max_retry", taskInfo.MaxRetry).Msg("Enqueued task")
	return nil
}

func (processor RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unable to un-marshal json for task %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("unable to get user record %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user %w", err)
	}

	// TODO: send email
	ve, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username: user.Username,
		Email: user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("Failed to create verify email %w", err)
	}

	subject := "Welcome to simple bank"
	verifyUrl := fmt.Sprintf("http://simple-bank.org?id=%d&secret_code=%s", ve.ID, ve.SecretCode)
	content := fmt.Sprintf(`Hello %s <br />
	Thankyou for registering with us
	Please <a href="%s">click here </a> to verify your email address 
	`, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil , nil);
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed task")
	return nil
}

func (processor RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(taskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
