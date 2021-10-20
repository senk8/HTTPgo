package main

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"

	"net"
	"strconv"
	"strings"
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

	reader := bufio.NewReader(socket)
	scanner := textproto.NewReader(reader)

	var contentLength int
	for {
		line, err := scanner.ReadLine()
		if line == "" {
			break
		}
		if err != nil {
			return err
		}
		if strings.HasPrefix(line, "Content-Length") {
			contentLength, err = strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[1]))
			if err != nil {
				return err
			}
		}
		fmt.Println(line)
	}

	buf := make([]byte, contentLength)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return err
	}
	fmt.Printf("Body:%s/n", string(buf))

	/*
		buf := make([]byte,1024)

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
	*/

	fmt.Println("Server: close listen...")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err)
	}
}
