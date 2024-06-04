package component

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"tutorial.io/tutorial-order/internal"
	"tutorial.io/tutorial-order/internal/server"
)

func NewGrpcServer(
	lc fx.Lifecycle,
	config *shared.Config,
	grpcServices []server.GrpcService,
) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	for _, svc := range grpcServices {
		svc.Register(grpcServer)
	}
	reflection.Register(grpcServer)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
				if err != nil {
					log.WithError(err).Fatal("gRPC server failed to listen")
				}
				log.Infof("grpc server listening at %v", lis.Addr())
				err = grpcServer.Serve(lis)
				if err != nil {
					log.WithError(err).Fatal("Error starting grpc server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
	return grpcServer
}
