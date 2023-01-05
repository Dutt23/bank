package gapi

import (
	"context"
	"database/sql"
	"github/dutt23/bank/customvalidators"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/pb"
	"github/dutt23/bank/util"
	"log"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
	violations := validateCreateUserReqest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

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
	violations := validateUserLoginRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

	metadata := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
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

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserReqest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "failed to hash password : %s", err)
	}

	// nullString := sql.NullString{
	// 	String: req.GetFullname(),
	// 	Valid:  len(req.GetFullname()) != 0,
	// }

	// nullEmail := sql.NullString{
	// 	String: req.GetFullname(),
	// 	Valid:  len(req.GetEmail()) != 0,
	// }

	// nullPassword := sql.NullString{
	// 	String: req.GetFullname(),
	// 	Valid:  len(req.GetPassword()) != 0,
	// }

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

func validateCreateUserReqest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := customvalidators.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldVioldation("username", err))
	}

	if err := customvalidators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldVioldation("password", err))
	}

	if err := customvalidators.ValidateFullName(req.GetFullname()); err != nil {
		violations = append(violations, fieldVioldation("full_name", err))
	}

	if err := customvalidators.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldVioldation("email", err))
	}
	return
}

func validateUpdateUserReqest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := customvalidators.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldVioldation("username", err))
	}

	if err := customvalidators.ValidatePassword(req.GetPassword()); err != nil && len(req.GetPassword()) != 0 {
		violations = append(violations, fieldVioldation("password", err))
	}

	if err := customvalidators.ValidateFullName(req.GetFullname()); err != nil && len(req.GetFullname()) != 0 {
		violations = append(violations, fieldVioldation("full_name", err))
	}

	if err := customvalidators.ValidateEmail(req.GetEmail()); err != nil && len(req.GetEmail()) != 0 {
		violations = append(violations, fieldVioldation("email", err))
	}
	return
}

func validateUserLoginRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := customvalidators.ValidateUserName(req.GetUsername()); err != nil {
		violations = append(violations, fieldVioldation("username", err))
	}

	if err := customvalidators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldVioldation("password", err))
	}
	return
}
