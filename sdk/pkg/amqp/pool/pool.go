package pool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"sync"
	"sync/atomic"

	"github.com/streadway/amqp"

	"github.com/jxo-me/plus-core/sdk/pkg/amqp/config"
)

var (
	// ErrTooManyConn means opened connection num is overflow the config.maxConnections limit
	ErrTooManyConn = errors.New("pool: too many connections")

	// ErrBadConn occured when openChannel called on a closed channel or bad Conn
	ErrBadConn = errors.New("pool: connection is bad")

	// ErrChannelNotAllClosed occured when the connection's channels not all closed
	// when close a connection, channels must be all closed
	ErrChannelNotAllClosed = errors.New("pool: connection's channels not all closed")

	// ErrNacked occured when send message not acked
	ErrNacked = errors.New("message not acked")

	// ErrPoolClosed occured when the pool was closed
	ErrPoolClosed = errors.New("pool closed")
)

const RESEND_MAX_NUM = 5

// Connection represent a amqp real connection, which record the connection to user
type Connection struct {
	Conn             *amqp.Connection
	l                sync.RWMutex
	numOpenedChannel int
	done             chan error
}

// Channel represent a amqp channel, which expose connection to user
type Channel struct {
	conn      *Connection
	cha       *amqp.Channel
	confirmCh chan amqp.Confirmation
}

// ConnPool is the real connection pool
type ConnPool struct {
	ctx   context.Context
	conf  *config.Amqp
	conns []*Connection

	// defer close the unused connection to reduce pool lock time
	connDelayCloseCh chan *Connection
	connDelayClosed  chan struct{}

	l sync.Mutex

	reqChaList *ReqChaList

	// idle channels
	idleChas []*Channel

	chaBusyNum int32

	closed bool
}

// ConnPoolStats contains current pool states returned by call pool.Stats()
type ConnPoolStats struct {
	IdleChaNum int
	ConnNum    int
	BusyChaNum int32
	ReqChaNum  int
}

func (conn *Connection) getNumOpenedChannel() int {
	conn.l.RLock()
	defer conn.l.RUnlock()
	return conn.numOpenedChannel
}

