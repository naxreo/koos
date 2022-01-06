package rha

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func RhaPing(client *redis.Client, tmout int) string {
	/*
		return PONG: Connection success
		return AUTH: Connected but Password failed
		return FAILED: Connection failed
	*/
	ctx := context.Background()
	retval := "TIMEOUT"

	for i := 0; i < tmout; i++ {
		p, err := client.Ping(ctx).Result()
		if p == "PONG" {
			return p
		}

		if err != nil {
			fmt.Printf("redis connect failed %d times, %s\n", tmout, err.Error())
			switch err.Error() {
			case "dial tcp: i/o timeout":
				retval = "TIMEOUT"
			default:
				retval = "TIMEOUT"
			}
		}
		time.Sleep(time.Second * 1)
	}

	ctx.Done()
	return retval
}
