package ahttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
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
	Client      *http.Client           //Client
	Response    *http.Response         //Response
	Url         string                 //Request address<请求地址>
	ContentType string                 //Request type<网络文件的类型>
	Header      map[string]string      //Request header<请求头>
	Cookies     map[string]string      //Request Cookies<请求Cookies>
	Timeout     time.Duration          //Request timeout<请求超时时间>
	Body        map[string]interface{} //Request body<请求体>
	Method      string                 //Request method<请求类型>
	Files       []string               //Request Files <多个文件>
	File        string                 //Request File <单个文件>
	FileKey     string                 //Request FileKey<文件Key>
	FileName    string                 //Request FileName<文件名称>
	BodyReader  io.Reader              //Request BodyReader<读取器>
	Err         error
	StatusCode  int
	debug       bool
	curl        bool
}

var singletonHttpSend *HttpSend
var once sync.Once

// Client Default request
func Client() *HttpSend {
	once.Do(func() {
		singletonHttpSend = &HttpSend{
			ContentType: "application/json",
			Err:         nil,
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
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,
					DisableKeepAlives:     false,
				},
			},
		}
	})
	return singletonHttpSend
}

// SetClient <设置http.Client>
func (h *HttpSend) SetClient(client *http.Client) *HttpSend {
	h.Client = client
	return h
}

// SetTransport <设置http.Transport>
func (h *HttpSend) SetTransport(transport *http.Transport) *HttpSend {
	h.Client.Transport = transport
	return h
}

// SetBody Set body<设置请求体>
func (h *HttpSend) SetBody(body map[string]interface{}) *HttpSend {
	if h.Body == nil {
		h.Body = make(map[string]interface{})
	}

	configData, err := json.Marshal(body)
	if err != nil && alog.Write != nil {
		alog.Write.Error(h.Url, zap.String("SetBody error:", conv.String(body)), zap.Error(h.Err))
	}
	h.BodyReader = bytes.NewBuffer(configData)
	h.Body = body
	return h
}

// SetBodyReader Set body<设置读取器>
func (h *HttpSend) SetBodyReader(BodyReader io.Reader) *HttpSend {
	h.BodyReader = BodyReader
	return h
}

// SetUrl Set url<设置请求地址>
func (h *HttpSend) SetUrl(url string) *HttpSend {
	h.Url = url
	return h
}

// SetTimeout Set Timeout<设置超时>
func (h *HttpSend) SetTimeout(timeout time.Duration) *HttpSend {
	h.Timeout = timeout * time.Second
	return h
}

// SetHeader set header<设置请求头>
func (h *HttpSend) SetHeader(header map[string]string) *HttpSend {
	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	h.Header = header
	return h
}

// setHeader set header<设置请求头>
func (h *HttpSend) setHeader(req *http.Request) {
	if len(h.Header) > 0 {
		for k, v := range h.Header {
			if strings.ToLower(k) == "host" {
				req.Host = v
			} else {
				req.Header.Set(k, v)
			}
		}
	}

	if len(h.Header) == 0 {
		req.Header.Add("Content-Type", h.ContentType)
	}
}

// SetCookie set cookie<设置cookie>
func (h *HttpSend) SetCookie(cookies map[string]string) *HttpSend {
	if h.Cookies == nil {
		h.Cookies = make(map[string]string)
	}
	h.Cookies = cookies

	return h
}

// setCookie set cookie<设置cookie>
func (h *HttpSend) setCookie(req *http.Request) {
	if len(h.Cookies) > 0 {
		headerCookie := ""
		for k, v := range h.Cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			req.Header.Set("Cookie", headerCookie)
		}
	}
}

// SetMethod set method<设置请求类型>
func (h *HttpSend) SetMethod(method string) *HttpSend {
	h.Method = method
	return h
}

// SetContentType Set Type<设置资源的MIME类型>
func (h *HttpSend) SetContentType(ContentType string) *HttpSend {
	h.ContentType = ContentType
	return h
}

