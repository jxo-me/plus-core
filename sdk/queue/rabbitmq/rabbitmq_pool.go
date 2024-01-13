package rabbitmq

import (
	"context"
	"github.com/jxo-me/rabbitmq-go"
	"sync"
	"time"
)

// ConnectionPool 是一个带有计数器和最大连接数限制的连接池
type ConnectionPool struct {
	mu                sync.Mutex
	cond              *sync.Cond
	connections       []*connectionWithCounter
	url               string
	maxConnections    int
	reconnectInterval int
	logger            rabbitmq.Logger
	config            rabbitmq.Config
}

type connectionWithCounter struct {
	conn     *rabbitmq.Conn
	counter  int
	isInUse  bool
	lastUsed time.Time
}

// NewConnectionPool 创建一个新的带有计数器和最大连接数限制的连接池
func NewConnectionPool(ctx context.Context, url string, maxConnections, reconnectInterval int,
	log rabbitmq.Logger, cfg *rabbitmq.Config) (*ConnectionPool, error) {

	pool := &ConnectionPool{
		connections:       make([]*connectionWithCounter, 0),
		url:               url,
		maxConnections:    maxConnections,
		reconnectInterval: reconnectInterval,
		logger:            log,
	}
	if cfg != nil {
		pool.config = *cfg
	}

	return pool, nil
}

// GetConnection 从连接池获取一个连接
func (pool *ConnectionPool) GetConnection(ctx context.Context) (*rabbitmq.Conn, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// 检查连接池是否已达到最大连接数限制
	if len(pool.connections) >= pool.maxConnections {
		for len(pool.connections) >= pool.maxConnections {
			// 选择计数器最小且未被占用的连接
			selectedConn := pool.connections[0]
			for _, conn := range pool.connections {
				if !conn.isInUse && conn.counter < selectedConn.counter {
					selectedConn = conn
				}
			}

			selectedConn.isInUse = true
			selectedConn.counter++
			selectedConn.lastUsed = time.Now()

			return selectedConn.conn, nil
		}
	}

	conn, err := pool.createConnection(ctx)
	if err != nil {
		return nil, err
	}
	conn.isInUse = true
	conn.counter++
	conn.lastUsed = time.Now()
	pool.connections = append(pool.connections, conn)
	return conn.conn, nil
}

// ReleaseConnection 将连接释放回连接池
func (pool *ConnectionPool) ReleaseConnection(conn *rabbitmq.Conn) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for _, c := range pool.connections {
		if c.conn == conn {
			c.isInUse = false
			c.lastUsed = time.Now()
			return
		}
	}
}

// createConnection 创建一个新的连接
func (pool *ConnectionPool) createConnection(ctx context.Context) (*connectionWithCounter, error) {
	conn, err := rabbitmq.NewConn(
		ctx,
		pool.url,
		rabbitmq.WithConnectionOptionsLogger(pool.logger),
		rabbitmq.WithConnectionOptionsConfig(pool.config),
		rabbitmq.WithConnectionOptionsReconnectInterval(time.Duration(pool.reconnectInterval)*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return &connectionWithCounter{
		conn:     conn,
		counter:  0,
		isInUse:  false,
		lastUsed: time.Now(),
	}, nil
}

func (pool *ConnectionPool) Close(ctx context.Context) error {
	for _, conn := range pool.connections {
		_ = conn.conn.Close(ctx)
	}
	return nil
}
