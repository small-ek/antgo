package test

import (
	"github.com/small-ek/antgo/crypto/password"
	"testing"
)

var pwd = "123456"

func TestGenerate(t *testing.T) {

	// 测试生成哈希
	hashed, err := password.Generate(pwd)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// 验证生成的哈希是否有效
	if len(hashed) == 0 {
		t.Error("Generated hash is empty")
	}

	// 验证相同的密码生成不同的哈希
	hashed2, err := password.Generate(pwd)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if hashed == hashed2 {
		t.Error("Generated hashes are the same for the same password")
	}
}

func TestVerify(t *testing.T) {
	// 生成哈希
	hashed, err := password.Generate(pwd)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	// 验证正确的密码
	if !password.Verify(hashed, pwd) {
		t.Error("Verify failed for correct password")
	}

	// 验证错误的密码
	wrongPassword := "wrongpassword"
	if password.Verify(hashed, wrongPassword) {
		t.Error("Verify succeeded for incorrect password")
	}

	// 验证空密码
	emptyPassword := ""
	if password.Verify(hashed, emptyPassword) {
		t.Error("Verify succeeded for empty password")
	}
}

func TestVerifyWithInvalidHash(t *testing.T) {
	// 使用无效的哈希值进行验证
	invalidHash := "invalidhash"

	if password.Verify(invalidHash, pwd) {
		t.Error("Verify succeeded with invalid hash")
	}
}
