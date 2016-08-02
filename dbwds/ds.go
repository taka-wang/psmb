// Package dbwds an in-memory data store for writer.
//
// By taka@cmwang.net
//
package dbwds

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	// RedisPool redis connection pool
	RedisPool *redis.Pool
	hostName  string
)

func init() {
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}

	RedisPool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0, // When zero, there is no limit on the number of connections in the pool.
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", hostName+":6379")
			if err != nil {
				fmt.Println(err.Error())
			}
			return conn, err
		},
	}
}

// @Implement IWriterTaskDataStore contract implicitly

// writerTaskDataStore write task map type
type writerTaskDataStore struct {
	redis redis.Conn
}

// NewDataStore instantiate mbtcp write task map
func NewDataStore(conf map[string]string) (interface{}, error) {
	conn := RedisPool.Get()
	if nil == conn {
		return nil, errors.New("Connect redis failed")
	}

	return &writerTaskDataStore{
		redis: conn,
	}, nil
}

// Add add request to write task map
func (ds *writerTaskDataStore) Add(tid, cmd string) {
	_, err := ds.redis.Do("HSET", "mbtcp:writer", tid, cmd)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
}

// Get get request from write task map
func (ds *writerTaskDataStore) Get(tid string) (string, bool) {
	ret, err := redis.String(ds.redis.Do("HGET", "mbtcp:writer", tid))
	if err != nil {
		fmt.Println("redis get failed:", err)
		return "", false
	}
	return ret, true
}

// Delete remove request from write task map
func (ds *writerTaskDataStore) Delete(tid string) {
	_, err := ds.redis.Do("HDEL", "mbtcp:writer", tid)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
}
