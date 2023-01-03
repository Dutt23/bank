package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func newJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size must be atleast %d characters but was %d", minSecretKeySize, len(secretKey))
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) Validate(token string) (*Payload, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("signing method is not valid here")
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyfunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)

		if ok && errors.Is(verr.Inner, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}

		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, errors.New("invalid token")
	}

	return payload, nil
}
