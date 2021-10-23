package HTTPgo

import (
	"fmt"
	"io"
	"net"
	"strings"
)

type HttpResponse struct {
	Scheme       string
	StatusCode   string
	StatusPhrase string
	ContentType  string
	Body         string
}

func (resp *HttpResponse) writeHttpResponse(socket *net.Conn) error {
	statusLine := strings.Join([]string{resp.Scheme, resp.StatusCode, resp.StatusPhrase}, " ") + "\r\n"
	_, err := io.WriteString(*socket, statusLine)
	if err != nil {
		return err
	}
	contentType := fmt.Sprintf("Content-Type: %s\r\n", resp.ContentType)
	_, err = io.WriteString(*socket, contentType)
	if err != nil {
		return err
	}
	_, err = io.WriteString(*socket, "\r\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(*socket, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func writeHttpResponseNotFound(socket *net.Conn) error {
	_, err := io.WriteString(*socket, "HTTP/1.1 404 Not Found\r\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(*socket, "Content-Type: text/html\r\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(*socket, "\r\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(*socket, string("<h1>Error 404</h1>"))
	if err != nil {
		return err
	}
	return nil
}
