package utils

import (
	"log"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
	redisConf RedisWriter
)

/**
* NewPool
* Prepares a pool of redis connection
 */

func (r RedisWriter) NewPool() *redis.Pool {
	config := r
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

func (r RedisWriter) InitRedisPool() {
	redisPool = r.NewPool()
}

type RedisWriter struct {
	RedisAddress         string
	RedisPoolMaxIdle     int
	RedisPoolMaxActive   int
	RedisPoolIdleTimeout int
	RedisPassword        string
	//RedisPtr             *Redis
}

/**
* Write
* split the string on . and make domin as key and full url as value to Links hash in redis
 */
func (r RedisWriter) Write(datapipe <-chan string) {
	rw := NewRedis()
	for data := range datapipe {
		key := strings.Split(data, ".")
		var title string
		if len(key)>1{
			title=key[0]+key[1]
		}else{
			title=key[0]
		}
		err := rw.Hset("Links", title, []byte(data))
		if err != nil {
			log.Fatalln("Unable to write in redis")
		}
	}
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
