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

type HttpSend struct {
	Link     string
	SendType string
	Header   map[string]string
	Body     map[string]interface{}
	sync.RWMutex
}

func Client() *HttpSend {
	return &HttpSend{
		SendType: SENDTYPE_JSON,
	}
}

func (h *HttpSend) SetBody(body map[string]interface{}) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Body = body
	return h
}

func (h *HttpSend) SetHeader(header map[string]string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Header = header
	return h
}

func (h *HttpSend) SetSendType(send_type string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.SendType = send_type
	return h
}

func (h *HttpSend) Get(url string) ([]byte, error) {
	h.Link = url
	return h.send(GET)
}

func (h *HttpSend) Post(url string) ([]byte, error) {
	h.Link = url
	return h.send(POST)
}

func (h *HttpSend) Put(url string) ([]byte, error) {
	h.Link = url
	return h.send(PUT)
}

func (h *HttpSend) Delete(url string) ([]byte, error) {
	h.Link = url
	return h.send(DELETE)
}

func (h *HttpSend) Connect(url string) ([]byte, error) {
	h.Link = url
	return h.send(CONNECT)
}

func (h *HttpSend) Head(url string) ([]byte, error) {
	h.Link = url
	return h.send(HEAD)
}

func (h *HttpSend) Options(url string) ([]byte, error) {
	h.Link = url
	return h.send(OPTIONS)
}

func (h *HttpSend) Trace(url string) ([]byte, error) {
	h.Link = url
	return h.send(TRACE)
}

func (h *HttpSend) Patch(url string) ([]byte, error) {
	h.Link = url
	return h.send(PATCH)
}

func GetUrlBuild(link string, data map[string]string) string {
	u, _ := url.Parse(link)
	q := u.Query()
	for k, v := range data {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (h *HttpSend) send(method string) ([]byte, error) {
	var (
		req    *http.Request
		resp   *http.Response
		client http.Client
		err    error
	)

	configdata, err := json.Marshal(h.Body)
	if err != nil {
		log.Println(err.Error())
	}
	var send_data = bytes.NewBuffer(configdata)

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	log.Println(method)
	log.Println(h.Link)
	log.Println(send_data)
	req, err = http.NewRequest(method, h.Link, send_data)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

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
			req.Host = v
		} else {
			req.Header.Add(k, v)
		}
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
