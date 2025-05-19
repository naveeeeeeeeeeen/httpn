package httpn

import "fmt"

type Response struct {
	StatusCode int
	Status     string
	Data       string
	Headers    map[string]string
}

func (res Response) convertToBytes() []byte {
	rawString := ""
	rawString += fmt.Sprintf("HTTP/1.1 %d %s \r\n", res.StatusCode, res.Status)
	for key, value := range res.Headers {
		rawString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	rawString += fmt.Sprintf("\r\n%s", res.Data)
	return []byte(rawString)
}
