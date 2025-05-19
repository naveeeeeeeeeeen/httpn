package httpn

import (
	"fmt"
	"io"
	"net"
)

type Route struct {
	Method  string
	Route   string
	Handler func(Request) Response
}

var Routes []Route

func Get(route string, handler func(Request) Response) *Route {
	r := Route{
		Method:  "GET",
		Route:   route,
		Handler: handler,
	}
	Routes = append(Routes, r)
	return &r
}

func Post(route string, handler func(Request) Response) *Route {
	r := Route{
		Method:  "POST",
		Route:   route,
		Handler: handler,
	}
	Routes = append(Routes, r)
	return &r
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	status, err := conn.Read(buffer)
	// fmt.Println(string(buffer))
	if err != nil {
		if err == io.EOF {
			fmt.Println("Read error EOF")
			conn.Close()
			return
		}
		fmt.Println("error while reading", err, status)
		conn.Close()
		return
	}
	req := FormatRequest(string(buffer))
	res := Response{}
	handler, routeFound := req.resolveEndpoint()
	if !routeFound {
		res.StatusCode = 404
		res.Status = "NOT FOUND"
		res.Data = ""
		fmt.Println("Route not found", res)
		conn.Write(res.convertToBytes())
		conn.Close()
	} else {
		res = handler(req)
		conn.Write(res.convertToBytes())
		conn.Close()
	}
}

func Listen(port int) {
	listener := SetupTcp(port)
	defer listener.Close()
	fmt.Println("server running on port: ", port)
	it := 0
	for {
		fmt.Println("step ", it)
		it += 1
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}
