package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/microservices-simulator-api/internal/utils/hashutil"
	"github.com/microservices-simulator-api/internal/utils/json"
	"github.com/microservices-simulator-api/internal/utils/redis"
	"github.com/microservices-simulator-api/internal/utils/strutil"
	"time"
)

var (
	Infinity time.Duration = 0
)

type (
	Repository interface {
		Create(ctx context.Context, user User) (User, error)
		UserOfId(ctx context.Context, id int64) (User, error)
		UserOfUsername(ctx context.Context, username string) (User, error)
	}

	RedisRepository struct {
		redis redis.Service
	}
)

func NewRepository(redis redis.Service) *RedisRepository {
	return &RedisRepository{
		redis,
	}
}

func (urr *RedisRepository) Create(ctx context.Context, user User) (User, error) {
	_, err := urr.redis.TxPipelined(ctx, func(pipe redis.Pipe) error {
		data, err := json.ToJson(user)
		if err != nil {
			return err
		}

		pipe.Set(ctx, userKey(user.Id), data, Infinity)
		pipe.Set(ctx, usernameKey(user.Email), strutil.Int64ToString(user.Id), Infinity)
		return nil
	})

	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (urr *RedisRepository) UserOfId(ctx context.Context, id int64) (User, error) {
	data, err := urr.redis.Get(ctx, userKey(id))
	if err != nil {
		return User{}, err
	}

	var user User
	if err = json.FromJson([]byte(data), &user); err != nil {
		return User{}, err
	}

	return user, nil
}

func (urr *RedisRepository) UserOfUsername(ctx context.Context, username string) (User, error) {
	id, err := urr.redis.Get(ctx, usernameKey(username))
	if err != nil {
		return User{}, err
	} else if id == "" {
		return User{}, errors.New("user not found")
	}

	id64, err := strutil.StringToInt64(id)
	if err != nil {
		return User{}, err
	}

	return urr.UserOfId(ctx, id64)
}

func userKey(userId int64) string {
	return fmt.Sprintf("users:%d", userId)
}

func usernameKey(username string) string {
	return fmt.Sprintf("usernames:%s", hashutil.Sha256(username))
}
