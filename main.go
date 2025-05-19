package main

import (
	"fmt"
	httpn "http/http"
)

func healthCheck(req httpn.Request) httpn.Response {
	queryParams := req.ParseQuery()
	fmt.Println(queryParams)
	res := httpn.Response{
		StatusCode: 200,
		Data:       "{'message': 'success', 'status': 'ok'}",
		Status:     "OK",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return res
}

func main() {
	httpn.Get("/health-check", healthCheck)
	httpn.Listen(8003)
}
