package ghttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
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
	Client   http.Client                           //Client
	resp     *http.Response                        //Response
	Req      *http.Request                         //Request
	Proxy    func(*http.Request) (*url.URL, error) //Request proxy
	Link     string                                //Request address
	SendType string                                //Request type
	Header   map[string]string                     //Request header
	Body     map[string]interface{}                //Request body
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

//SetProxy Set proxy
func (h *HttpSend) SetProxy(proxy string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Proxy = func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
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

//Get request
func (h *HttpSend) Get(url string) ([]byte, error) {
	h.Link = url
	return h.send(GET)
}

//Post request
func (h *HttpSend) Post(url string) ([]byte, error) {
	h.Link = url
	return h.send(POST)
}

//Put request
func (h *HttpSend) Put(url string) ([]byte, error) {
	h.Link = url
	return h.send(PUT)
}

//Delete request
func (h *HttpSend) Delete(url string) ([]byte, error) {
	h.Link = url
	return h.send(DELETE)
}

//Connect request
func (h *HttpSend) Connect(url string) ([]byte, error) {
	h.Link = url
	return h.send(CONNECT)
}

//Head request
func (h *HttpSend) Head(url string) ([]byte, error) {
	h.Link = url
	return h.send(HEAD)
}

//Options request
func (h *HttpSend) Options(url string) ([]byte, error) {
	h.Link = url
	return h.send(OPTIONS)
}

//Trace request
func (h *HttpSend) Trace(url string) ([]byte, error) {
	h.Link = url
	return h.send(TRACE)
}

//Patch ...
func (h *HttpSend) Patch(url string) ([]byte, error) {
	h.Link = url
	return h.send(PATCH)
}

//GetUrlBuild ...
func GetUrlBuild(link string, data map[string]string) string {
	u, _ := url.Parse(link)
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
	h.Client.Transport = Transport
	h.Req, err = http.NewRequest(method, h.Link, sendData)
	if err != nil {
		return nil, err
	}
	defer h.Req.Body.Close()

	//设置默认header
	if len(h.Header) == 0 {
		//json
		if strings.ToLower(h.SendType) == SENDTYPE_JSON {
			h.Header = map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			}
		} else { //form
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

	h.resp, err = h.Client.Do(h.Req)
	if err != nil {
		return nil, err
	}

	defer h.resp.Body.Close()

	return ioutil.ReadAll(h.resp.Body)
}
