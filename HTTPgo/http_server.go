package HTTPgo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
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
	header := make(map[string]string)
	reader := bufio.NewReader(socket)
	scanner := textproto.NewReader(reader)

	err := readHttpRequestHeader(scanner, header)
	if err != nil {
		return err
	}

	if header["Method"] == "GET" {
		err := processGetRequest(&socket, header)
		if err != nil {
			return err
		}
	} else if header["Method"] == "POST" {
		err := processPostRequest(&socket, reader, scanner, header)
		if err != nil {
			return err
		}
	} else {
		panic("un-implement methods")
	}

	return nil
}

func processGetRequest(socket *net.Conn, header map[string]string) error {
	path, ok := header["Path"]
	if !ok {
		return errors.New("no path found")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	resourcePath := filepath.Join(cwd, filepath.Clean(path))
	if !fileExists(resourcePath) {
		io.WriteString(*socket, "HTTP/1.1 404 Not Found\r\n")
		if err != nil {
			return err
		}
		io.WriteString(*socket, "Content-Type: text/html\r\n")
		if err != nil {
			return err
		}
		io.WriteString(*socket, "\r\n")
		if err != nil {
			return err
		}
		io.WriteString(*socket, string("<h1>Error 404</h1>"))
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

func processPostRequest(socket *net.Conn, reader *bufio.Reader, scanner *textproto.Reader, header map[string]string) error {
	transferEncoding, ok := header["Transfer-Encoding"]
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
			err := writeHttpResponse(socket)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Transfer-Encoding type is not defined.")
		}
	} else {
		contentLengthStr, ok := header["Content-Length"]
		if !ok {
			return errors.New("Content-Length must be specified. ")
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

		err = writeHttpResponse(socket)
		if err != nil {
			return err
		}
	}
	return nil
}
