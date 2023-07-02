package lives

import (
	"github.com/microservices-simulator-api/internal/boards"
	"sync"
)

type LiveBoardManager struct {
	sync.Mutex
	boards map[string]*LiveBoard
}

func NewLiveBoardManager() *LiveBoardManager {
	return &LiveBoardManager{
		boards: make(map[string]*LiveBoard),
	}
}

func (lbm *LiveBoardManager) ManageBoard(userId int64, board boards.Board, channel Channel) error {
	liveBoard, ok := lbm.boards[board.Id]
	if !ok {
		liveBoard = lbm.createLiveBoard(board)
	}

	return liveBoard.NewClientConnection(userId, channel)
}

func (lbm *LiveBoardManager) createLiveBoard(board boards.Board) *LiveBoard {
	lbm.Lock()
	liveBoard := NewLiveBoard(board)
	lbm.boards[board.Id] = liveBoard
	lbm.Unlock()
	return liveBoard
}
