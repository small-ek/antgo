package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

// 常量声明
const (
	defaultExpiration = time.Hour * 24 * 6 // 默认过期时间 6 天
	defaultClockSkew  = time.Minute * 1    // 允许的时钟偏移
)

// 错误定义
var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrKeyNotSet               = errors.New("key not set")
	ErrKeyConfiguration        = errors.New("key configuration error")
)

// JwtManager JWT 管理器结构体
type JwtManager struct {
	mu            sync.RWMutex      // 读写锁
	privateKey    *rsa.PrivateKey   // RSA 私钥
	publicKey     *rsa.PublicKey    // RSA 公钥
	defaultExp    time.Duration     // 默认过期时间
	signingMethod jwt.SigningMethod // 签名算法
	clockSkew     time.Duration     // 时钟偏移容忍时间
	initErrors    []error           // 初始化错误收集
}

var instance *JwtManager
var once sync.Once

// New 创建单例实例
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

// ================== 链式配置方法 ==================

// SetPublicKey 设置公钥（链式调用）
func (jm *JwtManager) SetPublicKey(pem []byte) *JwtManager {
	if len(pem) == 0 {
		jm.collectError(fmt.Errorf("%w: empty public key", ErrKeyConfiguration))
		return jm
	}

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

// SetPrivateKey 设置私钥（链式调用）
func (jm *JwtManager) SetPrivateKey(pem []byte) *JwtManager {
	if len(pem) == 0 {
		jm.collectError(fmt.Errorf("%w: empty private key", ErrKeyConfiguration))
		return jm
	}

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

// SetExpiration 设置默认过期时间（链式调用）
func (jm *JwtManager) SetExpiration(d time.Duration) *JwtManager {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.defaultExp = d
	return jm
}

// ================== 核心方法 ==================

// Generate 生成JWT令牌
func (jm *JwtManager) Generate(claims map[string]interface{}, exp ...time.Duration) (string, error) {
	// 先检查初始化错误
	if err := jm.checkInitErrors(); err != nil {
		return "", err
	}

	jm.mu.RLock()
	defer jm.mu.RUnlock()

	if jm.privateKey == nil {
		return "", fmt.Errorf("%w: private key required", ErrKeyNotSet)
	}

	// 创建标准声明
	now := time.Now()
	stdClaims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(jm.defaultExp)),
	}

	// 处理自定义过期时间
	if len(exp) > 0 {
		stdClaims.ExpiresAt = jwt.NewNumericDate(now.Add(exp[0]))
	}

	// 合并声明
	mapClaims := make(jwt.MapClaims)
	for k, v := range claims {
		mapClaims[k] = v
	}
	mapClaims["exp"] = stdClaims.ExpiresAt
	mapClaims["nbf"] = stdClaims.NotBefore
	mapClaims["iat"] = stdClaims.IssuedAt

	// 创建并签名Token
	token := jwt.NewWithClaims(jm.signingMethod, mapClaims)
	return token.SignedString(jm.privateKey)
}

// Parse 解析验证JWT令牌
func (jm *JwtManager) Parse(tokenString string) (jwt.MapClaims, error) {
	// 先检查初始化错误
	if err := jm.checkInitErrors(); err != nil {
		return nil, err
	}

	jm.mu.RLock()
	defer jm.mu.RUnlock()

	if jm.publicKey == nil {
		return nil, fmt.Errorf("%w: public key required", ErrKeyNotSet)
	}

	// 解析Token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jm.signingMethod.Alg() {
			return nil, ErrUnexpectedSigningMethod
		}
		return jm.publicKey, nil
	}, jwt.WithLeeway(jm.clockSkew))

	// 错误处理
	if err != nil {
		return nil, wrapJWTError(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ================== 辅助方法 ==================

// 收集初始化错误
func (jm *JwtManager) collectError(err error) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.initErrors = append(jm.initErrors, err)
}

// 检查初始化错误
func (jm *JwtManager) checkInitErrors() error {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	if len(jm.initErrors) > 0 {
		return fmt.Errorf("%w: %v", ErrKeyConfiguration, jm.initErrors[0])
	}
	return nil
}

// 包装JWT错误
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
