package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jxo-me/plus-core/pkg/v2/puller/model"
	"github.com/jxo-me/rabbitmq-go"
	"log"
)

const (
	ExchangeType     = "topic"
	ExchangeName     = "puller.compensate.topic"
	DefaultQueueName = "puller.compensate.allVendor"
	RoutingKeyPrefix = "puller.compensate"
)

type MQConsumer struct {
	conn   *rabbitmq.Conn
	client *rabbitmq.Consumer
	queue  string
}

func NewMQConsumer(ctx context.Context, url, queue string, vendorKeys []string) (*MQConsumer, error) {
	conn, err := rabbitmq.NewConn(
		ctx,
		url,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}
	if queue == "" {
		queue = DefaultQueueName
	}
	var routingKeys []string
	for _, vendor := range vendorKeys {
		routingKey := fmt.Sprintf("%s.%s", RoutingKeyPrefix, vendor)
		routingKeys = append(routingKeys, routingKey)
	}
	consumer, err := rabbitmq.NewConsumer(
		ctx,
		conn,
		func(ctx context.Context, rw *rabbitmq.ResponseWriter, d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		queue,
		rabbitmq.WithConsumerOptionsRoutingKeys(routingKeys),
		rabbitmq.WithConsumerOptionsExchangeName(ExchangeName),
		rabbitmq.WithConsumerOptionsExchangeKind(ExchangeType),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsQueueDurable,
	)
	return &MQConsumer{conn: conn, client: consumer, queue: queue}, nil
}

func (c *MQConsumer) StartConsuming(ctx context.Context, handler func(ctx context.Context, msg model.CompensateMessage) error) error {
	//defer c.client.Close(ctx)
	// block main thread - wait for shutdown signal
	err := c.client.Run(ctx, func(ctx context.Context, rw *rabbitmq.ResponseWriter, d rabbitmq.Delivery) rabbitmq.Action {
		var msg model.CompensateMessage
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			return rabbitmq.NackRequeue
		}
		if err := handler(ctx, msg); err != nil {
			return rabbitmq.NackRequeue
		}
		return rabbitmq.Ack
	})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (c *MQConsumer) Close(ctx context.Context) error {
	if c.client != nil {
		c.client.Close(ctx)
	}
	if c.conn != nil {
		err := c.conn.Close(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
