package server

import (
	aiv1 "github.com/glebateee/auto-proto/gen/go/inventory"
	"google.golang.org/grpc"
)

type serverApi struct {
	aiv1.UnimplementedInventoryServer
}

func Register(
	srv *grpc.Server,
) {
	aiv1.RegisterInventoryServer(srv, &serverApi{})
}
