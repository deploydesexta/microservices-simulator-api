package users

import (
	"context"
	"github.com/microservices-simulator-api/internal/utils/hashutil"
)

type (
	Service interface {
		Create(ctx context.Context, input NewUserInput) (User, error)
		UserOfId(ctx context.Context, id int64) (User, error)
		UserOfUsername(ctx context.Context, username string) (User, error)
	}

	UserService struct {
		hashing    hashutil.PasswordHashing
		repository Repository
	}
)

func NewService(hashing hashutil.PasswordHashing, repository Repository) *UserService {
	return &UserService{hashing, repository}
}

func (us *UserService) Create(ctx context.Context, input NewUserInput) (User, error) {
	pwd, err := us.hashing.Hash(input.Password)
	if err != nil {
		return User{}, err
	}

	return us.repository.Create(ctx, NewUser(input.Name, input.Email, pwd))
}

func (us *UserService) UserOfId(ctx context.Context, id int64) (User, error) {
	return us.repository.UserOfId(ctx, id)
}

func (us *UserService) UserOfUsername(ctx context.Context, username string) (User, error) {
	return us.repository.UserOfUsername(ctx, username)
}
