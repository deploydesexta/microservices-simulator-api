package authentication

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/microservices-simulator-api/internal/utils/redis"
	"github.com/microservices-simulator-api/internal/utils/strutil"
	"net/http"
	"time"
)

var (
	MaxAttemptsError     = errors.New("max attempts exceeded")
	Unauthorized         = http.StatusUnauthorized
	UnauthorizedResponse = map[string]string{
		"message": "unauthorized",
	}
)

type (
	Manager interface {
		OnBeforeAuthentication(c echo.Context) error
		Authenticate(ctx context.Context, username, password, ip string) (string, error)
		OnAuthenticationError(ctx context.Context, ip string)
	}

	AuthenticationManager struct {
		provider Provider
		redis    redis.Service
	}
)

func NewManager(provider Provider, redis redis.Service) *AuthenticationManager {
	return &AuthenticationManager{provider, redis}
}

func (am *AuthenticationManager) Authenticate(ctx context.Context, username, password, ip string) (string, error) {
	token, err := am.provider.Authenticate(ctx, username, password)
	if err != nil {
		am.OnAuthenticationError(ctx, ip)
		return "", echo.NewHTTPError(Unauthorized, UnauthorizedResponse)
	}

	return token, nil
}

func (am *AuthenticationManager) OnAuthenticationError(ctx context.Context, ip string) {
	_, err := am.redis.Pipelined(ctx, func(pipe redis.Pipe) error {
		key := ipsAuthKey(ip)
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Hour)
		return nil
	})

	if err != nil {
		log.Warnf("error suppressed on AuthenticationManager.OnAuthenticationError: %s", err.Error())
	}
}

func (am *AuthenticationManager) OnBeforeAuthentication(c echo.Context) error {
	ctx := c.Request().Context()
	key := ipsAuthKey(c.RealIP())

	sAtts, err := am.redis.Get(ctx, key)
	if err != nil {
		return err
	} else if sAtts == "" {
		return nil
	}

	atts, err := strutil.StringToInt(sAtts)
	if err != nil {
		return err
	}

	if atts > 5 {
		return MaxAttemptsError
	}

	return nil
}

func ipsAuthKey(ip string) string {
	return "security:" + ip + ":attempts"
}
