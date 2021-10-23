package HTTPgo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Header map[string]string
	Body   []byte
}

func (req *HttpRequest) readHttpRequestHeader(scanner *textproto.Reader) error {
	var method, path string
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
			req.Header["Method"] = requestLine[0]
			req.Header["Path"] = requestLine[1]
			fmt.Println(method, path)
			continue
		}

		headerFields := strings.SplitN(line, ": ", 2)
		fmt.Printf("%s: %s\n", headerFields[0], headerFields[1])
		req.Header[headerFields[0]] = headerFields[1]
	}
	return nil
}

func (req *HttpRequest) readBodyWithContentLength(reader *bufio.Reader) error {
	contentLength, err := req.ContentLength()
	if err != nil {
		return err
	}

	req.Body = make([]byte, contentLength)
	_, err = io.ReadFull(reader, req.Body)
	if err != nil {
		return err
	}

	return nil
}

func (req *HttpRequest) IsGet() bool {
	return req.Header["Method"] == "GET"
}

func (req *HttpRequest) IsPost() bool {
	return req.Header["Method"] == "POST"
}

func (req *HttpRequest) ContentLength() (int, error) {
	contentLengthStr, ok := req.Header["Content-Length"]
	if !ok {
		return 0, errors.New("Content-Length must be specified. ")
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return 0, err
	}

	return contentLength, nil
}
