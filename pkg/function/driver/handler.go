package driver

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails" // Pull in errdetails
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

type HandlerFunc func(grpcServer *grpc.Server, mux *runtime.ServeMux) http.HandlerFunc

func NewHandlerFunc(grpcServer *grpc.Server, mux *runtime.ServeMux) HandlerFunc {
	return func(grpcServer *grpc.Server, mux *runtime.ServeMux) http.HandlerFunc {
		return grpcHandlerFunc(grpcServer , mux)
	}
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}