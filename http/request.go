package httpn

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Request struct {
	method     string
	url        string
	endpoint   string
	Data       map[string]interface{}
	requestStr string
}

func (req Request) ParseQuery() map[string]string {
	var rawQuery []string = strings.Split(req.url, "?")
	if len(rawQuery) < 1 {
		return make(map[string]string)
	}
	result := make(map[string]string)
	query := strings.TrimPrefix(rawQuery[0], "?")
	if query == "" {
		return result
	}
	pairs := strings.Split(query, "&")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		key := kv[0]
		value := ""
		if len(kv) > 1 {
			value = kv[1]
		}
		result[key] = value
	}
	return result
}

func (request Request) resolveEndpoint() (func(Request) Response, bool) {
	for i := 0; i < len(Routes); {
		if Routes[i].Route == request.endpoint && Routes[i].Method == request.method {
			return Routes[i].Handler, true
		}
		i += 1
	}
	return nil, false
}

func getRequestPostData(message string) map[string]any {
	var lines []string = strings.Split(message, "\n")
	dataStartIndex := 0
	for i := 0; i < len(lines); {
		if len(lines[i]) == 1 {
			dataStartIndex = i
			break
		}
		i += 1
	}
	var rawDataStr string
	for i := dataStartIndex; i < len(lines); {
		rawDataStr += lines[i]
		i += 1
	}
	if len(rawDataStr) <= 0 {
		return make(map[string]any)
	}
	cleanStr := strings.ReplaceAll(rawDataStr, "\x00", "")
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(cleanStr), &jsonData)
	fmt.Println("RAW DATA ", jsonData)
	if err != nil {
		fmt.Println("ERROR CONVERT DATA TO JSON ", err)
	}
	return jsonData
}

func FormatRequest(message string) Request {
	var lines []string = strings.Split(message, "\n")
	var header string = lines[0]
	var method string = strings.Trim(strings.Split(header, " ")[0], "")
	var fullUrl string = strings.Trim(strings.Split(header, " ")[1], "")
	var endpoint string = strings.Split(fullUrl, "?")[0]

	var req Request = Request{
		requestStr: message,
		method:     method,
		url:        fullUrl,
		endpoint:   endpoint,
	}
	if method != "GET" {
		postData := getRequestPostData(message)
		req.Data = postData
	}
	return req
}
