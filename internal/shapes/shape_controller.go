package shapes

import (
	"github.com/labstack/echo/v4"
	"github.com/microservices-simulator-api/internal/utils/strutil"
	"io"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{
		service,
	}
}

func (bc *Controller) AllShapes(c echo.Context) error {
	ctx := c.Request().Context()

	shapes, err := bc.service.AllShapes(ctx)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(200, shapes)
}

func (bc *Controller) CreateShape(c echo.Context) error {
	ctx := c.Request().Context()

	image, err := c.FormFile("image")
	if err != nil {
		return err
	}

	is, err := image.Open()
	if err != nil {
		return err
	}
	defer is.Close()

	ibs, err := io.ReadAll(is)
	if err != nil {
		return err
	}

	var input ShapeInput
	if err := c.Bind(&input); err != nil {
		return err
	}

	shape, err := bc.service.Create(ctx, input, ibs)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(201, shape)
}

func (bc *Controller) ShapeOfId(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	if id == "" {
		return c.JSON(400, map[string]string{
			"message": "Param id is required",
		})
	}

	uid, err := strutil.StringToInt64(id)
	if err != nil {
		return c.JSON(400, map[string]string{
			"message": "Param id must be an integer",
		})
	}

	shape, err := bc.service.ShapeOfId(ctx, uid)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(200, shape)
}