func (conn *Connection) openChannel(confirm bool) (*Channel, error) {
	conn.l.Lock()
	defer conn.l.Unlock()

	if conn.Conn == nil {
		return nil, ErrBadConn
	}

	amqpCha, err := conn.Conn.Channel()
	if err != nil {
		return nil, err
	}

	cha := &Channel{
		conn: conn,
		cha:  amqpCha,
	}

	if confirm == true {
		// always in confirm mode
		if err := cha.cha.Confirm(false); err != nil {
			return nil, WrapError(err, "channel set to confirm mode failed")
		}
		cha.confirmCh = cha.cha.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	conn.numOpenedChannel++
	return cha, nil
}

// getChannel get a new channel from conn
func (conn *Connection) getChannel(ctx context.Context, confirm bool) (*Channel, error) {
	// open new channel
	cha, err := conn.openChannel(confirm)
	if err == amqp.ErrClosed || err == ErrBadConn {
		g.Log().Printf(ctx, "ConnPool.getChannel: %s", ErrBadConn)
		return nil, err
	}

	return cha, nil
}

func (conn *Connection) decrNumOpenedChannel() {
	conn.l.Lock()
	defer conn.l.Unlock()
	conn.numOpenedChannel--
}

// close the amqp connection and put it to cache pool
func (conn *Connection) close(ctx context.Context, force bool) error {
	conn.l.Lock()
	defer conn.l.Unlock()

	if !force && conn.numOpenedChannel > 0 {
		return ErrChannelNotAllClosed
	}

	if conn.Conn != nil {
		err := conn.Conn.Close()
		if err != nil {
			g.Log().Printf(ctx, "Connection.close: %s\n", err)
		}
	}
	conn.Conn = nil
	return nil
}

func (conn *Connection) Shutdown(ctx context.Context, cop *ConnPool) error {
	if err := conn.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	cop.l.Lock()
	defer cop.l.Unlock()
	err := cop.removeConn(conn)
	if err != nil {
		g.Log().Errorf(ctx, "AMQP Connection Pool remove error: %s", err)
	}

	defer g.Log().Printf(ctx, "AMQP shutdown OK")

	// wait for handle() to exit
	return <-conn.done
}

// close the amqp channel and decr it's connection numOpenedChannel
func (cha *Channel) close(ctx context.Context) {
	// not care about channel close error, because it's the client action
	err := cha.cha.Close()
	if err != nil {
		g.Log().Printf(ctx, "amqp channel close error: %s\n", err.Error())
	}

	cha.conn.decrNumOpenedChannel()
	cha.cha = nil
	cha.conn = nil
	g.Log().Printf(ctx, "[channel] channel close")
}

// NewPool return a pool inited with the given config
func NewPool(ctx context.Context, conf *config.Amqp) *ConnPool {
	pool := &ConnPool{
		ctx:   ctx,
		conf:  conf,
		conns: make([]*Connection, 0, conf.MaxConnections),

		connDelayCloseCh: make(chan *Connection, 10000),
		connDelayClosed:  make(chan struct{}),

		idleChas: make([]*Channel, 0, conf.MaxIdleChannels),

		reqChaList: &ReqChaList{},
	}

	go func() {
		for conn := range pool.connDelayCloseCh {
			err := conn.close(ctx, true)
			if err != nil {
				g.Log().Printf(ctx, "Conn.close: %s\n", err.Error())
			}
		}
		close(pool.connDelayClosed)
	}()

	return pool
}

func (cop *ConnPool) incrChaBusyNum() {
	atomic.AddInt32(&cop.chaBusyNum, int32(1))
}

func (cop *ConnPool) decrChaBusyNum() {
	atomic.AddInt32(&cop.chaBusyNum, int32(-1))
}

func (cop *ConnPool) getChaBusyNum() int32 {
	return atomic.LoadInt32(&cop.chaBusyNum)
}

// GetConf get a copy of current config
func (cop *ConnPool) GetConf() config.Amqp {
	conf := *cop.conf
	return conf
}

// CloseAll close the whole connections and channels
func (cop *ConnPool) CloseAll() {
	cop.l.Lock()
	defer cop.l.Unlock()
	cop.close()
}

// close will close the pool and it's connections
func (cop *ConnPool) close() {
	if cop.closed {
		return
	}
	cop.closed = true

	// wait for connection all closed
	close(cop.connDelayCloseCh)
	<-cop.connDelayClosed

	for i, cha := range cop.idleChas {
		cha.close(cop.ctx)
		cop.idleChas[i] = nil
	}
	cop.idleChas = nil

	for i, conn := range cop.conns {
		err := conn.close(cop.ctx, true)
		if err != nil {
			g.Log().Printf(cop.ctx, "amqp Conn colse error: %s\n", err.Error())
		}
		cop.conns[i] = nil
	}
	cop.conns = nil
}

func (cop *ConnPool) removeConn(conn *Connection) error {

	if cop.conf.Debug {
		g.Log().Printf(cop.ctx, "[Conn] old Conn closed")
	}

	foundIdx := -1
	for i := 0; i < len(cop.conns); i++ {
		if conn == cop.conns[i] {
			foundIdx = i
			break
		}
	}
	// delete from pool
	if foundIdx > -1 {
		copy(cop.conns[foundIdx:], cop.conns[foundIdx+1:])
		cop.conns[len(cop.conns)-1] = nil
		cop.conns = cop.conns[:len(cop.conns)-1]
	}

	// if pool closed , direct close the connection
	if cop.closed {
		if err := conn.close(cop.ctx, true); err != nil {
			return err
		}
		return nil
	}

	// put into Conn delay close channel
	cop.connDelayCloseCh <- conn
	return nil
}

func (cop *ConnPool) getConn(isNew bool) (*Connection, error) {
	if isNew == false && len(cop.conns) > 0 {
		for i := 0; i < len(cop.conns); i++ {
			// notice: the first connection may handel more channels
			if cop.conns[i].getNumOpenedChannel() < cop.conf.MaxChannelsPerConnection {
				return cop.conns[i], nil
			}
		}
	}

	if cop.conf.MaxConnections > 0 && len(cop.conns) >= cop.conf.MaxConnections {
		return nil, ErrTooManyConn
	}

	amqpConn, err := amqp.Dial(cop.conf.DSN)
	if err != nil {
		return nil, WrapError(err, "amqp.Dial")
	}
	conn := &Connection{Conn: amqpConn, numOpenedChannel: 0, done: make(chan error)}
	go func() {
		fmt.Printf("[Conn] closing: %s", <-conn.Conn.NotifyClose(make(chan *amqp.Error)))
	}()
	cop.conns = append(cop.conns, conn)

	if cop.conf.Debug {
		g.Log().Printf(cop.ctx, "[Conn] new Conn opened")
	}

	return conn, nil
}

// Stats return current pool states
func (cop *ConnPool) Stats() *ConnPoolStats {
	cop.l.Lock()
	defer cop.l.Unlock()

	return &ConnPoolStats{
		IdleChaNum: len(cop.idleChas),
		ConnNum:    len(cop.conns),
		BusyChaNum: cop.getChaBusyNum(),
		ReqChaNum:  cop.reqChaList.Len(),
	}
}

func (cop *ConnPool) putChannel(cha *Channel) {

	cop.decrChaBusyNum()

	// if channel request is notified, skip put into idleChas
	if cop.reqChaList.NotifyOne(cha) {
		return
	}

	cop.l.Lock()
	lenFree := len(cop.idleChas)
	cop.l.Unlock()

	if lenFree >= cop.conf.MaxIdleChannels {
		cop.probeCloseChannel(cha)
		return
	}

	cop.l.Lock()
	cop.idleChas = append(cop.idleChas, cha)
	cop.l.Unlock()

}

// getChannel get a free channel from pool
func (cop *ConnPool) getChannel(confirm bool) (*Channel, error) {

GETFREECHANNEL:
	cop.l.Lock()

	if cop.closed {
		cop.l.Unlock()
		return nil, ErrPoolClosed
	}

	// step1: reuse free channels
	if len(cop.idleChas) > 0 {
		var cha *Channel
		// shift from free pool
		cha, cop.idleChas = cop.idleChas[0], cop.idleChas[1:]
		cop.incrChaBusyNum()

		cop.l.Unlock()
		return cha, nil
	}

	// step2: get connection
	conn, err := cop.getConn(false)
	if err == ErrTooManyConn {
		// unlock
		cop.l.Unlock()

		// wait for available channel
		// if wait return and has a free channel, use it
		ch := make(chan *Channel)
		cop.reqChaList.Put(ch)
		if cha := <-ch; cha != nil {
			cop.incrChaBusyNum()
			return cha, nil
		}

		// retry
		goto GETFREECHANNEL
	} else if err != nil {
		cop.close()
		// unlock
		cop.l.Unlock()
		FailOnError(cop.ctx, err, "ConnPool.getChannel")
	}

	// step3: open new channel
	cha, err := conn.openChannel(confirm)

	if err == amqp.ErrClosed || err == ErrBadConn {
		g.Log().Printf(cop.ctx, "ConnPool.getChannel: %s", ErrBadConn)
		err := cop.removeConn(conn)
		if err != nil {
			g.Log().Printf(cop.ctx, "ConnPool.getChannel: %s", err.Error())
		}
		// unlock
		cop.l.Unlock()
		goto GETFREECHANNEL
	} else if err == amqp.ErrChannelMax {
		g.Log().Printf(cop.ctx, "ConnPool.getChannel: %s", amqp.ErrChannelMax)
		cop.l.Unlock()
		goto GETFREECHANNEL
	} else if err != nil {
		cop.close()
		// unlock
		cop.l.Unlock()
		FailOnError(cop.ctx, err, "ConnPool.getChannel")
	}

	cop.incrChaBusyNum()

	if cop.conf.Debug {
		g.Log().Printf(cop.ctx, "[channel] new channel opened")
	}

	// unlock
	cop.l.Unlock()
	return cha, nil
}

func (cop *ConnPool) probeCloseChannel(cha *Channel) {
	conn := cha.conn
	cha.close(cop.ctx)

	if cop.conf.Debug {
		g.Log().Printf(cop.ctx, "[channel] old channel closed")
	}

	cop.l.Lock()
	defer cop.l.Unlock()
	if conn.getNumOpenedChannel() == 0 && len(cop.conns) > cop.conf.MinConnections {
		err := cop.removeConn(conn)
		if err != nil {
			g.Log().Printf(cop.ctx, "remove Conn error:%s\n", err)
		}
	}
}

// ConfirmSendMsg send message with confirm mode
func (cop *ConnPool) ConfirmSendMsg(exchange string, routingKey string, data []byte) error {
	if cop.conf.Debug {
		defer func() {
			stats, _ := json.Marshal(cop.Stats())
			conf, _ := json.Marshal(cop.conf)
			g.Log().Printf(cop.ctx, "debugStats: %s %s\n", stats, conf)
		}()
	}

	var err error
	var cha *Channel

	for i := 0; i < RESEND_MAX_NUM; i++ {
		cha, err = cop.getChannel(true)
		if err != nil {
			return err
		}
		err = cha.cha.Publish(
			exchange,   // exchange
			routingKey, // routing key
			true,       // mandatory，若为true，则当没有对应的队列，不ack
			false,      // immediate，若为true，则当没有消费者消费，不ack
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         data,
				DeliveryMode: 2, // 持久消息
			})
		if err == nil {
			break
		}
		cop.probeCloseChannel(cha)
	}

	// waiting for the server confirm
	confirmed := <-cha.confirmCh

	// put current channel into idle pool
	cop.putChannel(cha)
	//cop.idleChaPutCh <- cha
	if confirmed.Ack {
		return nil
	}
	// @todo if channel closed by putChannel method, message would nacked
	g.Log().Printf(cop.ctx, "%s: %#v, %s, %x\n", ErrNacked, confirmed, data, data)
	return ErrNacked
}

