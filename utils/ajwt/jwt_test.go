package ajwt

import (
	"testing"
	"time"
)

// 用于测试的RSA密钥对（仅用于测试环境，请勿在生产环境中使用）
// Test RSA key pair (for testing purposes only, do not use in production)
const testRSAPrivateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAMMBf8hUaDk+oCsQ
ecIJDrRt+EktQaGE+11Un1YWgjelgJR68JcS+2HsSTjkd+aUYFuy7uo7Jc0ugQqN
xTwPzyJXfIX0J3niM0MM4SNLiqGD+UYGJ3bbBiw33NPT4CQ0/ATKk7Y4kdGXWl9z
OMqV6YqWau88SpFAUENHj1rJrFKrAgMBAAECgYBEFSTo61dMDSpcfq8T6JeitPZH
ji5o1wXvqtjKdKdYCEdhD58qD62GnblezJ1z+n+95DX3v1jOTxssdRzUgGx/Ys13
Ukso6hzNgMOw46zwd8qoWuFuytWWc953FGsr+atWFrvfU8aEjBjWzhlTtFVRPaiS
42zS7OZzYdmCw/PJaQJBAPu9R23SePtKIufwkg9ejWKyfIwJfLfdrfPL3RSTryGy
ph6CahWcr19/fMRTJaTYQtfxnXc1quFow3X+IBgfMr0CQQDGTmibUOJGyVrjSNMe
dTntWaBjMdFRwNdl7EzgUuLePUFs0gbZ6SW5dMsK/3JfzyxYF3XRki7Lupju3Ano
RGWHAkAN/lCJJ0j4Vv+nuvSzjAL5+If51NEs+1KfGbb5XNhAXEjlq0QwXVxWR6Ts
2N5f0nGsxU6GgOI103gCCBVKoflVAkBvBzVwSE/4XAJEIOD7O50MM9Ml1p2gjTzM
Nwovyph0340C9XCajvvtIuQPq0gJNoBYbgIsLRGARWAc1BvD7I9/AkEAgQEnpQEI
isXUlyKSsakm+M+hzkoJxlizUiM3tN9cIfsIBXdWv9LNGRp2gl8Sa69ri3EdqQXv
0PcStMOn2IX1kw==
-----END PRIVATE KEY-----`

const testRSAPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDAX/IVGg5PqArEHnCCQ60bfhJ
LUGhhPtdVJ9WFoI3pYCUevCXEvth7Ek45HfmlGBbsu7qOyXNLoEKjcU8D88iV3yF
9Cd54jNDDOEjS4qhg/lGBid22wYsN9zT0+AkNPwEypO2OJHRl1pfczjKlemKlmrv
PEqRQFBDR49ayaxSqwIDAQAB
-----END PUBLIC KEY-----`

// TestGenerateAndParse 测试生成和解析JWT令牌
// TestGenerateAndParse tests token generation and parsing.
func TestGenerateAndParse(t *testing.T) {
	// 初始化JwtManager并设置公钥和私钥
	// Initialize JwtManager and set RSA keys.
	jm := New().SetPrivateKey([]byte(testRSAPrivateKeyPEM)).SetPublicKey([]byte(testRSAPublicKeyPEM))

	// 设置自定义声明 / Define custom claims.
	claims := map[string]interface{}{
		"user": "john",
		"role": "admin",
	}

	// 生成Token / Generate token.
	token, err := jm.Generate(claims)
	if err != nil {
		t.Fatalf("Generate token failed: %v", err)
	}

	// 解析Token / Parse token.
	parsedClaims, err := jm.Parse(token)
	if err != nil {
		t.Fatalf("Parse token failed: %v", err)
	}

	// 验证声明中的数据 / Validate claims.
	if parsedClaims["user"] != "john" {
		t.Errorf("Expected user 'john', got %v", parsedClaims["user"])
	}
	if parsedClaims["role"] != "admin" {
		t.Errorf("Expected role 'admin', got %v", parsedClaims["role"])
	}
}

// TestInvalidToken 测试解析一个无效的Token
// TestInvalidToken tests parsing an invalid token string.
func TestInvalidToken(t *testing.T) {
	jm := New().SetPrivateKey([]byte(testRSAPrivateKeyPEM)).SetPublicKey([]byte(testRSAPublicKeyPEM))
	_, err := jm.Parse("invalid.token.here")
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}

// TestMissingKeys 测试缺少密钥时的行为
// TestMissingKeys tests the behavior when required keys are missing.
func TestMissingKeys(t *testing.T) {
	jm := New() // 未设置公钥和私钥
	claims := map[string]interface{}{
		"user": "john",
	}

	// 测试缺少私钥时生成Token应返回错误 / Generate should error if private key is not set.
	_, err := jm.Generate(claims)
	if err == nil {
		t.Error("Expected error for missing private key, got nil")
	}

	// 设置私钥，但不设置公钥
	jm.SetPrivateKey([]byte(testRSAPrivateKeyPEM))
	token, err := jm.Generate(claims)
	if err != nil {
		t.Fatalf("Generate token failed: %v", err)
	}

	// 解析Token时应返回缺少公钥的错误 / Parse should error if public key is not set.
	_, err = jm.Parse(token)
	if err == nil {
		t.Error("Expected error for missing public key, got nil")
	}
}

// TestCustomExpiration 测试自定义过期时间
// TestCustomExpiration tests token expiration with a custom expiration duration.
func TestCustomExpiration(t *testing.T) {
	// 初始化JwtManager，并设置公钥和私钥
	jm := New().SetPrivateKey([]byte(testRSAPrivateKeyPEM)).SetPublicKey([]byte(testRSAPublicKeyPEM))

	claims := map[string]interface{}{
		"user": "john",
	}

	// 生成一个1秒后过期的Token
	token, err := jm.Generate(claims, time.Second*1)
	if err != nil {
		t.Fatalf("Generate token failed: %v", err)
	}

	// 等待2秒以确保Token过期
	time.Sleep(2 * time.Second)

	_, err = jm.Parse(token)

	if err != nil {
		t.Error("Expected token to be expired, but no error returned")
	}
}
