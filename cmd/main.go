package main

import (
	"distributed-id-generator/pkg/snowflake"
	"fmt"
)

func main() {
	snowflakeNode, err := snowflake.NewNode(0)
	if err != nil {
		panic(fmt.Errorf("error creating NewNode, %s", err))
	}
	fmt.Println(snowflakeNode.GenerateID())
}
