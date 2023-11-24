package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

const (
	USER_JWT_TOKEN_KEY_TEMPLATE         = "%d_jwt_expired_at"
	EMAIL_JWT_TOKEN_KEY_PREFIX          = "email_%s"
	DEFAULT_EMAIL_REGISTER_TOKEN_EXPIRE = 15 * time.Minute
)

type JWTCache struct {
	logger  *zap.SugaredLogger
	cache   *redis.Client
	context context.Context
}

func NewJWTCache(cache *redis.Client, logger *zap.SugaredLogger) *JWTCache {
	return &JWTCache{
		logger:  logger,
		cache:   cache,
		context: context.Background(),
	}
}

func (c *JWTCache) CleanUserJWTTokenExpiredAt(user *User) error {
	key := fmt.Sprintf(USER_JWT_TOKEN_KEY_TEMPLATE, user.ExportID())
	return c.cache.Del(c.context, key).Err()
}

func (c *JWTCache) InitUserJWTTokenExpiredAt(user *User, jwtTokenExpireAt string) error {
	key := fmt.Sprintf(USER_JWT_TOKEN_KEY_TEMPLATE, user.ExportID())
	return c.cache.Set(c.context, key, jwtTokenExpireAt, 0).Err()
}

func (c *JWTCache) DoesUserJWTTokenAvaliable(user *User, jwtTokenExpireAt string) (bool, error) {
	key := fmt.Sprintf(USER_JWT_TOKEN_KEY_TEMPLATE, user.ExportID())
	expireAtInCache, errInGet := c.cache.Get(c.context, key).Result()

	// check error
	if errInGet == redis.Nil {
		return false, nil
	} else if errInGet != nil {
		return false, errInGet
	}
	// check expire
	return jwtTokenExpireAt >= expireAtInCache, nil
}

func (c *JWTCache) SetTokenForEmail(email string, jwtToken string) error {
	key := fmt.Sprintf(EMAIL_JWT_TOKEN_KEY_PREFIX, email)
	return c.cache.Set(c.context, key, jwtToken, DEFAULT_EMAIL_REGISTER_TOKEN_EXPIRE).Err()
}

func (c *JWTCache) GetTokenByEmail(email string) (string, error) {
	key := fmt.Sprintf(EMAIL_JWT_TOKEN_KEY_PREFIX, email)
	jwtTokenInCache, errInGet := c.cache.Get(c.context, key).Result()

	// check error
	if errInGet == redis.Nil {
		return "", errors.New("empty jwt token")
	} else if errInGet != nil {
		return "", errInGet
	}
	// check expire
	return jwtTokenInCache, nil
}
