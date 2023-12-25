package ahttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/small-ek/antgo/utils/conv"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// HttpSend Request parameter
type HttpSend struct {
	Client      *http.Client                                 //Client
	Response    *http.Response                               //Response
	Req         *http.Request                                //Request
	Proxy       func(*http.Request) (*url.URL, error)        //Request proxy<代理地址>
	Url         string                                       //Request address<请求地址>
	ContentType string                                       //Request type<网络文件的类型>
	Header      map[string]string                            //Request header<请求头>
	Cookies     map[string]string                            //Request Cookies<请求Cookies>
	Timeout     time.Duration                                //Request timeout<请求超时时间>
	Body        map[string]interface{}                       //Request body<请求体>
	Dial        func(network, addr string) (net.Conn, error) //Request Timeout<请求超时>
	Method      string                                       //Request method<请求类型>
	Files       []string                                     //Request Files <多个文件>
	File        string                                       //Request File <单个文件>
	FileKey     string                                       //Request FileKey<文件Key>
	FileName    string                                       //Request FileName<文件名称>
	BodyReader  io.Reader                                    //Request BodyReader<读取器>
	debug       bool
	sync.RWMutex
}

var singletonHttpSend *HttpSend
var once sync.Once

// Client Default request
func Client() *HttpSend {
	once.Do(func() {
		singletonHttpSend = &HttpSend{
			ContentType: "application/json",
			Client: &http.Client{
				Timeout: 30 * time.Second,
				Transport: &http.Transport{
					MaxIdleConns:        10000,
					MaxIdleConnsPerHost: 0,
					MaxConnsPerHost:     0,
					TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
					DialContext: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
					}).DialContext,
					ForceAttemptHTTP2:     true,
					IdleConnTimeout:       300 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
				},
			},
		}
	})

	return singletonHttpSend
}

// SetTransport <设置>
func (h *HttpSend) SetTransport(transport *http.Transport) *HttpSend {
	h.Client.Transport = transport
	return h
}

// SetBody Set body<设置请求体>
func (h *HttpSend) SetBody(body map[string]interface{}) *HttpSend {
	h.Lock()
	defer h.Unlock()
	configData, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	h.BodyReader = bytes.NewBuffer(configData)
	h.Body = body
	return h
}

// SetBodyReader Set body<设置读取器>
func (h *HttpSend) SetBodyReader(BodyReader io.Reader) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.BodyReader = BodyReader
	return h
}

// SetUrl Set url<设置请求地址>
func (h *HttpSend) SetUrl(url string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Url = url
	return h
}

// SetProxy Set proxy<>设置代理
func (h *HttpSend) SetProxy(proxy string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Proxy = func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}
	return h
}

// SetTimeout Set Timeout<设置超时>
func (h *HttpSend) SetTimeout(timeout time.Duration) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Timeout = timeout * time.Second
	h.Dial = func(netw, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(netw, addr, time.Second*timeout) //设置建立连接超时
		if err != nil {
			return nil, err
		}
		err2 := c.SetDeadline(time.Now().Add(timeout * time.Second))
		if err2 != nil {
			return nil, err2
		} //设置发送接收数据超时
		return c, nil
	}
	return h
}

// SetHeader set header<设置请求头>
func (h *HttpSend) SetHeader(header map[string]string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Header = header
	return h
}

