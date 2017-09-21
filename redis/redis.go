package redis_cache

import (
	"encoding"
	"net/url"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/wzshiming/cache"
)

type Redis struct {
	p *redis.Client
	u string
}

var _ cache.Cache = (*Redis)(nil)

func NewRedis(u string) (*Redis, error) {

	ur, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	que := ur.Query()
	op := &redis.Options{}
	for k, v := range que {
		r := v[0]
		switch k {
		case "dialTimeout":
			d, err := time.ParseDuration(r)
			if err != nil {
				return nil, err
			}
			op.DialTimeout = d
		case "readTimeout":
			d, err := time.ParseDuration(r)
			if err != nil {
				return nil, err
			}
			op.ReadTimeout = d
		case "writeTimeout":
			d, err := time.ParseDuration(r)
			if err != nil {
				return nil, err
			}
			op.WriteTimeout = d
		case "db":
			d, err := strconv.Atoi(r)
			if err != nil {
				return nil, err
			}
			op.DB = d
		case "password":
			op.Password = r
		}
	}
	op.Addr = ur.Host

	client := redis.NewClient(op)

	return &Redis{
		p: client,
		u: u,
	}, nil
}

func (rc *Redis) Scan(key string, val interface{}) error {
	v, ok := val.(encoding.BinaryUnmarshaler)
	if !ok {
		v = &cache.Unmarshaler{val}
	}
	err := rc.p.Get(key).Scan(v)
	if err != nil {
		return err
	}
	return nil
}

func (rc *Redis) Get(key string) interface{} {
	var i interface{}
	rc.Scan(key, &i)
	return i
}

func (rc *Redis) Put(key string, val interface{}, timeout time.Duration) error {
	if timeout < 0 {
		timeout = 0
	}
	v, ok := val.(encoding.BinaryMarshaler)
	if !ok {
		v = &cache.Marshaler{val}
	}
	_, err := rc.p.Set(key, v, timeout).Result()
	return err
}

func (rc *Redis) Delete(key string) error {
	_, err := rc.p.Del(key).Result()
	return err
}

func (rc *Redis) IsExist(key string) bool {
	d, _ := rc.p.Exists(key).Result()
	return d != 0
}
