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

	var method, path string
	header := make(map[string]string)

	isRequestLine := true
	for {
		line, err := scanner.ReadLine()
		if line == "" {
			break
		}

		if err != nil {
			return err
		}

		if isRequestLine {
			isRequestLine = false
			requestLine := strings.Fields(line)
			header["Method"] = requestLine[0]
			header["Path"] = requestLine[1]
			fmt.Println(method, path)
			fmt.Println("foo")
			continue
		}

		headerFields := strings.SplitN(line, ": ", 2)
		fmt.Printf("%s: %s\n", headerFields[0], headerFields[1])
		header[headerFields[0]] = headerFields[1]
	}

	contentLengthStr, ok := header["Content-Length"]
	if !ok {
		return fmt.Errorf("%d", 12)
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return err
	}

	buf := make([]byte, contentLength)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return err
	}
	fmt.Printf("Body:%s\n", string(buf))
	fmt.Println("Server: close listen...")
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err)
	}
}
