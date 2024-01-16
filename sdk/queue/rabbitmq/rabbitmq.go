package rabbitmq

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	messageLib "github.com/jxo-me/plus-core/core/v2/message"
	queueLib "github.com/jxo-me/plus-core/core/v2/queue"
	"github.com/jxo-me/plus-core/sdk/v2/message"
	"github.com/jxo-me/plus-core/sdk/v2/queue"
	"github.com/jxo-me/rabbitmq-go"
	"sync"
)

func NewRabbitMQ(
	ctx context.Context,
	dsn string,
	maxConnections int,
	reconnectInterval int,
	cfg *rabbitmq.Config,
	logger rabbitmq.Logger,
) (*RabbitMQ, error) {
	p, err := NewConnectionPool(ctx, dsn, maxConnections, reconnectInterval, logger, cfg)
	if err != nil {
		return nil, err
	}
	r := &RabbitMQ{
		producers:  map[string]*rabbitmq.Publisher{},
		consumers:  make([]*rabbitmq.Consumer, 0),
		rpcClients: map[string]*rabbitmq.RpcClient{},
		Logger:     logger,
		pool:       p,
	}

	return r, nil
}

// RabbitMQ cache implement
type RabbitMQ struct {
	Handler          []rabbitmq.Handler
	mux              sync.RWMutex
	consumers        []*rabbitmq.Consumer
	ConsumerOptions  *rabbitmq.ConsumerOptions
	producers        map[string]*rabbitmq.Publisher
	rpcClients       map[string]*rabbitmq.RpcClient
	PublisherOptions *rabbitmq.PublisherOptions
	Logger           rabbitmq.Logger
	pool             *ConnectionPool
}

func (r *RabbitMQ) String() string {
	return "rabbitmq"
}

func (r *RabbitMQ) newConsumer(ctx context.Context, queueName string, handler rabbitmq.Handler, options queueLib.ConsumeOptions) (*rabbitmq.Consumer, error) {
	var err error
	var conn *rabbitmq.Conn
	conn, err = r.pool.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		// 将连接释放回连接池
		r.pool.ReleaseConnection(conn)
	}()
	var optionFuncs []func(*rabbitmq.ConsumerOptions)
	if options.BindingExchange.Declare {
		optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsExchangeDeclare)
	}
	if options.BindingExchange.Durable {
		optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsExchangeDurable)
	}
	if options.BindingExchange.Passive {
		optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsExchangePassive)
	}
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsQueueDurable)
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsExchangeName(options.BindingExchange.Name))
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsExchangeKind(options.BindingExchange.Kind))
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsLogger(r.Logger))
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsRoutingKeys(options.BindingRoutingKeys))
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsConsumerName(options.ConsumerName))
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsConcurrency(options.Concurrency))
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsConsumerAutoAck(options.ConsumerAutoAck))
	optionFuncs = append(optionFuncs, rabbitmq.WithConsumerOptionsQOSPrefetch(options.QOSPrefetch))
	return rabbitmq.NewConsumer(ctx,
		conn,
		handler,
		queueName,
		optionFuncs...,
	)
}

func (r *RabbitMQ) newProducer(ctx context.Context) (*rabbitmq.Publisher, error) {
	var err error
	var conn *rabbitmq.Conn
	conn, err = r.pool.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		// 将连接释放回连接池
		r.pool.ReleaseConnection(conn)
	}()
	return rabbitmq.NewPublisher(ctx,
		conn,
		rabbitmq.WithPublisherOptionsLogger(r.Logger),
	)
}

// Publish 消息入生产者
func (r *RabbitMQ) Publish(ctx context.Context, message messageLib.IMessage, optionFuncs ...func(*queueLib.PublishOptions)) error {
	// exchange exchangeType routingKey
	rb, err := queue.Marshal(message.GetValue())
	if err != nil {
		return err
	}
	options := &queueLib.PublishOptions{
		ContentType: "application/json",
		MessageID:   gctx.CtxId(ctx),
	}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	var p *rabbitmq.Publisher
	var ok bool
	if p, ok = r.producers[options.Exchange]; !ok {
		p, err = r.newProducer(ctx)
		if err != nil {
			glog.Warning(ctx, "rabbitmq newProducer error:", err)
			return err
		}
		r.producers[options.Exchange] = p
	}

	err = p.PublishWithContext(
		ctx,
		rb,
		[]string{message.GetRoutingKey()},
		rabbitmq.WithPublishOptionsExchange(options.Exchange),
		rabbitmq.WithPublishOptionsContentType(options.ContentType),
		rabbitmq.WithPublishOptionsMessageID(options.MessageID),
		rabbitmq.WithPublishOptionsAppID(options.AppID),
		rabbitmq.WithPublishOptionsUserID(options.UserID),
		rabbitmq.WithPublishOptionsReplyTo(options.ReplyTo),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
	return err
}

func (r *RabbitMQ) newRpcClient(ctx context.Context, optionFuncs ...func(*rabbitmq.ClientOptions)) (*rabbitmq.RpcClient, error) {
	var err error
	var conn *rabbitmq.Conn
	conn, err = r.pool.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		// 将连接释放回连接池
		r.pool.ReleaseConnection(conn)
	}()
	return rabbitmq.NewRpcClient(ctx,
		conn,
		rabbitmq.WithClientOptionsLogger(r.Logger),
	)
}

