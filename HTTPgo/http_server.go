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

	socket, err := listener.Accept()
	if err != nil {
		return err
	}
	defer socket.Close()

	header := make(map[string]string)
	reader := bufio.NewReader(socket)
	scanner := textproto.NewReader(reader)

	err = readHttpRequestHeader(scanner, header)
	if err != nil {
		return err
	}

	if header["Method"] == "GET" {
		path, ok := header["Path"]
		if !ok {
			return errors.New("no path found")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		resourcePath := filepath.Join(cwd, filepath.Clean(path))
		resource, err := ioutil.ReadFile(resourcePath)
		if err != nil {
			return err
		}
		err = writeHttpResponseWithResource(&socket, resource)
		if err != nil {
			return err
		}
	} else {
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
				err = writeHttpResponse(&socket)
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

			err = writeHttpResponse(&socket)
			if err != nil {
				return err
			}
		}

	}

	fmt.Println("Server: close listen...")
	return nil
}
