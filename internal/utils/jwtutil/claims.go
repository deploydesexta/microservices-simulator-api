package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidUserId = errors.New("userId must be a valid number")
)

type UserClaims struct {
	jwt.RegisteredClaims
	Id   int64  `json:"user_id"`
	Name string `json:"user_name"`
}

func NewUserClaims(userId int64, name, domain, jit string, maxAge int) *UserClaims {
	return &UserClaims{
		Id:   userId,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{domain},
			Issuer:    domain,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(maxAge))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   name,
			ID:        jit,
		},
	}
}

func (uc UserClaims) Valid() error {
	if uc.Id < 0 {
		return ErrInvalidUserId
	}

	return nil
}
