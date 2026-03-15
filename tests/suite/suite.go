package suite

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/glebateee/auto-inventory/internal/config"
	aiv1 "github.com/glebateee/auto-proto/gen/go/inventory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient aiv1.InventoryClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	//t.Log(os.Getwd())

	cfg := config.MustLoadByPath("../config/local.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: aiv1.NewInventoryClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(cfg.GRPCConfig.Host, strconv.Itoa(cfg.GRPCConfig.Port))
}
