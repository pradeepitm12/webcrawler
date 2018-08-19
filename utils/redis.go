package utils

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
)

func NewPool() *redis.Pool {
	config := RedisWriter{}
	return &redis.Pool{
		MaxIdle:     config.RedisPoolMaxIdle,
		MaxActive:   config.RedisPoolMaxActive,
		IdleTimeout: time.Duration(config.RedisPoolIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			//return redis.Dial("tcp", config.RedisAddress) },
			c, err := redis.Dial("tcp", config.RedisAddress)
			if err != nil {
				return nil, err
			}
			c.Do("AUTH", config.RedisPassword)
			return c, err
		},
	}
}

func InitRedisPool() {
	config := RedisWriter{}
	log.Println("Redis *********")
	log.Printf("Init Redis Connection pool with params: RedisAddress = %s, MaxIdle = %d, MaxActive = %d, IdleTimeout = %d(s)",
		config.RedisAddress, config.RedisPoolMaxIdle, config.RedisPoolMaxActive, config.RedisPoolIdleTimeout)
	redisPool = NewPool()
}

/**
* Redis Commands
* Hset - to set hash in redis
 */

type Redis struct {
	conn redis.Conn
}

func NewRedis() *Redis {
	conn := redisPool.Get()
	log.Printf("Get conn from redis_conn_pool, active connections in the pool is %d", redisPool.ActiveCount())
	//conn.Do("AUTH", "redis@123@Azure")
	return &Redis{
		conn: conn,
	}
}

func (s *Redis) Close() {
	s.conn.Close()
	log.Printf("Return conn to redis_conn_pool, active connections in the pool is %d", redisPool.ActiveCount())
}

func (s *Redis) Hset(key, field string, value []byte) error {
	log.Printf("Redis HSet: %s %s %#v", key, field, value)
	_, err := s.conn.Do("HSET", key, field, value)
	return err
}
