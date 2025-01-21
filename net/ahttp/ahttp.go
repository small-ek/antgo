package ahttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Config 包含 HTTP 客户端的配置选项
// Config contains the configuration options for the HTTP client
type Config struct {
	MaxIdleConnections    int           // 连接池的最大空闲连接数 / Maximum idle connections in the connection pool
	IdleConnectionTimeout time.Duration // 空闲连接超时时间 / Timeout for idle connections
	DisableCompression    bool          // 禁用压缩 / Disable compression
	DisableKeepAlives     bool          // 禁用 keep-alive / Disable keep-alive
	InsecureSkipVerify    bool          // 跳过 TLS 证书验证 / Skip TLS certificate verification
	Timeout               time.Duration // 总超时时间 / Total timeout duration
	TLSHandshakeTimeout   time.Duration // TLS 握手超时时间 / Timeout for TLS handshake
	ExpectContinueTimeout time.Duration // 100-continue 超时时间 / Timeout for 100-continue
	MaxConnectionsPerHost int           // 每主机的最大连接数 / Maximum connections per host
	RetryAttempts         int           // 请求重试次数 / Number of retry attempts for failed requests
	DialerTimeout         time.Duration // Dialer 的连接超时时间 / Dialer connection timeout
	DialerKeepAlive       time.Duration // Dialer 的 Keep-Alive 时间 / Dialer keep-alive time
}

// HttpClient 是对 Resty 客户端的封装，支持自定义配置和日志 / HttpClient is a wrapper for the Resty client with custom configuration and logging support
type HttpClient struct {
	httpClient    *resty.Client   // Resty 客户端实例 / Resty client instance
	httpTransport *http.Transport // HTTP 传输层配置 / HTTP transport layer configuration
	logger        *zap.Logger     // 日志记录器 / Logger instance
	config        *Config
}

var (
	singletonClient *HttpClient
	once            sync.Once
)

// New 创建一个新的 HTTP 客户端实例 / NewHttpClient creates a new HTTP client instance
func New(config *Config) *HttpClient {
	once.Do(func() {
		if config == nil {
			config = defaultConfig()
		}

		singletonClient = &HttpClient{
			httpClient:    newRestyClient(config),
			httpTransport: newTransport(config),
			config:        config,
		}
	})
	return singletonClient
}

// newRestyClient 配置 Resty 客户端 / Configures the Resty client
func newRestyClient(config *Config) *resty.Client {
	client := resty.NewWithClient(&http.Client{
		Transport: newTransport(config),
		Timeout:   config.Timeout,
	})

	client.SetRetryCount(config.RetryAttempts)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(10 * time.Second)

	return client
}

// defaultConfig 返回默认的配置值 / Returns the default configuration values
func defaultConfig() *Config {
	return &Config{
		MaxIdleConnections:    runtime.GOMAXPROCS(0) * 200,
		IdleConnectionTimeout: 120 * time.Second,
		DisableCompression:    false,
		Timeout:               90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
		MaxConnectionsPerHost: runtime.GOMAXPROCS(0) * 100,
		RetryAttempts:         3,
		DialerTimeout:         15 * time.Second, // 默认 Dialer 超时时间
		DialerKeepAlive:       60 * time.Second, // 默认 Dialer Keep-Alive 时间
	}
}

// newTransport 创建并配置 HTTP 传输层 / Creates and configures the HTTP transport
func newTransport(config *Config) *http.Transport {
	tcpDialer := &net.Dialer{
		Timeout:   config.DialerTimeout,
		KeepAlive: config.DialerKeepAlive,
		DualStack: true,
	}

	return &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         tcpDialer.DialContext,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        config.MaxIdleConnections,
		MaxIdleConnsPerHost: config.MaxConnectionsPerHost,
		IdleConnTimeout:     config.IdleConnectionTimeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
			MinVersion:         tls.VersionTLS12,
			CurvePreferences:   []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ExpectContinueTimeout: config.ExpectContinueTimeout,
		DisableCompression:    config.DisableCompression,
	}
}

// SetLog 设置日志记录器 / SetLog sets the logger
func (h *HttpClient) SetLog(logger *zap.Logger) *HttpClient {
	h.logger = logger
	h.httpClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		logger.Info("Request",
			zap.String("URL", r.URL),
			zap.String("Method", r.Method),
			zap.Any("Headers", r.Header),
			zap.Any("Cookies", r.Cookies),
			zap.Any("FormData", r.FormData),
			zap.Any("QueryParam", r.QueryParam),
			zap.Any("Body", r.Body),
			zap.Duration("Timeout", h.config.Timeout))
		return nil
	})

	h.httpClient.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.StatusCode() >= 400 {
			logger.Error("Request failed",
				zap.Int("StatusCode", r.StatusCode()),
				zap.String("URL", r.Request.URL),
				zap.String("Method", r.Request.Method),
				zap.Any("Headers", r.Header),
				zap.Any("Cookies", r.Cookies),
				zap.Any("FormData", r.Request.FormData),
				zap.Any("QueryParam", r.Request.QueryParam),
				zap.ByteString("Body", r.Body()),
				zap.Duration("Duration", time.Since(r.Request.Time)),
			)
		} else {
			logger.Info("Response",
				zap.Int("StatusCode", r.StatusCode()),
				zap.String("URL", r.Request.URL),
				zap.String("Method", r.Request.Method),
				zap.Any("Headers", r.Header),
				zap.Any("Cookies", r.Cookies),
				zap.Any("FormData", r.Request.FormData),
				zap.Any("QueryParam", r.Request.QueryParam),
				zap.ByteString("Body", r.Body()),
				zap.Duration("Duration", time.Since(r.Request.Time)),
			)
		}
		return nil
	})
	return h
}

// Client 返回 Resty 客户端的请求实例 / Returns a request instance of the Resty client
func (h *HttpClient) Client() *resty.Request {
	h.httpClient.SetHeader("User-Agent", "antgo")
	return h.httpClient.R()
}
