package shapes

import "github.com/microservices-simulator-api/internal/utils/identifier"

type Shape struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Type  string `json:"type"`
	Group string `json:"group"`
}

type ShapeInput struct {
	Name  string `json:"name" form:"name"`
	Type  string `json:"type" form:"type"`
	Group string `json:"group" form:"group"`
}

func NewShapeFromInput(si ShapeInput, filepath string) Shape {
	return Shape{
		Id:    identifier.UniqueID(),
		Name:  si.Name,
		Image: filepath,
		Type:  si.Type,
		Group: si.Group,
	}
}
