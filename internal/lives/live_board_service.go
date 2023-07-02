package lives

import (
	"fmt"
	"github.com/microservices-simulator-api/internal/boards"
	"github.com/microservices-simulator-api/internal/users"
	"github.com/microservices-simulator-api/internal/utils/jwtutil"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type LiveBoardService struct {
	lives  *LiveBoardManager
	boards boards.Service
	users  *users.UserService
}

func NewLiveBoardService(
	lives *LiveBoardManager,
	boards boards.Service,
	users *users.UserService,
) *LiveBoardService {
	return &LiveBoardService{
		lives:  lives,
		boards: boards,
		users:  users,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (lbm *LiveBoardService) NewBoardConnection(c echo.Context) error {
	fmt.Printf("Attempting to connect to a board...\n")
	ctx := c.Request().Context()
	boardId := c.Param("id")

	user := c.Get("user").(*jwtutil.UserClaims)
	if user == nil {
		return echo.ErrUnauthorized
	}

	board, err := lbm.boards.BoardOfId(ctx, boardId)
	if err != nil {
		return err
	}

	// Upgrade HTTP request to WebSocket
	fmt.Printf("Attempting to connect [board:%s]...\n", boardId)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer ws.Close()

	fmt.Printf("WebSocket connected [board:%s]...\n", boardId)
	return lbm.lives.ManageBoard(user.Id, board, NewWebSocketChannel(ws))
}
