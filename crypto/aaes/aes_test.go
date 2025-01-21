package aaes

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"testing"
)

// 测试向量结构体
type testCase struct {
	name      string
	plaintext []byte
	key       []byte
	iv        []byte
	mode      string
	padding   string
	wantErr   bool
}

// 生成随机字节序列
func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// ======================== 通用测试函数 ========================

func testEncryptDecrypt(t *testing.T, tc testCase) {
	t.Run(tc.name, func(t *testing.T) {
		// 加密
		ciphertext, err := Encrypt(tc.plaintext, tc.key, tc.iv, tc.mode, tc.padding)
		if (err != nil) != tc.wantErr {
			t.Fatalf("Encrypt() error = %v, wantErr %v", err, tc.wantErr)
			return
		}
		if tc.wantErr {
			return
		}

		// 解密
		decrypted, err := Decrypt(ciphertext, tc.key, tc.iv, tc.mode, tc.padding)
		if (err != nil) != tc.wantErr {
			t.Fatalf("Decrypt() error = %v, wantErr %v", err, tc.wantErr)
			return
		}

		// 验证结果
		if !bytes.Equal(tc.plaintext, decrypted) {
			t.Errorf("Decrypted text mismatch\nGot:  %q\nWant: %q", decrypted, tc.plaintext)
		}
	})
}

// ======================== 具体模式测试 ========================

