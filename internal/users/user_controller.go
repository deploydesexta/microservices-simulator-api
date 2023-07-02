package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

func (bc *Controller) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var input NewUserInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	user, err := bc.service.Create(ctx, input)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(http.StatusOK, user)
}
