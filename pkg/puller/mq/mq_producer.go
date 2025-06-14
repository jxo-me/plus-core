package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jxo-me/plus-core/pkg/v2/puller/model"
	"time"

	"github.com/jxo-me/rabbitmq-go"
)

type MQProducer struct {
	conn     *rabbitmq.Conn
	client   *rabbitmq.Publisher
	exchange string
}

func Build() {

}

func NewMQProducer(ctx context.Context, url string) (*MQProducer, error) {
	conn, err := rabbitmq.NewConn(
		ctx,
		url,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}
	publisher, err := rabbitmq.NewPublisher(
		ctx,
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(ExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}
	return &MQProducer{conn: conn, client: publisher, exchange: ExchangeName}, nil
}

func (p *MQProducer) PushCompensateTask(ctx context.Context, vendor string, start, end time.Time) error {
	msg := model.CompensateMessage{Vendor: vendor, Start: start, End: end}
	payload, _ := json.Marshal(msg)
	routingKey := fmt.Sprintf("%s.%s", RoutingKeyPrefix, vendor)
	return p.client.Publish(payload, []string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(p.exchange),
	)
	//return p.client.Publish(ctx, p.exchange, routingKey, payload)
}

func (p *MQProducer) Close(ctx context.Context) error {
	if p.client != nil {
		p.client.Close(ctx)
	}
	if p.conn != nil {
		err := p.conn.Close(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
