package cache

import (
	"errors"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

const (
	timeout = 2
)

var (
	pool *redis.Pool
)

// Setup initialises redis
func Setup(addr string, auth string) {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err == nil && auth != "" {
				_, err = c.Do("AUTH", auth)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Process does some funky caching stuff
func Process(url string, regenerate func() ([]byte, error)) ([]byte, error) {

	key := "lobby:c:" + url

	// If we haven't yet expired return data
	if !checkExpired(url) {
		// Grab data
		data, err := get(key)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	setExpiry(url)
	regen, err := regenerate()
	if err != nil {
		return nil, err
	}

	err = put(key, regen)
	if err != nil {
		return nil, err
	}

	setExpiry(url)
	return regen, nil
}

func setExpiry(url string) error {
	key := "lobby:c:" + url + ":expiry"
	to := time.Now().Unix() + timeout
	return put(key, []byte(strconv.FormatInt(to, 10)))
}

func checkExpired(url string) bool {
	key := "lobby:c:" + url + ":expiry"
	resp, err := get(key)
	if err != nil {
		return true
	}

	s := string(resp[:])
	expiry, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return true
	}

	if time.Now().Unix() > expiry {
		return true
	}

	return false
}

func get(key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("HGET", key, "response"))
}

func put(key string, response []byte) error {
	conn := pool.Get()
	defer conn.Close()

	r, err := conn.Do("HDEL", key, "response")
	if err != nil {
		return err
	}

	if r == 0 {
		return errors.New("Error adding response to redis")
	}

	r, err = conn.Do("HSET", key, "response", response)
	if err != nil {
		return err
	}

	if r == 0 {
		return errors.New("Error adding response to redis")
	}

	return nil

}
