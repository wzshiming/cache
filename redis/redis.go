package redis_cache

import (
	"crypto/tls"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/wzshiming/cache"
)

type Redis struct {
	cli *redis.Client
	u   string
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
	if ur.Scheme == "rediss" {
		h, _, _ := net.SplitHostPort(ur.Host)
		op.TLSConfig = &tls.Config{ServerName: h}

	}
	op.Addr = ur.Host
	cli := redis.NewClient(op)

	return &Redis{
		cli: cli,
		u:   u,
	}, nil
}

func (rc *Redis) Scan(key string, val interface{}) error {
	//	v, ok := val.(encoding.BinaryUnmarshaler)
	//	if !ok {
	v := &cache.Unmarshaler{val}
	//	}
	err := rc.cli.Get(key).Scan(v)
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
	//	v, ok := val.(encoding.BinaryMarshaler)
	//	if !ok {
	v := &cache.Marshaler{val}
	//	}
	_, err := rc.cli.Set(key, v, timeout).Result()
	return err
}

func (rc *Redis) Delete(key string) error {
	_, err := rc.cli.Del(key).Result()
	return err
}

func (rc *Redis) IsExist(key string) bool {
	d, _ := rc.cli.Exists(key).Result()
	return d != 0
}