func TestCBC(t *testing.T) {
	key128 := randomBytes(16)
	key192 := randomBytes(24)
	key256 := randomBytes(32)
	iv := randomBytes(aes.BlockSize)

	tests := []testCase{
		{
			name:      "CBC-PKCS7-128",
			plaintext: []byte("Hello CBC Mode!"),
			key:       key128,
			iv:        iv,
			mode:      ModeCBC,
			padding:   PaddingPKCS7,
		},
		{
			name:      "CBC-ISO10126-192",
			plaintext: []byte("Hello CBC-ISO10126-192"), // 随机35字节
			key:       key192,
			iv:        iv,
			mode:      ModeCBC,
			padding:   PaddingISO10126,
		},
		{
			name:      "CBC-ANSIX923-256",
			plaintext: []byte("X923 padding test"),
			key:       key256,
			iv:        iv,
			mode:      ModeCBC,
			padding:   PaddingANSIX923,
		},
		{
			name:      "Invalid-IV",
			plaintext: []byte("test"),
			key:       key128,
			iv:        []byte("short iv"),
			mode:      ModeCBC,
			padding:   PaddingPKCS7,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		testEncryptDecrypt(t, tc)
	}
}

func TestECB(t *testing.T) {
	key := randomBytes(32)
	plaintext := []byte("ECB Mode Testing!")

	tests := []testCase{
		{
			name:      "ECB-PKCS7",
			plaintext: plaintext,
			key:       key,
			iv:        []byte{},
			mode:      ModeECB,
			padding:   PaddingPKCS7,
		},
		{
			name:      "ECB-ZeroPad",
			plaintext: []byte("DataWithZero"),
			key:       key,
			iv:        []byte{},
			mode:      ModeECB,
			padding:   PaddingZero,
		},
		{
			name:      "ECB-NoPaddingError",
			plaintext: []byte("InvalidLength"), // 13 bytes
			key:       key,
			iv:        []byte{},
			mode:      ModeECB,
			padding:   PaddingNone,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		testEncryptDecrypt(t, tc)
	}
}

func TestCTR(t *testing.T) {
	key := randomBytes(24)
	iv := randomBytes(aes.BlockSize)

	tests := []testCase{
		{
			name:      "CTR-NoPadding",
			plaintext: []byte("CTR Stream Mode"),
			key:       key,
			iv:        iv,
			mode:      ModeCTR,
			padding:   PaddingNone,
		},
		{
			name:      "CTR-LongText",
			plaintext: bytes.Repeat([]byte("A"), 12345),
			key:       key,
			iv:        iv,
			mode:      ModeCTR,
			padding:   PaddingNone,
		},
		{
			name:      "CTR-EmptyData",
			plaintext: []byte{},
			key:       key,
			iv:        iv,
			mode:      ModeCTR,
			padding:   PaddingNone,
		},
	}

	for _, tc := range tests {
		testEncryptDecrypt(t, tc)
	}
}

func TestOFB(t *testing.T) {
	key := randomBytes(16)
	iv := randomBytes(aes.BlockSize)

	tests := []testCase{
		{
			name:      "OFB-Simple",
			plaintext: []byte("OFB Mode Test"),
			key:       key,
			iv:        iv,
			mode:      ModeOFB,
			padding:   PaddingNone,
		},
		{
			name:      "OFB-RandomData",
			plaintext: randomBytes(255),
			key:       key,
			iv:        iv,
			mode:      ModeOFB,
			padding:   PaddingNone,
		},
	}

	for _, tc := range tests {
		testEncryptDecrypt(t, tc)
	}
}

func TestCFB(t *testing.T) {
	key := randomBytes(32)
	iv := randomBytes(aes.BlockSize)

	tests := []testCase{
		{
			name:      "CFB-Normal",
			plaintext: []byte("CFB Mode Example"),
			key:       key,
			iv:        iv,
			mode:      ModeCFB,
			padding:   PaddingNone,
		},
		{
			name:      "CFB-SpecialChars",
			plaintext: []byte("\x00\xff\x55\xaa\x11\x22\x33"),
			key:       key,
			iv:        iv,
			mode:      ModeCFB,
			padding:   PaddingNone,
		},
	}

	for _, tc := range tests {
		testEncryptDecrypt(t, tc)
	}
}

// ======================== 特殊案例测试 ========================

func TestEdgeCases(t *testing.T) {
	key := randomBytes(16)
	iv := randomBytes(aes.BlockSize)

	tests := []testCase{
		{
			name:      "EmptyPlaintext",
			plaintext: []byte{},
			key:       key,
			iv:        iv,
			mode:      ModeCBC,
			padding:   PaddingPKCS7,
		},
		{
			name:      "ExactBlockSize",
			plaintext: bytes.Repeat([]byte{0x41}, aes.BlockSize),
			key:       key,
			iv:        iv,
			mode:      ModeCBC,
			padding:   PaddingPKCS7,
		},
		{
			name:      "SpacePadding",
			plaintext: []byte("Text with spaces"),
			key:       key,
			iv:        iv,
			mode:      ModeCBC,
			padding:   PaddingSpace,
		},
	}

	for _, tc := range tests {
		testEncryptDecrypt(t, tc)
	}
}

func TestInvalidInputs(t *testing.T) {
	validKey := randomBytes(16)
	validIV := randomBytes(aes.BlockSize)

	tests := []testCase{
		{
			name:      "InvalidKeyLength",
			plaintext: []byte("test"),
			key:       []byte("short"),
			iv:        validIV,
			mode:      ModeCBC,
			padding:   PaddingPKCS7,
			wantErr:   true,
		},
		{
			name:      "InvalidMode",
			plaintext: []byte("test"),
			key:       validKey,
			iv:        validIV,
			mode:      "INVALID",
			padding:   PaddingPKCS7,
			wantErr:   true,
		},
		{
			name:      "InvalidPadding",
			plaintext: []byte("test"),
			key:       validKey,
			iv:        validIV,
			mode:      ModeCBC,
			padding:   "UNKNOWN",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		testEncryptDecrypt(t, tc)
	}
}

// ======================== 性能测试 ========================

func BenchmarkCBC(b *testing.B) {
	key := randomBytes(32)
	iv := randomBytes(aes.BlockSize)
	data := randomBytes(1024 * 1024) // 1MB数据

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encrypt(data, key, iv, ModeCBC, PaddingPKCS7)
	}
}

func BenchmarkCTR(b *testing.B) {
	key := randomBytes(32)
	iv := randomBytes(aes.BlockSize)
	data := randomBytes(1024 * 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encrypt(data, key, iv, ModeCTR, PaddingNone)
	}
}