func (cop *ConnPool) Declare(exchange, exchangeType, queueName, bindingKey string, queueArgs amqp.Table) error {
	var err error

	var channel *Channel
	channel, err = cop.getChannel(false)
	if err != nil {
		return err
	}
	//args := amqp.Table{"x-queue-type": "quorum"}
	err = channel.cha.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return WrapError(err, "Exchange Declare error !")
	}
	queue, err := channel.cha.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // noWait
		queueArgs, // arguments
	)
	if err != nil {
		return WrapError(err, "Queue Declare error !")
	}

	err = channel.cha.QueueBind(
		queue.Name, // name of the queue
		bindingKey, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return WrapError(err, "Queue Bind error !")
	}
	channel.close(cop.ctx)
	cop.l.Lock()
	cop.decrChaBusyNum()
	cop.l.Unlock()
	return nil
}

func (cop *ConnPool) DeclareRandom(exchange, exchangeType, queueName, bindingKey string, queueArgs amqp.Table) error {
	var err error

	var channel *Channel
	channel, err = cop.getChannel(false)
	if err != nil {
		return err
	}

	err = channel.cha.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		true,         // noWait
		nil,          // arguments
	)

	if err != nil {
		return WrapError(err, "Exchange Declare error !")
	}

	queue, err := channel.cha.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		queueArgs, // arguments
	)
	if err != nil {
		return WrapError(err, "Queue Declare error !")
	}

	err = channel.cha.QueueBind(
		queue.Name,
		bindingKey,
		exchange,
		false,
		queueArgs,
	)
	if err != nil {
		return WrapError(err, "Queue Bind error !")
	}

	channel.close(cop.ctx)
	cop.l.Lock()
	cop.decrChaBusyNum()
	cop.l.Unlock()
	return nil
}

