package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	gatewayUserAgentKey = "grpcgateway-user-agent"
	gatewayClientIPKey  = "x-forwarded-for"
	grpcClientKey       = "grpc-client"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md : %+v\n", md)

		if gtusrAgnt := md.Get(gatewayUserAgentKey); len(gtusrAgnt) > 0 {
			mtdt.UserAgent = gtusrAgnt[0]
		}

		if gtClientIP := md.Get(gatewayClientIPKey); len(gtClientIP) > 0 {
			mtdt.ClientIP = gtClientIP[0]
		}

		if grpcAgent := md.Get(grpcClientKey); len(grpcAgent) > 0 {
			mtdt.UserAgent = grpcAgent[0]
		}

	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
