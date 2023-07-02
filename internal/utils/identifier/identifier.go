package identifier

import (
	"github.com/bwmarrin/snowflake"
	"github.com/oklog/ulid/v2"
)

type ID int64

var (
	node, _ = snowflake.NewNode(1)
)

func UniqueULID() string {
	return ulid.Make().String()
}

func UniqueID() int64 {
	return node.Generate().Int64()
}
