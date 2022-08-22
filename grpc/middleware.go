package grpc

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func loggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	// log.InfoWithFields(
	// 	"Request_Response",
	// 	log.Fields{
	// 		"method": info.FullMethod,
	// 		"req":    common.PrintJson(req),
	// 		"resp":   common.PrintJson(resp),
	// 	},
	// )
	return resp, err
}

func recoverInterceptor(p interface{}) (err error) {
	// log.ErrorWithFields(
	// 	"recover",
	// 	log.Fields{
	// 		"err": p,
	// 	},
	// )
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}
