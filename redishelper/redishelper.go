package redishelper

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	redsync "gopkg.in/redsync.v1"
)

const (
	maxIdleConns = 400
	idleTimeout  = 240 * time.Second

	lockExpriyTime     = 50 * time.Second
	retryDelayTime     = 500 * time.Millisecond
	bufferTimeForError = 5 * time.Second
)

type Pool struct {
	pool redis.Pool
}

type Client struct {
	redis.Conn
}

type RedisPool interface {
	Get() redis.Conn
	GetDB(databaseID int) RedisClient
}

type RedisClient interface {
	redis.Conn
	Get() redis.Conn
}

func NewRedisPool(server string) *Pool {
	return &Pool{
		pool: redis.Pool{
			MaxIdle:     maxIdleConns,
			IdleTimeout: idleTimeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func (p *Pool) GetDB(databaseID int) RedisClient {
	conn := p.Get()
	conn.Do("SELECT", databaseID)
	return &Client{
		Conn: conn,
	}
}

func (p *Pool) Get() redis.Conn {
	return p.pool.Get()
}

func (c *Client) Get() redis.Conn {
	return c.Conn
}

type locker struct {
	mutex *redsync.Mutex
}

type SharedLock interface {
	Unlock()
}

func NewLock(redis RedisPool, name string) (SharedLock, error) {
	mutex := redsync.New(
		[]redsync.Pool{redis}).NewMutex(
		name,
		redsync.SetExpiry(lockExpriyTime),
		redsync.SetTries(int((lockExpriyTime-bufferTimeForError)/retryDelayTime)),
		redsync.SetRetryDelay(retryDelayTime),
	)
	fmt.Println("retry ", int((lockExpriyTime-bufferTimeForError)/retryDelayTime))
	err := mutex.Lock()
	return &locker{
		mutex: mutex,
	}, err
}

func (l *locker) Unlock() {
	l.mutex.Unlock()
}
