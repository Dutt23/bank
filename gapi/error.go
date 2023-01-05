package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldVioldation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badrequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")
	status, err := statusInvalid.WithDetails(badrequest)

	if err != nil {
		return statusInvalid.Err()
	}

	return status.Err()
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized : %s", err)
}
