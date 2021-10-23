package HTTPgo

import (
	"io"
	"net"
)

func writeHttpResponse(socket *net.Conn) error {
	_, err := io.WriteString(*socket, "HTTP/1.1 200 OK\r\n")
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
	_, err = io.WriteString(*socket, "<h1>Hello World!!</h1>")
	if err != nil {
		return err
	}
	return nil
}

func writeHttpResponseWithResource(socket *net.Conn, resource []byte) error {
	_, err := io.WriteString(*socket, "HTTP/1.1 200 OK\r\n")
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
	_, err = io.WriteString(*socket, string(resource))
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
