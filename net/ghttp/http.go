package ghttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

//GET POST PUT DELETE HEAD PATCH CONNECT OPTIONS TRACE SENDTYPE_JSON Requet Type
var (
	GET           = "GET"
	POST          = "POST"
	PUT           = "PUT"
	DELETE        = "DELETE"
	HEAD          = "HEAD"
	PATCH         = "PATCH"
	CONNECT       = "CONNECT"
	OPTIONS       = "OPTIONS"
	TRACE         = "TRACE"
	SENDTYPE_JSON = "json"
)

//HttpSend Request parameter
type HttpSend struct {
	Client   http.Client                                  //Client
	Response *http.Response                               //Response
	Req      *http.Request                                //Request
	Proxy    func(*http.Request) (*url.URL, error)        //Request proxy
	Url      string                                       //Request address
	SendType string                                       //Request type
	Header   map[string]string                            //Request header
	Body     map[string]interface{}                       //Request body
	Dial     func(network, addr string) (net.Conn, error) //Request Timeout
	sync.RWMutex
}

//Client Default request
func Client() *HttpSend {
	return &HttpSend{
		SendType: SENDTYPE_JSON,
	}
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

//SetSendType Set Type
func (h *HttpSend) SetSendType(sendType string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.SendType = sendType
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
	return h.send(GET)
}

//Post request
func (h *HttpSend) Post(url string) ([]byte, error) {
	h.Url = url
	return h.send(POST)
}

//Put request
func (h *HttpSend) Put(url string) ([]byte, error) {
	h.Url = url
	return h.send(PUT)
}

//Delete request
func (h *HttpSend) Delete(url string) ([]byte, error) {
	h.Url = url
	return h.send(DELETE)
}

//Connect request
func (h *HttpSend) Connect(url string) ([]byte, error) {
	h.Url = url
	return h.send(CONNECT)
}

//Head request
func (h *HttpSend) Head(url string) ([]byte, error) {
	h.Url = url
	return h.send(HEAD)
}

//Options request
func (h *HttpSend) Options(url string) ([]byte, error) {
	h.Url = url
	return h.send(OPTIONS)
}

//Trace request
func (h *HttpSend) Trace(url string) ([]byte, error) {
	h.Url = url
	return h.send(TRACE)
}

//Patch ...
func (h *HttpSend) Patch(url string) ([]byte, error) {
	h.Url = url
	return h.send(PATCH)
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

//send ...
func (h *HttpSend) send(method string) ([]byte, error) {
	configData, err := json.Marshal(h.Body)
	if err != nil {
		log.Println(err.Error())
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

	h.Req, err = http.NewRequest(method, h.Url, sendData)
	if err != nil {
		return nil, err
	}
	defer h.Req.Body.Close()

	if len(h.Header) == 0 {
		if strings.ToLower(h.SendType) == SENDTYPE_JSON {
			h.Header = map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			}
		} else {
			h.Header = map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			}
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
	defer h.Response.Body.Close()
	return ioutil.ReadAll(h.Response.Body)
}

//Send<扩展一般用于手动请求>
func (h *HttpSend) Send(method string) (io.ReadCloser, error) {
	configData, err := json.Marshal(h.Body)
	if err != nil {

		log.Println(err.Error())
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

	h.Req, err = http.NewRequest(method, h.Url, sendData)
	if err != nil {
		return nil, err
	}

	if len(h.Header) == 0 {
		if strings.ToLower(h.SendType) == SENDTYPE_JSON {
			h.Header = map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			}
		} else {
			h.Header = map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			}
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
