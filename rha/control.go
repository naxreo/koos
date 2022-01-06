package rha

import (
	"github.com/go-redis/redis/v8"
)

type RedisPod struct {
	Role   string
	Init   bool
	Addr   string
	Passwd string
	Db     int
}

func Connect(i *RedisPod) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     i.Addr,
		Password: i.Passwd,
		DB:       i.Db,
	})
	return rdb
}
