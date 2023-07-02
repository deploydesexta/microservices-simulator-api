package shapes

import (
	"context"
	"fmt"
	"github.com/microservices-simulator-api/internal/utils/json"
	"github.com/microservices-simulator-api/internal/utils/redis"
	"github.com/microservices-simulator-api/internal/utils/strutil"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

const (
	AllShapesSetKey = "shapes:all"
)

var (
	Infinity time.Duration = 0
)

type (
	Repository interface {
		AllShapes(ctx context.Context) ([]Shape, error)
		Create(ctx context.Context, shape Shape) (Shape, error)
		ShapeOfId(ctx context.Context, id int64) (Shape, error)
		StoreImage(ctx context.Context, name string, content []byte) (string, error)
	}

	RedisRepository struct {
		blobStore sync.Map
		redis     redis.Service
	}
)

func NewRepository(redis redis.Service) *RedisRepository {
	return &RedisRepository{
		blobStore: sync.Map{},
		redis:     redis,
	}
}

func (urr *RedisRepository) AllShapes(ctx context.Context) ([]Shape, error) {
	data, err := urr.redis.SMembers(ctx, AllShapesSetKey)
	if err != nil {
		return []Shape{}, err
	}

	g, ctx := errgroup.WithContext(ctx)

	shapesC := make(chan Shape, len(data))

	GetShape := func(shapeId string) func() error {
		return func() error {
			i, err := strutil.StringToInt64(shapeId)
			if err != nil {
				return err
			}

			shape, err := urr.ShapeOfId(ctx, i)
			if err != nil {
				return err
			}

			shapesC <- shape
			return nil
		}
	}

	for _, shapeId := range data {
		g.Go(GetShape(shapeId))
	}

	if err = g.Wait(); err != nil {
		return []Shape{}, err
	}

	shapes := make([]Shape, 0, len(data))
	for i := 0; i < len(data); i++ {
		shapes = append(shapes, <-shapesC)
	}

	return shapes, nil
}

func (urr *RedisRepository) Create(ctx context.Context, shape Shape) (Shape, error) {
	_, err := urr.redis.TxPipelined(ctx, func(pipe redis.Pipe) error {
		data, err := json.ToJson(shape)
		if err != nil {
			return err
		}

		pipe.Set(ctx, shapeKey(shape.Id), data, Infinity)
		pipe.SAdd(ctx, AllShapesSetKey, shape.Id)
		return nil
	})
	if err != nil {
		return Shape{}, err
	}

	return shape, nil
}

func (urr *RedisRepository) ShapeOfId(ctx context.Context, id int64) (Shape, error) {
	data, err := urr.redis.Get(ctx, shapeKey(id))
	if err != nil {
		return Shape{}, err
	}

	var shape Shape
	if err = json.FromJson([]byte(data), &shape); err != nil {
		return Shape{}, err
	}

	return shape, nil
}

func (urr *RedisRepository) StoreImage(ctx context.Context, name string, content []byte) (string, error) {
	urr.blobStore.Store(name, content)
	return name, nil
}

func shapeKey(shapeId int64) string {
	return fmt.Sprintf("shapes:%d", shapeId)
}
