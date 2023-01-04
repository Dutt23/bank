package gapi

import (
	"context"
	"database/sql"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/pb"
	"github/dutt23/bank/util"
	"log"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		Fullname:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.ChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "failed to hash password : %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		Email:          req.GetEmail(),
		FullName:       req.GetFullname(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())

			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "Username already exists : %s", err)
			}
		}
		return nil, status.Errorf(codes.AlreadyExists, "Failed to create user : %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: createUser(user),
	}

	return rsp, nil
}

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Error logging in: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error logging in: %s", err)
	}

	token, payload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error logging in: %s", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error logging in: %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
		Username:     req.Username,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error logging in: %s", err)
	}

	resp := &pb.LoginUserResponse{
		User:                  createUser(user),
		SessionId:             session.ID.String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		AccessToken:           token,
		AccessTokenExpiredAt:  timestamppb.New(payload.ExpiredAt),
	}

	return resp, nil
}
