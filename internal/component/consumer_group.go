package component

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx"
	"strings"
	shared "tutorial.io/tutorial-order/internal"
	"tutorial.io/tutorial-order/internal/consumer"
)

func NewConsumerGroup(
	lc fx.Lifecycle,
	config *shared.Config,
	consumer *consumer.Consumer,
) sarama.ConsumerGroup {
	topics := []string{
		shared.KafkaTopicSomething,
	}
	version, err := sarama.ParseKafkaVersion(config.KafkaVersion)
	if err != nil {
		log.WithError(err).Fatalf("error parsing Kafka version: %s", config.KafkaVersion)
	}

	// Setup Kafka consumer
	log.Infof("Consumer connecting to Kafka broker at %s", config.KafkaBroker)
	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = version
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(config.KafkaBroker, ","), shared.KafkaGroupIdOrder, consumerConfig)
	if err != nil {
		log.WithError(err).Fatalf("Error creating consumer group client: %v", err)
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				propagators := propagation.TraceContext{}
				consumer := otelsarama.WrapConsumerGroupHandler(consumer, otelsarama.WithPropagators(propagators))
				for {
					err := consumerGroup.Consume(context.Background(), topics, consumer)
					if err != nil {
						if errors.Is(err, sarama.ErrClosedConsumerGroup) {
							return
						}
						log.Panicf("Error from consumer: %v", err)
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return consumerGroup.Close()
		},
	})
	return consumerGroup
}
