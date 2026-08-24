package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const refreshTokenTTL = 30 * 24 * time.Hour

type AuthServiceInterface interface {
	CreateToken(ctx context.Context, userID int, username string) (string, error)
	Validate(ctx context.Context, token string) (int, error)
	Revoke(ctx context.Context, token string) error
}

type AuthService struct {
	rdb *redis.Client
}

func NewAuthService(rdb *redis.Client) *AuthService {
	return &AuthService{rdb: rdb}
}

func generateRefreshToken() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func (auth *AuthService) CreateToken(ctx context.Context, userID int, username string) (string, error) {
	token := generateRefreshToken()

	err := auth.rdb.Set(ctx, "refresh:"+token, strconv.Itoa(userID), refreshTokenTTL).Err()
	return token, err
}

func (auth *AuthService) Validate(ctx context.Context, token string) (int, error) {
	val, err := auth.rdb.Get(ctx, "refresh:"+token).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (auth *AuthService) Revoke(ctx context.Context, token string) error {
	return auth.rdb.Del(ctx, "refresh:"+token).Err()
}
