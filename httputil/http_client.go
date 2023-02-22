package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	httpClient = http.DefaultClient

	bufferPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 4096))
		},
	}
)

// http请求客户端，实现链式调用，包含http请求和响应
// 默认使用http.DefaultClient，可以自定义Transport
// 通过New进行初始化，Do进行请求发送，ParseResponseBody进行响应体接受
type RestClient struct {
	// input
	url    string
	method string
	param  interface{}
	header map[string]string

	// output
	resCode   int
	resHeader http.Header
	resBytes  []byte
	
	restErr  error
}

// New http客户端初始化
// url http请求地址
// methed http请求类型 GET POST DELETE...
// param 请求参数
// header 自定义请求头
func (client *RestClient) New(url string, methed string, param interface{}, header map[string]string) *RestClient {
	return &RestClient{url: url, method: methed, param: param, header: header}
}

// Do 发送http请求返回响应
func (client *RestClient) Do() *RestClient {
	req:= client.assemblyRequest()
	if client.restErr != nil {
		return client
	}
	res, err := httpClient.Do(req)
	if err != nil {
		client.restErr = err
		return client
	}
	defer res.Body.Close()

	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)

	_, err = io.Copy(buffer, res.Body)
	if err != nil {
		client.restErr = err
		return client
	}

	client.resBytes = make([]byte, buffer.Len())
	copy(client.resBytes, buffer.Bytes())
	client.resCode = res.StatusCode
	client.resHeader = res.Header
	return client
}

// assemblyRequest 组装请求体
func (client *RestClient) assemblyRequest() (*http.Request) {
	var body io.Reader
	switch instance := client.param.(type) {
	case string:
		body = strings.NewReader(instance)
	default:
		paramBytes, err := json.Marshal(client.param)
		if err != nil {
			client.restErr = err
			return nil
		}
		body = bytes.NewBuffer(paramBytes)
	}

	req, err := http.NewRequest(client.method, client.url, body)
	if err != nil {
		client.restErr = err
		return nil
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	for k, v := range client.header {
		req.Header.Add(k, v)
	}
	return req

}

// ParseResponseBody 获取http返回体，根据传入结构体参数自动转换
// receiver 响应体接收者，目前只支持结构体指针类型
func (client *RestClient) ParseResponseBody(receiver interface{}) {
	if client.restErr != nil {
		return
	}
	if err := json.Unmarshal(client.resBytes, receiver); err != nil {
		client.restErr = err
	}
}
