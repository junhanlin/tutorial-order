package server

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	tc "go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	shared "tutorial.io/tutorial-order/internal"
	"tutorial.io/tutorial-order/internal/models"
	"tutorial.io/tutorial-order/internal/workflow"
	giaekv1 "tutorial.io/tutorial-order/pkg/generated/tutorial.io/proto/v1"
)

type GrpcService interface {
	Register(server *grpc.Server)
}

type OrderServer struct {
	giaekv1.UnimplementedOrderServiceServer
	snowflake   *snowflake.Node
	db          *gorm.DB
	producer    sarama.SyncProducer
	restyClient *resty.Client
	redisClient *redis.Client
	temporal    tc.Client
}

func NewOrderServer(
	snowflakeNode *snowflake.Node,
	db *gorm.DB,
	producer sarama.SyncProducer,
	restyClient *resty.Client,
	redisClient *redis.Client,
	temporal tc.Client,
) *OrderServer {
	srv := &OrderServer{
		snowflake:   snowflakeNode,
		db:          db,
		producer:    producer,
		restyClient: restyClient,
		redisClient: redisClient,
		temporal:    temporal,
	}
	return srv

}

func (s OrderServer) CreateOrder(ctx context.Context, in *giaekv1.CreateOrderRequest) (*giaekv1.Order, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (s OrderServer) GetOrder(ctx context.Context, in *giaekv1.GetOrderRequest) (*giaekv1.Order, error) {
	orderId := s.snowflake.Generate().String()
	input := models.IncrMemberOrderRequest{
		OrderId: orderId,
	}
	options := tc.StartWorkflowOptions{
		ID:        fmt.Sprintf("GetOrder_%s", orderId),
		TaskQueue: shared.TemporalTaskQueueOrder,
	}

	we, err := s.temporal.ExecuteWorkflow(ctx, options, workflow.OrderWorkflow{}.Demo, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	return &giaekv1.Order{
		OrderId: orderId,
	}, nil
}
func (s OrderServer) ListOrder(ctx context.Context, in *giaekv1.ListOrderRequest) (*giaekv1.ListOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrder not implemented")
}

func (s OrderServer) Register(server *grpc.Server) {
	giaekv1.RegisterOrderServiceServer(server, s)
}
