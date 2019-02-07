package driver

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
)

type Function struct {
	HandlerFunc HandlerFunc
}

func NewFunction(handler HandlerFunc) Function {
	return Function{
		HandlerFunc: handler,
	}
}

func (f *Function) Handle (s *grpc.Server, mux *runtime.ServeMux) http.HandlerFunc {
	return f.HandlerFunc(s, mux)
}

