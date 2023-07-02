package boards

import (
	"context"
	"fmt"
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
		Create(ctx context.Context, board Board) (Board, error)
		BoardOfId(ctx context.Context, id string) (Board, error)
		BoardsOfUser(ctx context.Context, userId int64) ([]Board, error)
	}

	RedisRepository struct {
		redis redis.Service
	}
)

func NewRepository(redis redis.Service) Repository {
	return &RedisRepository{
		redis,
	}
}

func (brr *RedisRepository) Create(ctx context.Context, board Board) (Board, error) {
	_, err := brr.redis.TxPipelined(ctx, func(pipe redis.Pipe) error {
		data, err := json.ToJson(board)
		if err != nil {
			return err
		}

		pipe.Set(ctx, boardKey(board.Id), data, Infinity)
		pipe.SAdd(ctx, userBoardsKey(board.OwnerId), board.Id)
		return nil
	})
	if err != nil {
		return Board{}, err
	}
	return board, nil
}

func (brr *RedisRepository) BoardOfId(ctx context.Context, id string) (Board, error) {
	data, err := brr.redis.Get(ctx, boardKey(id))
	if err != nil {
		return Board{}, err
	}

	var board Board
	if json.FromJson([]byte(data), &board) != nil {
		return Board{}, err
	}

	return board, nil
}

func (brr *RedisRepository) BoardsOfUser(ctx context.Context, userId int64) ([]Board, error) {
	boards, err := brr.redis.SMembers(ctx, userBoardsKey(userId))
	if err != nil {
		return nil, err
	}

	var userBoards []Board
	for _, boardId := range boards {
		board, err := brr.BoardOfId(ctx, boardId)
		if err != nil {
			return nil, err
		}

		userBoards = append(userBoards, board)
	}

	return userBoards, nil
}

func boardKey(boardId string) string {
	return fmt.Sprintf("boards:%s", boardId)
}

func userBoardsKey(userId int64) string {
	return fmt.Sprintf("users:%s:boards", strutil.Int64ToString(userId))
}
