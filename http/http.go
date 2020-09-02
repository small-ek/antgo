package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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
	SENDTYPE_JSON = "json"
)

type HttpSend struct {
	Link     string
	SendType string
	Header   map[string]string
	Body     map[string]string
	sync.RWMutex
}

func NewClient(link string) *HttpSend {
	return &HttpSend{
		Link:     link,
		SendType: SENDTYPE_JSON,
	}
}

func (h *HttpSend) SetBody(body map[string]string) {
	h.Lock()
	defer h.Unlock()
	h.Body = body
}

func (h *HttpSend) SetHeader(header map[string]string) {
	h.Lock()
	defer h.Unlock()
	h.Header = header
}

func (h *HttpSend) SetSendType(send_type string) {
	h.Lock()
	defer h.Unlock()
	h.SendType = send_type
}

func (h *HttpSend) Get() ([]byte, error) {
	return h.send(GET)
}

func (h *HttpSend) Post() ([]byte, error) {
	return h.send(POST)
}

func (h *HttpSend) Put() ([]byte, error) {
	return h.send(PUT)
}

func (h *HttpSend) DELETE() ([]byte, error) {
	return h.send(DELETE)
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
		req       *http.Request
		resp      *http.Response
		client    http.Client
		send_data string
		err       error
	)

	if len(h.Body) > 0 {
		if strings.ToLower(h.SendType) == SENDTYPE_JSON {
			send_body, json_err := json.Marshal(h.Body)
			if json_err != nil {
				return nil, json_err
			}
			send_data = string(send_body)
		} else {
			send_body := http.Request{}
			err = send_body.ParseForm()
			if err != nil {
				log.Println(err.Error())
			}
			for k, v := range h.Body {
				send_body.Form.Add(k, v)
			}
			send_data = send_body.Form.Encode()
		}
	}

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err = http.NewRequest(method, h.Link, strings.NewReader(send_data))
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error http code :%d", resp.StatusCode))
	}

	return ioutil.ReadAll(resp.Body)
}
