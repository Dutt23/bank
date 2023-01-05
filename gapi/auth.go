package gapi

import (
	"context"
	"fmt"
	"github/dutt23/bank/token"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	auth := md.Get(authorizationHeader)

	if len(auth) == 0 {
		return nil, fmt.Errorf("missing auth header")
	}

	authHeader := auth[0]

	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization type")
	}

	authType := strings.ToLower(fields[0])

	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsuported auth type: %s", authType)
	}

	payload, err := server.tokenMaker.Validate(fields[1])

	if err != nil {
		return nil, fmt.Errorf("invalid access token ")
	}

	return payload, nil
}