func (cop *ConnPool) HashDeclare(exchange, exchangeType, dstExchange, bindingKey string) error {
	var err error

	var channel *Channel
	channel, err = cop.getChannel(false)
	if err != nil {
		return err
	}

	err = channel.cha.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return WrapError(err, "Exchange Declare error !")
	}

	err = channel.cha.ExchangeBind(
		dstExchange, // name of the queue
		bindingKey,  // bindingKey
		exchange,    // sourceExchange
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return WrapError(err, "Exchange Bind error !")
	}
	channel.close(cop.ctx)
	cop.l.Lock()
	cop.decrChaBusyNum()
	cop.l.Unlock()
	return nil
}

func (cop *ConnPool) Publish(exchange, exchangeType, routingKey, body string) error {
	if cop.conf.Debug {
		defer func() {
			stats, _ := json.Marshal(cop.Stats())
			conf, _ := json.Marshal(cop.conf)
			g.Log().Printf(cop.ctx, "debugStats: %s %s\n", stats, conf)
		}()
	}

	var err error
	var cha *Channel
	cha, err = cop.getChannel(false)
	if err != nil {
		return err
	}

	if err := cha.cha.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	for i := 0; i < RESEND_MAX_NUM; i++ {
		err = cha.cha.Publish(
			exchange,   // publish to an exchange
			routingKey, // routing to 0 or more queues
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            []byte(body),
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        0,              // 0-9
				// a bunch of application/implementation-specific fields
			},
		)
		if err == nil {
			break
		}
		cop.probeCloseChannel(cha)
	}

	if err != nil {
		return WrapError(err, "Failed to publish a message")
	}

	// put current channel into idle pool
	cop.putChannel(cha)

	return err
}
func (cop *ConnPool) Consumer(queueName, ctag string, handle func(deliveries <-chan amqp.Delivery, done chan error), queueArgs amqp.Table) (*Connection, error) {
	if cop.conf.Debug {
		defer func() {
			stats, _ := json.Marshal(cop.Stats())
			conf, _ := json.Marshal(cop.conf)
			g.Log().Printf(cop.ctx, "debugStats: %s %s\n", stats, conf)
		}()
	}

	var err error
	var conn *Connection
	conn, err = cop.getConn(false)
	if err != nil {
		return nil, err
	}

	var channel *Channel
	channel, err = cop.getChannel(true)
	if err != nil {
		return nil, err
	}

	queue, err := channel.cha.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		queueArgs, // arguments
	)
	if err != nil {
		return nil, WrapError(err, "Queue Declare error !")
	}
	// 新增限流
	err = channel.cha.Qos(500, 0, false)
	if err != nil {
		return nil, err
	}
	deliveries, err := channel.cha.Consume(
		queue.Name, // name
		ctag,       // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, WrapError(err, "Consume error !")
	}
	handle(deliveries, conn.done)

	//cop.probeCloseChannel(channel)
	// put current channel into idle pool
	cop.putChannel(channel)
	return conn, nil
}

