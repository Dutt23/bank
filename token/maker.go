package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	Validate(token string) (*Payload, error)
}
