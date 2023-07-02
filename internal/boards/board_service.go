package boards

import "context"

type (
	Service interface {
		NewBoard(ctx context.Context, ownerId int64) (Board, error)
		Create(ctx context.Context, board Board) (Board, error)
		BoardOfId(ctx context.Context, id string) (Board, error)
		BoardsOfUser(ctx context.Context, userId int64) ([]Board, error)
	}

	BoardService struct {
		repository Repository
	}
)

func NewService(repository Repository) Service {
	return &BoardService{
		repository,
	}
}

func (bs *BoardService) NewBoard(ctx context.Context, ownerId int64) (Board, error) {
	return bs.repository.Create(ctx, NewBoard(ownerId))
}

func (bs *BoardService) Create(ctx context.Context, board Board) (Board, error) {
	return bs.repository.Create(ctx, board)
}

func (bs *BoardService) BoardOfId(ctx context.Context, id string) (Board, error) {
	return bs.repository.BoardOfId(ctx, id)
}

func (bs *BoardService) BoardsOfUser(ctx context.Context, userId int64) ([]Board, error) {
	return bs.repository.BoardsOfUser(ctx, userId)
}
