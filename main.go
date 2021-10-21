package main

import (
	"HTTPgo/HTTPgo"
	"fmt"
)

func main() {
	if err := HTTPgo.Run(); err != nil {
		fmt.Printf("Error: %s", err)
	}
}
