package ahttp

import (
	"context"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
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
	RetryWaitTime         time.Duration // 重试等待时间 / RetryWaitTime
	RetryMaxWaitTime      time.Duration // 最大重试等待时间 / RetryMaxWaitTime
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
		singletonClient.init()
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
	client.SetRetryWaitTime(config.RetryWaitTime)
	client.SetRetryMaxWaitTime(config.RetryMaxWaitTime)

	return client
}

// defaultConfig 返回默认的配置值 / Returns the default configuration values
func defaultConfig() *Config {
	return &Config{
		MaxIdleConnections:    runtime.GOMAXPROCS(0) * 200,
		IdleConnectionTimeout: 120 * time.Second,
		DisableCompression:    false,
		Timeout:               60 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
		MaxConnectionsPerHost: runtime.GOMAXPROCS(0) * 100,
		RetryAttempts:         3,
		DialerTimeout:         15 * time.Second, // 默认 Dialer 超时时间
		DialerKeepAlive:       60 * time.Second, // 默认 Dialer Keep-Alive 时间
		RetryWaitTime:         1 * time.Second,  // 默认重试等待时间
		RetryMaxWaitTime:      10 * time.Second, // 默认最大重试等待时间
		InsecureSkipVerify:    true,             // 跳过 TLS 证书验证 / Skip TLS certificate verification
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
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS13,
			CurvePreferences:   []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ExpectContinueTimeout: config.ExpectContinueTimeout,
		DisableCompression:    config.DisableCompression,
	}
}

// getRequestId 获取请求标识
func getRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	val := ctx.Value("request_id")
	requestID, ok := val.(string)
	if !ok || requestID == "" {
		return ""
	}
	return requestID
}

// SetLog 设置日志记录器 / SetLog sets the logger
func (h *HttpClient) SetLog(logger *zap.Logger) *HttpClient {
	if logger == nil {
		return h
	}
	h.logger = logger
	h.httpClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		requestID := getRequestId(r.Context())
		reqLogger := logger
		if requestID != "" {
			reqLogger = logger.With(zap.String("request_id", requestID))
		}
		reqLogger.Info("Request",
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
		respLog := []zap.Field{
			zap.Int("StatusCode", r.StatusCode()),
			zap.String("URL", r.Request.URL),
			zap.String("Method", r.Request.Method),
			zap.Any("Headers", r.Header()),
			zap.Any("Cookies", r.Cookies()),
			zap.Any("FormData", r.Request.FormData),
			zap.Any("QueryParam", r.Request.QueryParam),
			zap.ByteString("Body", r.Body()),
			zap.Duration("Duration", time.Since(r.Request.Time)),
		}
		requestID := getRequestId(r.Request.Context())
		reqLogger := logger
		if requestID != "" {
			reqLogger = logger.With(zap.String("request_id", requestID))
		}
		switch {
		case r.StatusCode() >= 500:
			reqLogger.Error("Server error", respLog...)
		case r.StatusCode() >= 400:
			reqLogger.Warn("Client error", respLog...)
		default:
			reqLogger.Info("Response success", respLog...)
		}
		return nil
	})
	return h
}

// Request 返回 Resty 请求实例 / Returns the Resty request instance
func (h *HttpClient) Request() *resty.Request {
	return h.httpClient.R()
}

// Client 返回 Resty 客户端实例 / Returns the Resty client instance
func (h *HttpClient) Client() *resty.Client {
	return h.httpClient
}

// init 初始化时统一设置 User-Agent 头部 / Set the User-Agent header during initialization
func (h *HttpClient) init() {
	h.httpClient.SetHeader("User-Agent", "antgo")
}

// SetCommonHeader 设置通用请求头
func (h *HttpClient) SetCommonHeader(key, value string) *HttpClient {
	h.httpClient.SetHeader(key, value)
	return h
}
