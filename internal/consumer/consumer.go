package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/propagation"
	shared "tutorial.io/tutorial-order/internal"
)

type Consumer struct {
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (consumer Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumer Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (consumer Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		propagators := propagation.TraceContext{}
		ctx := propagators.Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(msg))
		log.WithContext(ctx).
			WithField("topic", msg.Topic).
			WithField("key", string(msg.Key)).
			WithField("value", string(msg.Value)).
			Info("Consumed message")
		switch msg.Topic {
		case shared.KafkaTopicSomething:
			//var value models.FlushReplyBufferRequest
			//err := json.Unmarshal(msg.Value, &value)
			//if err != nil {
			//	log.WithContext(ctx).
			//		WithError(err).
			//		Error("Error parsing JSON")
			//}
			//consumer.srv.HandleIt(ctx, &value)
		}
		// after processing the message, mark the offset
		session.MarkMessage(msg, "")
	}
	return nil
}
