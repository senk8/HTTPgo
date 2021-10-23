package HTTPgo

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
)

func Run() error {
	fmt.Println("Server: start listen...")

	listener, err := net.Listen("tcp", "127.0.0.1:50000")
	if err != nil {
		return err
	}
	defer listener.Close()

	socket, err := listener.Accept()
	if err != nil {
		return err
	}

	go func(socket net.Conn) {
		defer socket.Close()
		if err := service(socket); err != nil {
			fmt.Printf("%+v", err)
		}
	}(socket)

	fmt.Println("Server: close listen...")
	return nil
}

func service(socket net.Conn) error {
	reader := bufio.NewReader(socket)
	scanner := textproto.NewReader(reader)

	request := HttpRequest{
		Header: make(map[string]string),
	}

	err := request.readHttpRequestHeader(scanner)
	if err != nil {
		return err
	}

	if request.IsGet() {
		err := responseGetRequest(&socket, &request)
		if err != nil {
			return err
		}
	} else if request.IsPost() {
		err := processPostRequest(reader, scanner, &request)
		if err != nil {
			return err
		}
		err = writeHttpResponse(&socket)
		if err != nil {
			return err
		}
	} else {
		panic("un-implement methods")
	}

	return nil
}

func responseGetRequest(socket *net.Conn, request *HttpRequest) error {
	path, ok := request.Header["Path"]
	if !ok {
		return errors.New("no path found")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	resourcePath := filepath.Join(cwd, filepath.Clean(path))
	if !fileExists(resourcePath) {
		err := writeHttpResponseNotFound(socket)
		if err != nil {
			return err
		}
	} else {
		resource, err := ioutil.ReadFile(resourcePath)
		if err != nil {
			return err
		}
		err = writeHttpResponseWithResource(socket, resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func processPostRequest(reader *bufio.Reader, scanner *textproto.Reader, request *HttpRequest) error {
	transferEncoding, ok := request.Header["Transfer-Encoding"]
	if ok {
		if transferEncoding == "chunked" {
			for {
				line, err := scanner.ReadLine()
				if line == "0" {
					break
				}
				if err != nil {
					return err
				}
				fmt.Println(line)
			}
		} else {
			return errors.New("Transfer-Encoding type is not defined.")
		}
	} else {
		err := request.readBodyWithContentLength(reader)
		if err != nil {
			return err
		}
		fmt.Printf("Body:%s\n", string(request.Body))
	}

	return nil
}
