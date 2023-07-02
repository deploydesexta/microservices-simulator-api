package authentication

import (
	"github.com/labstack/echo/v4"
	"github.com/microservices-simulator-api/internal/config"
	"github.com/microservices-simulator-api/internal/utils/jwtutil"
	"net/http"
)

type Controller struct {
	cfg     config.SecurityConfig
	manager Manager
}

func NewController(cfg config.SecurityConfig, manager Manager) *Controller {
	return &Controller{cfg, manager}
}

func (ctrl *Controller) Login(c echo.Context) error {
	ctx := c.Request().Context()
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return echo.ErrUnauthorized
	}

	token, err := ctrl.manager.Authenticate(ctx, username, password, c.RealIP())
	if err != nil {
		return echo.ErrUnauthorized
	}

	c.SetCookie(&http.Cookie{
		Name:     ctrl.cfg.TokenName,
		Value:    token,
		MaxAge:   ctrl.cfg.TokenMaxAge,
		Path:     "/",
		HttpOnly: true,
		Domain:   ctrl.cfg.TokenDomain,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": token,
	})
}

func (ctrl *Controller) Me(c echo.Context) error {
	user := c.Get("user").(*jwtutil.UserClaims)
	if user == nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, user)
}