// RpcRequest AMQP RPC Request
func (r *RabbitMQ) RpcRequest(ctx context.Context, key string, data []byte, optionFuncs ...func(*queueLib.ClientOptions)) ([]byte, error) {
	options := queueLib.GetDefaultClientOptions(key)
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}
	var (
		p        *rabbitmq.RpcClient
		ok       bool
		err      error
		optFuncs []func(*rabbitmq.ClientOptions)
		resp     []byte
	)
	if options.PublishOptions.Mandatory {
		optFuncs = append(optFuncs, rabbitmq.WithClientPublishOptionsMandatory)
	}
	if options.PublishOptions.Immediate {
		optFuncs = append(optFuncs, rabbitmq.WithClientPublishOptionsImmediate)
	}
	if options.ConsumeOptions.ConsumerName != "" {
		optFuncs = append(optFuncs, rabbitmq.WithClientOptionsConsumerName(options.ConsumeOptions.ConsumerName))
	}
	if options.ConsumeOptions.ConsumerAutoAck {
		optFuncs = append(optFuncs, rabbitmq.WithClientOptionsConsumerAutoAck(true))
	}
	if options.ConsumeOptions.Exclusive {
		optFuncs = append(optFuncs, rabbitmq.WithClientOptionsConsumerExclusive)
	}
	if options.QueueOptions.Durable {
		optFuncs = append(optFuncs, rabbitmq.WithClientOptionsQueueDurable)
	}
	if options.QueueOptions.AutoDelete {
		optFuncs = append(optFuncs, rabbitmq.WithClientOptionsQueueAutoDelete)
	}
	if options.QueueOptions.Exclusive {
		optFuncs = append(optFuncs, rabbitmq.WithClientOptionsQueueExclusive)
	}
	if options.QueueOptions.Args != nil {
		optFuncs = append(optFuncs, rabbitmq.WithClientOptionsQueueArgs(options.QueueOptions.Args))
	}
	if p, ok = r.rpcClients[options.GroupName]; !ok {
		p, err = r.newRpcClient(ctx, optFuncs...)
		if err != nil {
			glog.Warning(ctx, "rabbitmq new rpcClient error:", err)
			return nil, err
		}
		r.rpcClients[options.GroupName] = p
	}

	resp, err = p.PublishWithContext(
		ctx,
		data,
		key,
	)
	return resp, err
}

// Consumer 监听消费者
func (r *RabbitMQ) Consumer(ctx context.Context, queueName string, consumerFunc queueLib.ConsumerFunc, optionFuncs ...func(*queueLib.ConsumeOptions)) {
	options := queueLib.GetDefaultConsumeOptions()
	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}
	var c *rabbitmq.Consumer
	var err error
	//var ok bool
	r.mux.Lock()
	defer r.mux.Unlock()
	//if c, ok = r.consumers[queueName]; !ok {
	handler := func(ctx context.Context, rw *rabbitmq.ResponseWriter, d rabbitmq.Delivery) rabbitmq.Action {
		//glog.Debug(ctx, fmt.Sprintf("rabbitmq handle raw message: %s", gconv.String(d)))
		m := new(message.Message)
		m.SetValue(d.Body)
		m.SetRoutingKey(d.RoutingKey)
		m.SetExchange(d.Exchange)
		m.SetId(d.MessageId)
		if d.Redelivered {
			m.SetErrorCount(d.DeliveryTag)
		}
		err = consumerFunc(ctx, rw, m)
		if err != nil {
			glog.Warning(ctx, fmt.Sprintf("RabbitMQ Requeue error:%s msg: %v", err.Error(), gconv.String(m)))
			return rabbitmq.NackRequeue
		}
		// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		return rabbitmq.Ack
	}
	c, err = r.newConsumer(ctx, queueName, handler, options)
	if err != nil {
		glog.Error(ctx, "rabbitmq newConsumer error:", err)
		return
	}
	r.consumers = append(r.consumers, c)
	//r.consumers[queueName] = c
	//}

}

func (r *RabbitMQ) Run(ctx context.Context) {
	return
}

func (r *RabbitMQ) Shutdown(ctx context.Context) {
	for _, pd := range r.producers {
		pd.Close(ctx)
	}
	for _, pushConsumer := range r.consumers {
		pushConsumer.Close(ctx)
	}
	if r.pool != nil {
		err := r.pool.Close(ctx)
		if err != nil {
			glog.Warning(ctx, "rabbitmq conn close error:", err.Error())
		}
	}
}
