// exists, _ := redis.RDB.Exists(redis.Ctx, key).Result()
// redis.RDB.HSet(redis.Ctx, key, "password", hashedPwd)
// redis.RDB.SAdd(redis.Ctx, "users", user.Username)
// storedPwd, err := redis.RDB.HGet(redis.Ctx, key, "password").Result()
// err = redis.RDB.Set(redis.Ctx, "session:"+token, user.Username, 24*time.Hour).Err()

package repositories

import (
	"context"
	"time"

	"github.com/haileamlak/chat-system/infrastructure"
)

type UserRepository interface {
	UserExists(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, username string, passwordHash string) error
	GetUserPassword(ctx context.Context, username string) (string, error)
	SaveSession(ctx context.Context, token string, username string) error
}

type userRepository struct{
	redisService infrastructure.RedisService
}

func NewUserRepository(redisService infrastructure.RedisService) UserRepository {
	return &userRepository{
		redisService: redisService,
	}
}

func (r *userRepository) UserExists(ctx context.Context, username string) (bool, error) {
	key := "user:" + username
	exists, err := r.redisService.GetClient().Exists(ctx, key).Result()
	return exists == 1, err
}

func (r *userRepository) CreateUser(ctx context.Context, username string, passwordHash string) error {
	key := "user:" + username
	r.redisService.GetClient().HSet(ctx, key, "password", passwordHash)
	r.redisService.GetClient().SAdd(ctx, "users", username)
	return nil
}

func (r *userRepository) GetUserPassword(ctx context.Context, username string) (string, error) {
	key := "user:" + username
	return r.redisService.GetClient().HGet(ctx, key, "password").Result()
}

func (r *userRepository) SaveSession(ctx context.Context, token string, username string) error {
	return r.redisService.GetClient().Set(ctx, "session:"+token, username, 24*time.Hour).Err()
}
func (r *userRepository) GetUserGroups(ctx context.Context, username string) ([]string, error) {
	key := "user:" + username + ":groups"
	groups, err := r.redisService.GetClient().SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return groups, nil
}