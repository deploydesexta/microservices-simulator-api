package boards

import (
	"github.com/labstack/echo/v4"
	"github.com/microservices-simulator-api/internal/utils/jwtutil"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{
		service,
	}
}

func (bc *Controller) CreateBoard(c echo.Context) error {
	ctx := c.Request().Context()

	user := c.Get("user").(*jwtutil.UserClaims)
	if user == nil {
		return echo.ErrUnauthorized
	}

	board, err := bc.service.NewBoard(ctx, user.Id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(200, board)
}

func (bc *Controller) BoardOfId(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return c.JSON(400, map[string]string{
			"message": "Param id is required",
		})
	}

	board, err := bc.service.BoardOfId(ctx, id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(200, board)
}

func (bc *Controller) BoardsOfUser(c echo.Context) error {
	ctx := c.Request().Context()

	user := c.Get("user").(*jwtutil.UserClaims)
	if user == nil {
		return echo.ErrUnauthorized
	}

	boards, err := bc.service.BoardsOfUser(ctx, user.Id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(200, boards)
}