func (cop *ConnPool) ConsumerHash(queueName, ctag string, handle func(deliveries <-chan amqp.Delivery, done chan error)) (*Connection, error) {
	if cop.conf.Debug {
		defer func() {
			stats, _ := json.Marshal(cop.Stats())
			conf, _ := json.Marshal(cop.conf)
			g.Log().Printf(cop.ctx, "debugStats: %s %s\n", stats, conf)
		}()
	}

	var err error
	var conn *Connection
	conn, err = cop.getConn(false)
	if err != nil {
		return nil, err
	}

	var channel *Channel
	channel, err = cop.getChannel(true)
	if err != nil {
		return nil, err
	}

	// 新增限流
	err = channel.cha.Qos(500, 0, false)
	if err != nil {
		return nil, err
	}
	deliveries, err := channel.cha.Consume(
		queueName, // name
		ctag,      // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, WrapError(err, "Consume error !")
	}
	handle(deliveries, conn.done)

	//cop.probeCloseChannel(channel)
	// put current channel into idle pool
	cop.putChannel(channel)
	return conn, nil
}

func (cop *ConnPool) DeclareConsumer(exchange, exchangeType, queueName, bindingKey, ctag string, handle func(deliveries <-chan amqp.Delivery)) error {
	if cop.conf.Debug {
		defer func() {
			stats, _ := json.Marshal(cop.Stats())
			conf, _ := json.Marshal(cop.conf)
			g.Log().Printf(cop.ctx, "debugStats: %s %s\n", stats, conf)
		}()
	}

	var err error
	var channel *Channel
	channel, err = cop.getChannel(false)
	if err != nil {
		return err
	}

	err = channel.cha.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return WrapError(err, "Exchange Declare error !")
	}
	queue, err := channel.cha.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return WrapError(err, "Queue Declare error !")
	}

	err = channel.cha.QueueBind(
		queue.Name, // name of the queue
		bindingKey, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return WrapError(err, "Queue Bind error !")
	}
	deliveries, err := channel.cha.Consume(
		queue.Name, // name
		ctag,       // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return WrapError(err, "Consume error !")
	}
	go handle(deliveries)

	return nil
}