// SetCookie set cookie<设置cookie>
func (h *HttpSend) SetCookie(cookies map[string]string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Cookies = cookies
	for k, v := range cookies {
		h.Req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return h
}

// SetMethod set method<设置请求类型>
func (h *HttpSend) SetMethod(method string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Method = method
	return h
}

// SetSendType Set Type<设置资源的MIME类型>
func (h *HttpSend) SetContentType(ContentType string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.ContentType = ContentType
	return h
}

// GetHeader Get Response Header<获取请求头>
func (h *HttpSend) GetHeader() map[string][]string {
	h.Lock()
	defer h.Unlock()
	if h.Response != nil {
		return h.Response.Header
	}
	return nil
}

// Get request<GET 请求>
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

// SetFiles Set Files<设置多个文件路径上传>
func (h *HttpSend) SetFiles(files []string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.Files = files
	return h
}

// SetFile Set File<设置单个文件上传路径>
func (h *HttpSend) SetFile(file string, key ...string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	if len(key) > 0 {
		h.FileKey = key[0]
	}
	h.File = file
	return h
}

// SetFileKey Set File<设置文件的Key>
func (h *HttpSend) SetFileKey(fileKey string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.FileKey = fileKey
	return h
}

// SetFileName Set File<设置文件的名称>
func (h *HttpSend) SetFileName(fileName string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.FileName = fileName
	return h
}

// SetFileKeyAndName Set File<设置文件Key和名称>
func (h *HttpSend) SetFileKeyAndName(fileKey, fileName string) *HttpSend {
	h.Lock()
	defer h.Unlock()
	h.FileKey = fileKey
	h.FileName = fileName
	return h
}

// PostForm request<Post 表单提交>
func (h *HttpSend) PostForm(url string) ([]byte, error) {
	h.Url = url
	h.Method = "POST"
	var result, err = h.SendForm()
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return ioutil.ReadAll(result)
}

// Post request<POST 请求>
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

// Put request<PUT 请求>
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

// Delete request<DELETE 请求>
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

// Connect request<CONNECT 请求>
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

// Head request<HEAD 请求>
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

// Options request<OPTIONS 请求嗅探>
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

// Trace request<TRACE 请求>
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

// Patch request<PATCH 请求>
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

// SetUrlBuild <根据Map对象设置url地址拼接参数>
func SetUrlBuild(urls string, data map[string]string) string {
	u, _ := url.Parse(urls)
	q := u.Query()
	for k, v := range data {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// defaultFileName<默认文件名称>
func (h *HttpSend) defaultFileName() error {
	if h.FileName == "" {
		info, err := os.Stat(h.File) //Stat获取文件属性
		if err != nil {
			return err
		}
		h.FileName = info.Name()
	}
	return nil
}

// sendFile<发送单个文件>
func (h *HttpSend) sendFile(sendData io.Writer) error {
	//判断单个文件
	if h.File != "" && h.FileKey != "" {
		bodyWrite := multipart.NewWriter(sendData)
		file, err := os.Open(h.File)
		defer file.Close()
		if err != nil {
			return err
		}
		h.defaultFileName()

		fileWrite, err := bodyWrite.CreateFormFile(h.FileKey, h.FileName)
		if _, err = io.Copy(fileWrite, file); err != nil {
			return err
		}
		for key, val := range h.Body {
			if err = bodyWrite.WriteField(key, conv.String(val)); err != nil {
				return err
			}
		}
		err2 := bodyWrite.Close()
		if err2 != nil {
			return err2
		}
		h.ContentType = bodyWrite.FormDataContentType()
	}
	return nil
}

// sendFiles<发送多个文件>
func (h *HttpSend) sendFiles(sendData io.Writer) error {
	if len(h.Files) > 0 {
		bodyWrite := multipart.NewWriter(sendData)
		for _, val := range h.Files {
			file, err := os.Open(val)
			defer file.Close()
			if err != nil {
				return err
			}
			fileWrite, err := bodyWrite.CreateFormFile(h.FileKey, val)
			if _, err = io.Copy(fileWrite, file); err != nil {
				return err
			}
		}
		for key, val := range h.Body {
			if err := bodyWrite.WriteField(key, conv.String(val)); err != nil {
				return err
			}
		}
		bodyWrite.Close()
		h.ContentType = bodyWrite.FormDataContentType()
	}
	return nil
}

// SendForm <一般用于发送表单请求>
func (h *HttpSend) SendForm() (io.ReadCloser, error) {
	sendData := &bytes.Buffer{}
	var err error

	var Transport = &http.Transport{}
	if h.Proxy != nil {
		Transport.Proxy = h.Proxy
	}

	if h.Dial != nil {
		Transport.Dial = h.Dial
	}

	h.Client.Transport = Transport

	//发送单个文件
	if err2 := h.sendFile(sendData); err2 != nil {
		return nil, err2
	}

	//发送多个文件
	if err3 := h.sendFiles(sendData); err3 != nil {
		return nil, err3
	}

	h.Req, err = http.NewRequest(h.Method, h.Url, sendData)
	if err != nil {
		return nil, err
	}
	log.Println(h.ContentType)
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
	h.print()
	return h.Response.Body, nil
}

// PostFormFile Request file byte stream<请求文件字节流>
func (h *HttpSend) PostFormFile(url, files string) ([]byte, error) {
	h.Url = url
	h.Method = "POST"
	file, err := os.Open(files)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	h.Response, err = http.Post(url, "binary/octet-stream", file)
	if err != nil {
		return nil, err
	}
	defer h.Close()
	h.print()
	return ioutil.ReadAll(h.Response.Body)
}

// Send <扩展一般用于手动请求>
func (h *HttpSend) Send() (io.ReadCloser, error) {
	var err error
	var Transport = &http.Transport{}
	if h.Proxy != nil {
		Transport.Proxy = h.Proxy
	}

	if h.Dial != nil {
		Transport.Dial = h.Dial
	}

	h.Client.Transport = Transport

	h.Req, err = http.NewRequest(h.Method, h.Url, h.BodyReader)
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
	h.print()
	return h.Response.Body, nil
}

// Debug<用于最后打印>
func (h *HttpSend) Debug() *HttpSend {
	h.debug = true
	return h
}

// print<打印>
func (h *HttpSend) print() {
	if h.debug == true {
		body, _ := ioutil.ReadAll(h.Response.Body)
		fmt.Printf("[HttpRequest]\n")
		fmt.Printf("-------------------------------------------------------------------\n")
		fmt.Printf("Request: %s %s %s\nHeaders: %v\nCookies: %v\nTimeout: %ds\nReqBody: %v\n\n", h.Method, h.Url, h.Body,
			h.Header, h.Cookies, h.Timeout, string(body))
	}
}

// Close <必须默认关闭>
func (h *HttpSend) Close() {
	defer h.Response.Body.Close()
	if h.Req.Body != nil {
		defer h.Req.Body.Close()
	}
}