// Get request<GET 请求>
func (h *HttpSend) Get(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "GET"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// SetFiles Set Files<设置多个文件路径上传>
func (h *HttpSend) SetFiles(files []string) *HttpSend {
	h.Files = files
	return h
}

// SetFile Set File<设置单个文件上传路径>
func (h *HttpSend) SetFile(file string, key ...string) *HttpSend {
	if len(key) > 0 {
		h.FileKey = key[0]
	}
	h.File = file
	return h
}

// SetFileKey Set File<设置文件的Key>
func (h *HttpSend) SetFileKey(fileKey string) *HttpSend {
	h.FileKey = fileKey
	return h
}

// SetFileName Set File<设置文件的名称>
func (h *HttpSend) SetFileName(fileName string) *HttpSend {
	h.FileName = fileName
	return h
}

// SetFileKeyAndName Set File<设置文件Key和名称>
func (h *HttpSend) SetFileKeyAndName(fileKey, fileName string) *HttpSend {
	h.FileKey = fileKey
	h.FileName = fileName
	return h
}

// PostForm request<Post 表单提交>
func (h *HttpSend) PostForm(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "POST"
	result, err = h.SendForm()
	if err != nil {
		return nil, err
	}
	return
}

// Post request<POST 请求>
func (h *HttpSend) Post(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "POST"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// Put request<PUT 请求>
func (h *HttpSend) Put(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "PUT"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// Delete request<DELETE 请求>
func (h *HttpSend) Delete(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "DELETE"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// Connect request<CONNECT 请求>
func (h *HttpSend) Connect(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "CONNECT"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// Head request<HEAD 请求>
func (h *HttpSend) Head(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "HEAD"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// Options request<OPTIONS 请求嗅探>
func (h *HttpSend) Options(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "OPTIONS"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// Trace request<TRACE 请求>
func (h *HttpSend) Trace(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "TRACE"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
}

// Patch request<PATCH 请求>
func (h *HttpSend) Patch(url string) (result []byte, err error) {
	h.Url = url
	h.Method = "PATCH"
	result, err = h.Send()
	if err != nil {
		return nil, err
	}
	return
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
		info, err := os.Stat(h.File)
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

		if err := h.defaultFileName(); err != nil {
			return err
		}

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
		if err := bodyWrite.Close(); err != nil {
			return err
		}
		h.ContentType = bodyWrite.FormDataContentType()
	}
	return nil
}

// SendForm <一般用于发送表单请求>
func (h *HttpSend) SendForm() (body []byte, err error) {
	sendData := &bytes.Buffer{}

	//发送单个文件
	if err2 := h.sendFile(sendData); err2 != nil {
		return nil, err2
	}

	//发送多个文件
	if err3 := h.sendFiles(sendData); err3 != nil {
		return nil, err3
	}

	return h.request(h.Method, h.Url, sendData)
}

// PostFormFile Request file byte stream<请求文件字节流>
func (h *HttpSend) PostFormFile(url, files string) (body []byte, err error) {
	h.Url = url
	h.Method = "POST"
	file, err := os.Open(files)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "binary/octet-stream", file)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	h.print(body)

	if err != nil {
		return nil, err
	} else {
		resp.Body.Close()
	}
	return
}

// Send <扩展一般用于手动请求>
func (h *HttpSend) Send() (body []byte, err error) {
	return h.request(h.Method, h.Url, h.BodyReader)
}

// request <发送请求>
func (h *HttpSend) request(method, url string, readerBody io.Reader) (body []byte, err error) {
	req, err := http.NewRequest(method, url, readerBody)
	if err != nil {
		return nil, err
	}

	h.setCookie(req)
	h.setHeader(req)
	resp, err := h.Client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	h.StatusCode = resp.StatusCode
	h.Err = err
	defer h.print(body)
	defer h.curlPrint()
	return
}

// Debug <用于最后打印>
func (h *HttpSend) Debug(debug ...bool) *HttpSend {
	if len(debug) > 0 {
		h.debug = debug[0]
	} else {
		h.debug = true
	}
	return h
}

// Debug <用于最后打印>
func (h *HttpSend) Curl(curl ...bool) *HttpSend {
	if len(curl) > 0 {
		h.curl = curl[0]
	} else {
		h.curl = true
	}
	return h
}

func (h *HttpSend) curlPrint() {
	if !h.curl {
		return // 如果 debug 模式关闭，直接返回
	}
	fmt.Printf("[HttpRequest]\n")
	fmt.Printf("-------------------------------------------------------------------\n")
	fmt.Printf("curl -X %s '%s' \\\n", h.Method, h.Url)

	// 打印 Headers
	for k, v := range h.Header {
		fmt.Printf("  -H '%s: %s' \\\n", k, v)
	}

	// 打印 Body
	if h.Body != nil {
		bodyBytes, err := json.Marshal(h.Body)
		if err != nil {
			fmt.Printf("  -d 'serialization error: %s'\n", err.Error())
		} else {
			fmt.Printf("  -d '%s'\n", string(bodyBytes))
		}
	}

	// 打印 Error
	if h.Err != nil {
		fmt.Printf("[ERROR] %s\n", h.Err.Error())
	}
	fmt.Printf("-------------------------------------------------------------------\n")
}

// print<打印>
func (h *HttpSend) print(body []byte) {
	if !h.debug {
		return // 如果 debug 模式关闭，直接返回
	}

	// 准备日志数据
	logData := map[string]interface{}{
		"URL":          h.Url,
		"Method":       h.Method,
		"Timeout":      h.Timeout.String(),
		"Headers":      conv.String(h.Header),
		"Cookies":      conv.String(h.Cookies),
		"RequestBody":  conv.String(h.Body),
		"ResponseBody": string(body),
		"StatusCode":   h.StatusCode,
	}

	// 日志输出
	if alog.Write != nil {
		// 使用 zap 记录日志
		if h.Err != nil || h.StatusCode != 200 {
			alog.Write.Error("HTTP Request Error", zap.Any("LogData", logData))
		} else {
			alog.Write.Debug("HTTP Request Success", zap.Any("LogData", logData))
		}
	} else {
		// 控制台打印
		fmt.Printf("[HttpRequest]\n")
		fmt.Printf("-------------------------------------------------------------------\n")
		for key, value := range logData {
			fmt.Printf("%s: %v\n", key, value)
		}
		if h.Err != nil {
			fmt.Printf("Error: %s\n", h.Err.Error())
		}
	}
}

// errorPrint<错误打印>
func (h *HttpSend) errorPrint(body []byte) {
	if h.debug == true && alog.Write != nil {
		alog.Write.Error(h.Url, zap.String("Method:", h.Method), zap.String("Timeout", h.Timeout.String()), zap.String("Headers", conv.String(h.Header)), zap.String("Cookies", conv.String(h.Cookies)), zap.String("Body:", conv.String(h.Body)), zap.String("Response:", string(body[0])), zap.Int("statusCode", h.StatusCode), zap.Error(h.Err))
	}
}
