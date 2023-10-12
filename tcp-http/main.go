package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Request struct {
	Method      string
	Host        string
	Path        string
	UserAgent   string
	ContentType string
	Body        string
}

type Response struct {
	Status      int
	ContentType string
	Body        string
}

var STATUS = map[int]string{
	200: "HTTP/1.1 200 OK",
	404: "HTTP/1.1 404 NOT FOUND",
}

func ParseRequest(conn net.Conn) Request {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err.Error())
	}
	reqString := string(buffer)
	r := Request{
		Method: strings.Split(reqString, " ")[0],
		Path:   strings.Split(reqString, " ")[1],
		Body:   strings.Split(reqString, "\r\n\r\n")[1],
	}
	for _, headers := range strings.Split(reqString, "\r\n") {
		if strings.HasPrefix(headers, "Host") {
			r.Host = strings.Split(headers, "Host: ")[1]
		} else if strings.HasPrefix(headers, "User-Agent") {
			r.UserAgent = strings.Split(headers, "User-Agent: ")[1]
		} else if strings.HasPrefix(headers, "Content-Type") {
			r.ContentType = strings.Split(headers, "Content-Type: ")[1]
		}
	}
	return r
}

func (res Response) Respond() string {
	resString := STATUS[res.Status] + "\r\n"
	if res.ContentType != "" {
		resString += "Content-Type: " + res.ContentType + "\r\n"
	}
	if res.Body != "" {
		resString += fmt.Sprintf("Content-Length: %d\r\n", len(res.Body))
		resString += "\r\n" + res.Body
	}

	return resString
}

func Handler(req Request) string {
	// respond a html content
	if req.Path == "/" {
		return Response{
			Status:      200,
			ContentType: "text/html",
			Body:        `<a href="https://github.com/chamanbravo">CallbackCAT</a>`,
		}.Respond()
	}

	// respond a json
	if req.Path == "/json" {
		return Response{
			Status:      200,
			ContentType: "application/json",
			Body:        `{"name": "callbackCAT"}`,
		}.Respond()
	}

	return Response{
		Status: 404,
	}.Respond()
}

func HandleConn(conn net.Conn) {
	req := ParseRequest(conn)
	res := Handler(req)
	conn.Write([]byte(res))
	conn.Close()
}

func main() {
	li, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go HandleConn(conn)
	}
}
