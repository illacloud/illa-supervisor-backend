package model

import (
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Cache struct {
	JWTCache *JWTCache
}

func NewCache(redisDriver *redis.Client, logger *zap.SugaredLogger) *Cache {
	jwtCache := NewJWTCache(redisDriver, logger)
	return &Cache{
		JWTCache: jwtCache,
	}
}
