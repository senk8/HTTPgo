package HTTPgo

import (
	"fmt"
	"net/textproto"
	"strings"
)

func ReadHttpRequestHeader(scanner *textproto.Reader, header map[string]string) error {
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
			header["Method"] = requestLine[0]
			header["Path"] = requestLine[1]
			fmt.Println(method, path)
			continue
		}

		headerFields := strings.SplitN(line, ": ", 2)
		fmt.Printf("%s: %s\n", headerFields[0], headerFields[1])
		header[headerFields[0]] = headerFields[1]
	}
	return nil
}
