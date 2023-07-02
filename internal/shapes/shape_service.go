package shapes

import "context"

type (
	Service interface {
		AllShapes(ctx context.Context) ([]Shape, error)
		Create(ctx context.Context, shape ShapeInput, image []byte) (Shape, error)
		ShapeOfId(ctx context.Context, id int64) (Shape, error)
	}

	ShapeService struct {
		repository Repository
	}
)

func NewService(repository Repository) Service {
	return &ShapeService{
		repository,
	}
}

func (ss *ShapeService) AllShapes(ctx context.Context) ([]Shape, error) {
	return ss.repository.AllShapes(ctx)
}

func (ss *ShapeService) Create(ctx context.Context, shape ShapeInput, image []byte) (Shape, error) {
	filepath, err := ss.repository.StoreImage(ctx, shape.Name, image)
	if err != nil {
		return Shape{}, err
	}

	return ss.repository.Create(ctx, NewShapeFromInput(shape, filepath))
}

func (ss *ShapeService) ShapeOfId(ctx context.Context, id int64) (Shape, error) {
	return ss.repository.ShapeOfId(ctx, id)
}
