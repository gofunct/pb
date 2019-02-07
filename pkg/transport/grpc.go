package transport

import (
	"github.com/gofunct/pb/pkg/transport/api"
	"github.com/gofunct/pb/pkg/transport/engine"
)

func ServeGrpc(servers ...api.Server) error {
	s := engine.New(
		engine.WithDefaultLogger(),
		engine.WithServers(
			servers...,
		),
	)
	return s.Serve()
}
