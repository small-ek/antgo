package ajwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

// 常量声明 / Constant declarations
const (
	defaultExpiration = time.Hour * 24 * 6 // 默认过期时间 6 天 / Default expiration time: 6 days
	defaultClockSkew  = time.Minute * 1    // 允许的时钟偏移 / Allowed clock skew duration: 1 minute
)

// 错误定义 / Error definitions
var (
	ErrInvalidToken            = errors.New("invalid token")             // 无效的Token / Invalid token error
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method") // 未知的签名算法 / Unexpected signing method error
	ErrKeyNotSet               = errors.New("key not set")               // 未设置密钥 / Key not set error
	ErrKeyConfiguration        = errors.New("key configuration error")   // 密钥配置错误 / Key configuration error
)

// JwtManager JWT 管理器结构体 / JwtManager manages JWT operations.
type JwtManager struct {
	mu            sync.RWMutex      // 读写锁 / Mutex for concurrent access
	privateKey    *rsa.PrivateKey   // RSA 私钥 / RSA private key
	publicKey     *rsa.PublicKey    // RSA 公钥 / RSA public key
	defaultExp    time.Duration     // 默认过期时间 / Default expiration duration
	signingMethod jwt.SigningMethod // 签名算法 / Signing algorithm
	clockSkew     time.Duration     // 时钟偏移容忍时间 / Allowed clock skew duration
	initErrors    []error           // 初始化错误收集 / Collection of initialization errors
}

var instance *JwtManager
var once sync.Once

// New 创建单例实例 / New creates a singleton instance of JwtManager.
func New() *JwtManager {
	once.Do(func() {
		instance = &JwtManager{
			defaultExp:    defaultExpiration,
			signingMethod: jwt.SigningMethodRS256,
			clockSkew:     defaultClockSkew,
			initErrors:    make([]error, 0),
		}
	})
	return instance
}

// ================== 链式配置方法 / Chainable Configuration Methods ==================

// SetPublicKey 设置公钥（链式调用）/ SetPublicKey sets the RSA public key (chainable).
func (jm *JwtManager) SetPublicKey(pem []byte) *JwtManager {
	if len(pem) == 0 {
		// 收集错误 / Collect error
		jm.collectError(fmt.Errorf("%w: empty public key", ErrKeyConfiguration))
		return jm
	}

	// 解析 PEM 格式的公钥 / Parse the PEM encoded public key.
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		jm.collectError(fmt.Errorf("%w: %v", ErrKeyConfiguration, err))
		return jm
	}

	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.publicKey = pubKey
	return jm
}

// SetPrivateKey 设置私钥（链式调用）/ SetPrivateKey sets the RSA private key (chainable).
func (jm *JwtManager) SetPrivateKey(pem []byte) *JwtManager {
	if len(pem) == 0 {
		jm.collectError(fmt.Errorf("%w: empty private key", ErrKeyConfiguration))
		return jm
	}

	// 解析 PEM 格式的私钥 / Parse the PEM encoded private key.
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		jm.collectError(fmt.Errorf("%w: %v", ErrKeyConfiguration, err))
		return jm
	}

	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.privateKey = privKey
	return jm
}

// SetExpiration 设置默认过期时间（链式调用）/ SetExpiration sets the default token expiration duration (chainable).
func (jm *JwtManager) SetExpiration(d time.Duration) *JwtManager {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.defaultExp = d
	return jm
}

// ================== 核心方法 / Core Methods ==================

// Generate 生成JWT令牌 / Generate creates and signs a JWT token.
// 参数说明: claims 为自定义声明，exp 可选参数指定过期时间 /
// Parameters: claims is a map of custom claims, exp is an optional custom expiration duration.
func (jm *JwtManager) Generate(claims map[string]interface{}, exp ...time.Duration) (string, error) {
	// 先检查初始化错误 / Check for initialization errors first.
	if err := jm.checkInitErrors(); err != nil {
		return "", err
	}

	jm.mu.RLock()
	defer jm.mu.RUnlock()

	// 检查私钥是否设置 / Ensure private key is set.
	if jm.privateKey == nil {
		return "", fmt.Errorf("%w: private key required", ErrKeyNotSet)
	}

	// 获取当前时间 / Get current time.
	now := time.Now()
	// 创建标准声明 / Create standard claims.
	stdClaims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(jm.defaultExp)),
	}

	// 处理自定义过期时间 / Handle custom expiration time if provided.
	if len(exp) > 0 {
		stdClaims.ExpiresAt = jwt.NewNumericDate(now.Add(exp[0]))
	}

	// 合并用户自定义声明和标准声明 / Merge custom claims with standard claims.
	mapClaims := make(jwt.MapClaims)
	for k, v := range claims {
		mapClaims[k] = v
	}
	mapClaims["exp"] = stdClaims.ExpiresAt
	mapClaims["nbf"] = stdClaims.NotBefore
	mapClaims["iat"] = stdClaims.IssuedAt

	// 创建并签名 Token / Create and sign the token.
	token := jwt.NewWithClaims(jm.signingMethod, mapClaims)
	return token.SignedString(jm.privateKey)
}

// Parse 解析验证JWT令牌 / Parse verifies and parses a JWT token.
// 参数说明: tokenString 是待解析的 Token 字符串 / Parameter: tokenString is the JWT token string.
func (jm *JwtManager) Parse(tokenString string) (jwt.MapClaims, error) {
	// 先检查初始化错误 / Check for initialization errors first.
	if err := jm.checkInitErrors(); err != nil {
		return nil, err
	}

	jm.mu.RLock()
	defer jm.mu.RUnlock()

	// 检查公钥是否设置 / Ensure public key is set.
	if jm.publicKey == nil {
		return nil, fmt.Errorf("%w: public key required", ErrKeyNotSet)
	}

	// 解析 Token / Parse the token.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 验证签名算法 / Verify the signing algorithm.
		if t.Method.Alg() != jm.signingMethod.Alg() {
			return nil, ErrUnexpectedSigningMethod
		}
		return jm.publicKey, nil
	}, jwt.WithLeeway(jm.clockSkew))

	// 错误处理 / Handle errors.
	if err != nil {
		return nil, wrapJWTError(err)
	}

	// 如果解析成功且 Token 有效，则返回声明 / If valid, return the claims.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ================== 辅助方法 / Helper Methods ==================

// collectError 收集初始化错误 / collectError appends an error to the initialization errors slice.
func (jm *JwtManager) collectError(err error) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.initErrors = append(jm.initErrors, err)
}

// checkInitErrors 检查初始化错误 / checkInitErrors checks if any initialization errors have been collected.
func (jm *JwtManager) checkInitErrors() error {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	if len(jm.initErrors) > 0 {
		return fmt.Errorf("%w: %v", ErrKeyConfiguration, jm.initErrors[0])
	}
	return nil
}

// wrapJWTError 包装JWT错误 / wrapJWTError wraps JWT errors into more descriptive errors.
func wrapJWTError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return fmt.Errorf("%w: malformed", ErrInvalidToken)
	case errors.Is(err, jwt.ErrTokenExpired):
		return fmt.Errorf("%w: expired", ErrInvalidToken)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return fmt.Errorf("%w: not active yet", ErrInvalidToken)
	default:
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
}
