package component

import (
	"context"
	log "github.com/sirupsen/logrus"
	tc "go.temporal.io/sdk/client"
	totel "go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
	shared "tutorial.io/tutorial-order/internal"
	"tutorial.io/tutorial-order/internal/activity"
	"tutorial.io/tutorial-order/internal/workflow"
)

func NewTemporalClient(
	lc fx.Lifecycle,
	config *shared.Config,
) tc.Client {
	// Setup temporal client
	tracingInterceptor, err := totel.NewTracingInterceptor(totel.TracerOptions{})
	if err != nil {
		log.WithError(err).Fatal("failed to create temporal tracing interceptor")
	}
	log.Infof("Dailing Temporal at %s", config.TemporalEndpoint)
	client, err := tc.Dial(tc.Options{
		HostPort:     config.TemporalEndpoint,
		Namespace:    config.TemporalNamespace,
		Interceptors: []interceptor.ClientInterceptor{tracingInterceptor},
		Logger:       logur.LoggerToKV(logrusadapter.New(log.StandardLogger())),
	})
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})

	return client
}

func NewTemporalWorker(
	lc fx.Lifecycle,
	client tc.Client,
	orderWorkflow workflow.OrderWorkflow,
	orderActivities activity.OrderActivities,
) worker.Worker {
	w := worker.New(client, shared.TemporalTaskQueueOrder, worker.Options{})
	w.RegisterWorkflow(orderWorkflow.Demo)
	w.RegisterActivity(orderActivities.IncrMemberOrder)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := w.Run(worker.InterruptCh())
				if err != nil {
					log.WithError(err).Fatalln("unable to start worker")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			w.Stop()
			return nil
		},
	})
	return w
}
