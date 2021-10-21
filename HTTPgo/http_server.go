package HTTPgo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
)

func Run() error {
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

	header := make(map[string]string)
	reader := bufio.NewReader(socket)

	err = ReadHttpRequestHeader(reader, header)
	if err != nil {
		return err
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
