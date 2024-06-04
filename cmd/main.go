package main

import (
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	fxlogrus "github.com/takt-corp/fx-logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	tc "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"os"
	shared "tutorial.io/tutorial-order/internal"
	"tutorial.io/tutorial-order/internal/activity"
	"tutorial.io/tutorial-order/internal/component"
	"tutorial.io/tutorial-order/internal/consumer"
	"tutorial.io/tutorial-order/internal/server"
	"tutorial.io/tutorial-order/internal/workflow"
)

var (
	config shared.Config
)

func main() {
	app := &cli.App{
		Name:  "order",
		Usage: "Tutorial order service server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "trace-service-name",
				Usage:       "The name of the service for tracing",
				EnvVars:     []string{"TRACE_SERVICE_NAME"},
				Destination: &config.TraceServiceName,
			},
			&cli.IntFlag{
				Name:        "grpc-port",
				Usage:       "gRPC server port",
				Value:       8000,
				EnvVars:     []string{"GRPC_PORT"},
				Destination: &config.GrpcPort,
			},
			&cli.StringFlag{
				Name:        "postgres-host",
				Usage:       "PostgresSQL DB host address",
				EnvVars:     []string{"POSTGRES_HOST"},
				Destination: &config.PostgresHost,
			},
			&cli.IntFlag{
				Name:        "postgres-port",
				Usage:       "PostgresSQL DB port number",
				Value:       5432,
				EnvVars:     []string{"POSTGRES_PORT"},
				Destination: &config.PostgresPort,
			},
			&cli.StringFlag{
				Name:        "postgres-user",
				Usage:       "PostgresSQL DB user",
				EnvVars:     []string{"POSTGRES_USER"},
				Destination: &config.PostgresUser,
			},
			&cli.StringFlag{
				Name:        "postgres-password",
				Usage:       "PostgresSQL DB password",
				EnvVars:     []string{"POSTGRES_PASSWORD"},
				Destination: &config.PostgresPassword,
			},
			&cli.StringFlag{
				Name:        "postgres-db",
				Usage:       "PostgresSQL DB name",
				EnvVars:     []string{"POSTGRES_DB"},
				Destination: &config.PostgresDb,
			},
			&cli.StringFlag{
				Name:        "postgres-schema",
				Usage:       "PostgresSQL DB schema",
				EnvVars:     []string{"POSTGRES_SCHEMA"},
				Destination: &config.PostgresSchema,
			},
			&cli.IntFlag{
				Name:        "db-max-idle-conns",
				Usage:       "PostgresSQL DB max idle connections",
				Value:       2,
				EnvVars:     []string{"DB_MAX_IDLE_CONNS"},
				Destination: &config.DbMaxIdleConns,
			},
			&cli.IntFlag{
				Name:        "db-max-open-conns",
				Usage:       "PostgresSQL DB max open connections",
				Value:       20,
				EnvVars:     []string{"DB_MAX_OPEN_CONNS"},
				Destination: &config.DbMaxOpenConns,
			},
			&cli.StringFlag{
				Name:        "kafka-version",
				Usage:       "Kafka version",
				Value:       sarama.DefaultVersion.String(),
				EnvVars:     []string{"KAFKA_VERSION"},
				Destination: &config.KafkaVersion,
			},
			&cli.StringFlag{
				Name:        "kafka-broker",
				Usage:       "Kafka broker address",
				EnvVars:     []string{"KAFKA_BROKER"},
				Destination: &config.KafkaBroker,
			},
			&cli.StringFlag{
				Name:        "redis-host",
				Usage:       "Redis host address",
				EnvVars:     []string{"REDIS_HOST"},
				Destination: &config.RedisHost,
			},
			&cli.IntFlag{
				Name:        "redis-port",
				Usage:       "Redis port number",
				Value:       6379,
				EnvVars:     []string{"REDIS_PORT"},
				Destination: &config.RedisPort,
			},
			&cli.StringFlag{
				Name:        "redis-password",
				Usage:       "Redis password",
				EnvVars:     []string{"REDIS_PASSWORD"},
				Destination: &config.RedisPassword,
			},
			&cli.IntFlag{
				Name:        "redis-db",
				Usage:       "Redis database number",
				Value:       0,
				EnvVars:     []string{"REDIS_DB"},
				Destination: &config.RedisDb,
			},
			&cli.StringFlag{
				Name:        "temporal-endpoint",
				Usage:       "Temporal endpoint",
				EnvVars:     []string{"TEMPORAL_ENDPOINT"},
				Destination: &config.TemporalEndpoint,
			},
			&cli.StringFlag{
				Name:        "temporal-namespace",
				Usage:       "Temporal namespace",
				EnvVars:     []string{"TEMPORAL_NAMESPACE"},
				Destination: &config.TemporalNamespace,
			},
		},
		Action: execute,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func execute(cCtx *cli.Context) error {
	log.Infof("Starting %s", config.TraceServiceName)

	// Setup and instrument logrus.
	log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
	)))
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return &fxlogrus.LogrusLogger{Logger: log.StandardLogger()}
		}),
		fx.Supply(&config),
		fx.Provide(
			component.NewOtelClient,
			component.NewSnowflake,
			component.NewDb,
			component.NewTemporalClient,
			component.NewTemporalWorker,
			component.NewRedisClient,
			component.NewProducer,
			component.NewConsumerGroup,
			fx.Annotate(
				component.NewGrpcServer,
				fx.ParamTags("", "", `group:"grpcServices"`),
			),
			component.NewRestyClient,
			consumer.NewConsumer,
			workflow.NewOrderWorkflow,
			activity.NewOrderActivities,
			AsGrpcService(server.NewOrderServer),
		),
		fx.Invoke(
			func(*otlptrace.Client) {},
			func(*gorm.DB) {},
			func(tc.Client) {},
			func(worker.Worker) {},
			func(*redis.Client) {},
			func(sarama.ConsumerGroup) {},
			func(sarama.SyncProducer) {},
			func(*grpc.Server) {},
		),
	).Run()
	return nil
}

func AsGrpcService(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(server.GrpcService)),
		fx.ResultTags(`group:"grpcServices"`),
	)
}
