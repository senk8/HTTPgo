package main

import (
	"fmt"
	"net"
)

func run() error {
	fmt.Println("Server: start listen...")

	listener, err := net.Listen("tcp", "127.0.0.1:50000")
	if err != nil {
		return err
	}

	socket, err := listener.Accept()
	if err != nil {
		return err
	}
	defer socket.Close()

	buf := make([]byte, 1024)

	for {
		n, err := socket.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(string(buf[:n]))
	}

	fmt.Println("Server: close listen...")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err)
	}
}
