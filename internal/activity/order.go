package activity

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"tutorial.io/tutorial-order/internal/models"
	giaekv1 "tutorial.io/tutorial-order/pkg/generated/tutorial.io/proto/v1"
)

type OrderActivities struct {
	snowflake *snowflake.Node
}

func NewOrderActivities(
	snowflake *snowflake.Node,
) OrderActivities {
	act := OrderActivities{
		snowflake: snowflake,
	}
	return act
}

func (act OrderActivities) IncrMemberOrder(ctx context.Context, in models.IncrMemberOrderRequest) (*giaekv1.Order, error) {
	return &giaekv1.Order{
		OrderId: act.snowflake.Generate().String(),
	}, nil
}
