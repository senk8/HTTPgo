package main

import (
	"HTTPgo/HTTPgo"
	"fmt"
)

func main() {
	for {
		if err := HTTPgo.Run(); err != nil {
			fmt.Printf("Error: %s", err)
		}
	}
}
