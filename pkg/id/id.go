package id

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
}

func GenerateString() string {
	return node.Generate().String()
}
func GenerateInt64() int64 {
	return node.Generate().Int64()
}
