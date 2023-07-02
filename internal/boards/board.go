package boards

import (
	"github.com/microservices-simulator-api/internal/utils/identifier"
	"strings"
	"time"
)

type Board struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerId   int64     `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBoard(ownerId int64) Board {
	return Board{
		Id:        uniqueId(),
		Name:      "Untitled",
		OwnerId:   ownerId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func uniqueId() string {
	ulid := identifier.UniqueULID()
	return strings.ToLower(ulid[0:4] + "-" + ulid[4:8] + "-" + ulid[8:12])
}
