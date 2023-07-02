package users

import (
	"github.com/microservices-simulator-api/internal/utils/identifier"
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type NewUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(name, email, password string) User {
	return BuildUser(identifier.UniqueID(), name, email, password, time.Now())
}

func BuildUser(id int64, name, email, password string, createdAt time.Time) User {
	return User{
		Id:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	}
}
