package xrequestid

import (
	"fmt"

	"github.com/renstrom/shortuuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// DefaultXRequestIDKey is metadata key name for request ID
var DefaultXRequestIDKey = "x-request-id"

func HandleRequestID(ctx context.Context, validator requestIDValidator) string {
	var requestID string
	defer addToResponse(ctx, &requestID)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newRequestID()
	}

	header, ok := md[DefaultXRequestIDKey]
	if !ok || len(header) == 0 {
		requestID = newRequestID()
		return requestID
	}

	requestID = header[0]
	if requestID == "" {
		requestID = newRequestID()
		return requestID
	}

	if !validator(requestID) {
		requestID = newRequestID()
		return requestID
	}

	return requestID
}

func HandleRequestIDChain(ctx context.Context, validator requestIDValidator) string {
	var res string
	defer addToResponse(ctx, &res)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		res = newRequestID()
		return res
	}

	header, ok := md[DefaultXRequestIDKey]
	if !ok || len(header) == 0 {
		res = newRequestID()
		return res
	}

	requestID := header[0]
	if requestID == "" {
		res = newRequestID()
		return res
	}

	if !validator(requestID) {
		res = newRequestID()
		return res
	}

	newValue := fmt.Sprintf("%s,%s", requestID, newRequestID())
	addToResponse(ctx, &newValue)

	return newValue
}

func addToResponse(ctx context.Context, value *string) {
	if value == nil || *value == "" {
		x := newRequestID()
		value = &x
	}
	newHeader := metadata.Pairs(DefaultXRequestIDKey, *value)
	grpc.SendHeader(ctx, newHeader)
}

func newRequestID() string {
	return shortuuid.New()
}
