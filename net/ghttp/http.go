package ghttp

import (
	"bytes"
	"crypto/tls"
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
	req      *http.Request
	Link     string                 //Request address
	SendType string                 //Request type
	Header   map[string]string      //Request header
	Body     map[string]interface{} //Request body
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
	h.req.AddCookie(c)
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
	var (
		resp   *http.Response
		client http.Client
		err    error
	)
	configData, err := json.Marshal(h.Body)
	if err != nil {
		log.Println(err.Error())
	}
	var sendData = bytes.NewBuffer(configData)

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	h.req, err = http.NewRequest(method, h.Link, sendData)
	if err != nil {
		return nil, err
	}
	defer h.req.Body.Close()

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
			h.req.Host = v
		} else {
			h.req.Header.Add(k, v)
		}
	}

	resp, err = client.Do(h.req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
