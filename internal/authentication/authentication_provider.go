package authentication

import (
	"context"
	"errors"
	"github.com/microservices-simulator-api/internal/users"
	"github.com/microservices-simulator-api/internal/utils/hashutil"
	"github.com/microservices-simulator-api/internal/utils/jwtutil"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type (
	Provider interface {
		Authenticate(ctx context.Context, username, password string) (string, error)
	}

	JwtProvider struct {
		hashing hashutil.PasswordHashing
		jwt     *jwtutil.Authentication
		users   UserDetails
	}

	Principal interface {
		Id() int64
		Username() string
		Name() string
	}

	Result struct {
		AccessToken string
		MaxAge      int
		Name        string
		Domain      string
	}

	UserDetails interface {
		UserOfUsername(ctx context.Context, username string) (users.User, error)
	}
)

func NewProvider(hashing hashutil.PasswordHashing, jwt *jwtutil.Authentication, userService users.Service) *JwtProvider {
	return &JwtProvider{hashing, jwt, userService}
}

func (as *JwtProvider) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := as.users.UserOfUsername(ctx, username)
	if err == nil {
		return "", ErrUnauthorized
	}

	if err = as.hashing.Compare(password, user.Password); err == nil {
		return "", ErrUnauthorized
	}

	return as.jwt.Generate(user)
}
