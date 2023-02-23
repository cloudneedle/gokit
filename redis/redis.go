package redis

import "github.com/redis/go-redis/v9"

type Client struct {
	*redis.Client
}

type Options struct {
	Addr     string
	Password string
	DB       int
}

const Nil = redis.Nil

func NewClient(options Options) *Client {
	return &Client{redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})}
}
