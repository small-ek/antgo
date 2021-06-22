package ahttp

import (
	"bytes"
	"encoding/json"
	"github.com/small-ek/antgo/os/logs"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

//HttpSend Request parameter
type HttpSend struct {
	Client      http.Client                                  //Client
	Response    *http.Response                               //Response
	Req         *http.Request                                //Request
	Proxy       func(*http.Request) (*url.URL, error)        //Request proxy
	Url         string                                       //Request address
	ContentType string                                       //Request type
	Header      map[string]string                            //Request header
	Body        map[string]interface{}                       //Request body
	Dial        func(network, addr string) (net.Conn, error) //Request Timeout
	Method      string                                       //Request method
	FormData    url.Values
	sync.RWMutex
}

//Client Default request
func Client() *HttpSend {
	return &HttpSend{
		ContentType: "application/json",
	}
}

//SetPostFormData Set Post Form Data
func (h *HttpSend) SetPostFormData(FormData url.Values) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.FormData = FormData
	return h
}

//SetBody Set body
func (h *HttpSend) SetBody(body map[string]interface{}) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Body = body
	return h
}

//SetUrl Set url
func (h *HttpSend) SetUrl(url string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Url = url
	return h
}

//SetProxy Set proxy
func (h *HttpSend) SetProxy(proxy string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Proxy = func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}
	return h
}

//SetTimeout Set Timeout
func (h *HttpSend) SetTimeout(timeout time.Duration) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Dial = func(netw, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(netw, addr, time.Second*timeout) //设置建立连接超时
		if err != nil {
			return nil, err
		}
		c.SetDeadline(time.Now().Add(timeout * time.Second)) //设置发送接收数据超时
		return c, nil
	}
	return h
}

//SetHeader set header
func (h *HttpSend) SetHeader(header map[string]string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Header = header
	return h
}

//SetCookie set cookie
func (h *HttpSend) SetCookie(c *http.Cookie) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Req.AddCookie(c)
	return h
}

//SetMethod set method
func (h *HttpSend) SetMethod(method string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Method = method
	return h
}

//SetSendType Set Type
func (h *HttpSend) SetContentType(ContentType string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.ContentType = ContentType
	return h
}

//GetHeader Get Response Header
func (h *HttpSend) GetHeader() map[string][]string {
	h.Lock()
	defer h.Unlock()
	if h.Response != nil {
		return h.Response.Header
	}
	return nil
}

//Get request
func (h *HttpSend) Get(url string) ([]byte, error) {
	h.Url = url
	h.Method = "GET"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//PostForm request
func (h *HttpSend) PostForm(urls string) ([]byte, error) {
	h.Url = urls
	h.Method = "POST"
	resp, err := http.PostForm(urls, h.FormData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	return body, err2
}

//Post request
func (h *HttpSend) Post(url string) ([]byte, error) {
	h.Url = url
	h.Method = "POST"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//Put request
func (h *HttpSend) Put(url string) ([]byte, error) {
	h.Url = url
	h.Method = "PUT"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//Delete request
func (h *HttpSend) Delete(url string) ([]byte, error) {
	h.Url = url
	h.Method = "DELETE"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//Connect request
func (h *HttpSend) Connect(url string) ([]byte, error) {
	h.Url = url
	h.Method = "CONNECT"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//Head request
func (h *HttpSend) Head(url string) ([]byte, error) {
	h.Url = url
	h.Method = "HEAD"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//Options request
func (h *HttpSend) Options(url string) ([]byte, error) {
	h.Url = url
	h.Method = "OPTIONS"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//Trace request
func (h *HttpSend) Trace(url string) ([]byte, error) {
	h.Url = url
	h.Method = "TRACE"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//Patch ...
func (h *HttpSend) Patch(url string) ([]byte, error) {
	h.Url = url
	h.Method = "PATCH"
	var result, err = h.Send()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

//GetUrlBuild ...
func GetUrlBuild(urls string, data map[string]string) string {
	u, _ := url.Parse(urls)
	q := u.Query()
	for k, v := range data {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

//Send <扩展一般用于手动请求>
func (h *HttpSend) Send() (io.ReadCloser, error) {
	configData, err := json.Marshal(h.Body)
	if err != nil {
		logs.Error(err.Error())
	}
	var sendData = bytes.NewBuffer(configData)

	var Transport = &http.Transport{}
	if h.Proxy != nil {
		Transport.Proxy = h.Proxy
	}

	if h.Dial != nil {
		Transport.Dial = h.Dial
	}

	h.Client.Transport = Transport
	if h.Method == "" {
		h.Method = "GET"
	}
	h.Req, err = http.NewRequest(h.Method, h.Url, sendData)
	if err != nil {
		return nil, err
	}

	if len(h.Header) == 0 {
		h.Header = map[string]string{
			"Content-Type": h.ContentType,
		}
	}

	for k, v := range h.Header {
		if strings.ToLower(k) == "host" {
			h.Req.Host = v
		} else {
			h.Req.Header.Add(k, v)
		}
	}

	h.Response, err = h.Client.Do(h.Req)
	if err != nil {
		return nil, err
	}
	return h.Response.Body, nil
}

//Close <必须默认关闭>
func (h *HttpSend) Close() {
	defer h.Response.Body.Close()
	defer h.Req.Body.Close()
}
