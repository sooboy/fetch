package fetch

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 放弃安全验证
		Dial: func(netw, addr string) (net.Conn, error) {
			client, err := net.DialTimeout(netw, addr, time.Second*3)
			if err != nil {
				return nil, err
			}
			client.SetDeadline(time.Now().Add(time.Second * 5)) // 5 秒等待建立链接
			return client, nil
		},
		ResponseHeaderTimeout: time.Second * 3, // 3秒等待数据响应
	},
}

// RequestData 请求信息
type RequestData struct {
	Header http.Header
	Query  url.Values
	Form   url.Values
	Body   io.Reader
}

// request 发送一条请求
func request(method int, url string, model interface{}, requestData *RequestData) error {
	var input io.Reader
	var header http.Header
	var methodName = GET_TAG

	if requestData.Query != nil {
		url += "?" + requestData.Query.Encode()
	}
	if requestData.Header != nil {
		header = requestData.Header
	}
	switch method {
	case POST:
		methodName = POST_TAG
		if requestData.Body != nil {
			input = requestData.Body
		}
	}
	req, err := http.NewRequest(methodName, url, input)
	if err != nil {
		return err
	}

	for k, v := range header {
		for i := 0; i < len(v); i++ {
			req.Header.Set(k, v[i])
		}
	}
	// 针对表单写的content-type  json的content-type暂时没使用到
	if method == POST {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(model)
}

func get(url string, header map[string]string, query map[string]string, model interface{}) error {

	requestData := &RequestData{
		Header: parseHeader(header),
		Query:  parseQuery(query),
	}

	return request(GET, url, model, requestData)
}

func post(url string, header map[string]string, param map[string]string, model interface{}) error {

	requestData := &RequestData{
		Header: parseHeader(header),
		Body:   strings.NewReader(parseQuery(param).Encode()),
	}

	// log.Println(parseQuery(param).Encode())
	return request(POST, url, model, requestData)
}
func post_json(url string, header map[string]string, param map[string]string, model interface{}) error {
	jsonStr, err := parseJSONstr(param)
	if err != nil {
		return err
	}
	requestData := &RequestData{
		Header: parseHeader(header),
		Body:   strings.NewReader(jsonStr),
	}

	return request(POST, url, model, requestData)
}

func parseHeader(header map[string]string) http.Header {
	if header == nil {
		return nil
	}
	h := make(http.Header)
	for key, value := range header {
		h.Add(key, value)
	}
	return h
}

func parseQuery(query map[string]string) url.Values {
	if query == nil {
		return nil
	}
	q := make(url.Values)
	for key, value := range query {
		q.Add(key, value)
	}
	return q
}

// parseJSONstr 将map转为string 供post方法使用
func parseJSONstr(query map[string]string) (string, error) {
	bytes, err := json.Marshal(query)
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}
