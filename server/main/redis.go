package main

import "github.com/go-redis/redis/v8"

var rdb *redis.Client

func initPool(addr string, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
