package main

import (
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/microservices-simulator-api/internal/authentication"
	"github.com/microservices-simulator-api/internal/utils/jwtutil"
)

type AuthenticationMiddleware struct {
	authenticationManager authentication.Manager
	jwt                   *jwtutil.Authentication
}

func NewAuthenticationMiddleware(container *Container) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authenticationManager: container.authenticationManager,
		jwt:                   container.jwt,
	}
}

func (am *AuthenticationMiddleware) BlacklistMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := am.authenticationManager.OnBeforeAuthentication(c)
			if err != nil {
				return echo.NewHTTPError(authentication.Unauthorized, err.Error())
			} else {
				return next(c)
			}
		}
	}
}

func (am *AuthenticationMiddleware) JwtMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:     "userToken",
		ErrorHandler:   am.onAuthenticationError,
		KeyFunc:        am.jwt.GetPublicKey,
		NewClaimsFunc:  func(c echo.Context) jwtutil.Claims { return am.jwt.NewClaims() },
		TokenLookup:    "cookie:accessToken,header:Authorization:Bearer ",
		SuccessHandler: am.onAuthenticationSuccess,
	})
}

func (am *AuthenticationMiddleware) onAuthenticationSuccess(c echo.Context) {
	userToken := c.Get("userToken")
	claims := userToken.(*jwtutil.Token).Claims.(*jwtutil.UserClaims)
	c.Set("user", claims)
}

func (am *AuthenticationMiddleware) onAuthenticationError(c echo.Context, err error) error {
	if err != nil {
		c.Logger().Error(err)
	}
	am.authenticationManager.OnAuthenticationError(c.Request().Context(), c.RealIP())
	return echo.NewHTTPError(authentication.Unauthorized, authentication.UnauthorizedResponse)
}
