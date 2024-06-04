package workflow

import (
	"fmt"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
	"tutorial.io/tutorial-order/internal/activity"
	"tutorial.io/tutorial-order/internal/models"
	giaekv1 "tutorial.io/tutorial-order/pkg/generated/tutorial.io/proto/v1"
)

type OrderWorkflow struct {
	activity.OrderActivities
}

func NewOrderWorkflow(
	activities activity.OrderActivities,
) OrderWorkflow {
	wf := OrderWorkflow{
		OrderActivities: activities,
	}
	return wf
}

func (w OrderWorkflow) Demo(ctx workflow.Context, input models.IncrMemberOrderRequest) (string, error) {
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        500, // 0 is unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retryPolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	// Withdraw money.
	var incrMemberOutput giaekv1.Order

	err := workflow.ExecuteActivity(ctx, w.OrderActivities.IncrMemberOrder, input).Get(ctx, &incrMemberOutput)

	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Demo complete: %s", incrMemberOutput.OrderId)
	return result, nil
}
