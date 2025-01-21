package aaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"unicode/utf8"
)

// Encryption Modes 加密模式
const (
	ModeCBC = "CBC" // Cipher Block Chaining (密码分组链接)
	ModeECB = "ECB" // Electronic Codebook (电子密码本)
	ModeCTR = "CTR" // Counter (计数器模式)
	ModeOFB = "OFB" // Output Feedback (输出反馈)
	ModeCFB = "CFB" // Cipher Feedback (密码反馈)
)

// Padding Types 填充类型
const (
	PaddingPKCS7    = "PKCS7"    // PKCS#7/PKCS5标准填充
	PaddingISO10126 = "ISO10126" // ISO 10126随机填充
	PaddingANSIX923 = "ANSIX923" // ANSI X.923零填充
	PaddingZero     = "Zero"     // 零字节填充
	PaddingSpace    = "Space"    // 空格字符填充
	PaddingNone     = "None"     // 无填充
)

// Error Definitions 错误类型
var (
	ErrInvalidKeyLength    = errors.New("invalid key length (16/24/32 bytes required)")
	ErrInvalidIVLength     = errors.New("invalid IV length")
	ErrInvalidPadding      = errors.New("invalid padding detected")
	ErrUnsupportedMode     = errors.New("unsupported encryption mode")
	ErrUnsupportedPadding  = errors.New("unsupported padding type")
	ErrInvalidDataLength   = errors.New("invalid data length")
	ErrPaddingSizeMismatch = errors.New("padding size mismatch")
)

// Encrypt AES加密
func Encrypt(plaintext, key, iv []byte, mode string, padding string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}

	blockSize := block.BlockSize()

	// 参数校验
	if err := validateParams(mode, iv, blockSize); err != nil {
		return nil, err
	}

	// 处理填充
	if requiresPadding(mode) {
		if plaintext, err = applyPadding(plaintext, blockSize, padding); err != nil {
			return nil, err
		}
	}

	// 执行加密
	switch mode {
	case ModeCBC:
		return encryptCBC(block, plaintext, iv)
	case ModeECB:
		return encryptECB(block, plaintext)
	case ModeCTR:
		return encryptCTR(block, plaintext, iv)
	case ModeOFB:
		return encryptOFB(block, plaintext, iv)
	case ModeCFB:
		return encryptCFB(block, plaintext, iv)
	default:
		return nil, ErrUnsupportedMode
	}
}

// Decrypt AES解密
func Decrypt(ciphertext, key, iv []byte, mode string, padding string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}

	blockSize := block.BlockSize()

	// 参数校验
	if err := validateParams(mode, iv, blockSize); err != nil {
		return nil, err
	}

	// 执行解密
	var plaintext []byte
	switch mode {
	case ModeCBC:
		if plaintext, err = decryptCBC(block, ciphertext, iv); err != nil {
			return nil, err
		}
	case ModeECB:
		if plaintext, err = decryptECB(block, ciphertext); err != nil {
			return nil, err
		}
	case ModeCTR:
		plaintext = decryptCTR(block, ciphertext, iv)
	case ModeOFB:
		plaintext = decryptOFB(block, ciphertext, iv)
	case ModeCFB:
		if plaintext, err = decryptCFB(block, ciphertext, iv); err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedMode
	}

	// 去除填充
	if requiresPadding(mode) {
		return removePadding(plaintext, blockSize, padding)
	}
	return plaintext, nil
}

// ======================== 加密模式实现 ========================

func encryptCBC(block cipher.Block, plaintext, iv []byte) ([]byte, error) {
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func decryptCBC(block cipher.Block, ciphertext, iv []byte) ([]byte, error) {
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("ciphertext length not multiple of block size(%d)", block.BlockSize())
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}

func encryptECB(block cipher.Block, plaintext []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(plaintext)%blockSize != 0 {
		return nil, fmt.Errorf("plaintext length must be multiple of %d", blockSize)
	}

	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += blockSize {
		block.Encrypt(ciphertext[i:], plaintext[i:])
	}
	return ciphertext, nil
}

func decryptECB(block cipher.Block, ciphertext []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext length must be multiple of %d", blockSize)
	}

	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += blockSize {
		block.Decrypt(plaintext[i:], ciphertext[i:])
	}
	return plaintext, nil
}

func encryptCTR(block cipher.Block, plaintext, iv []byte) ([]byte, error) {
	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

func decryptCTR(block cipher.Block, ciphertext, iv []byte) []byte {
	// CTR模式加解密相同
	plaintext, _ := encryptCTR(block, ciphertext, iv)
	return plaintext
}

func encryptOFB(block cipher.Block, plaintext, iv []byte) ([]byte, error) {
	stream := cipher.NewOFB(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

func decryptOFB(block cipher.Block, ciphertext, iv []byte) []byte {
	// OFB模式加解密相同
	plaintext, _ := encryptOFB(block, ciphertext, iv)
	return plaintext
}

func encryptCFB(block cipher.Block, plaintext, iv []byte) ([]byte, error) {
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

func decryptCFB(block cipher.Block, ciphertext, iv []byte) ([]byte, error) {
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

// ======================== 填充方法 ========================

func applyPadding(data []byte, blockSize int, paddingType string) ([]byte, error) {
	switch paddingType {
	case PaddingPKCS7:
		return pkcs7Pad(data, blockSize), nil
	case PaddingISO10126:
		return iso10126Pad(data, blockSize)
	case PaddingANSIX923:
		return ansiX923Pad(data, blockSize)
	case PaddingZero:
		return zeroPad(data, blockSize)
	case PaddingSpace:
		return spacePad(data, blockSize)
	case PaddingNone:
		if len(data)%blockSize != 0 {
			return nil, ErrInvalidDataLength
		}
		return data, nil
	default:
		return nil, ErrUnsupportedPadding
	}
}

func removePadding(data []byte, blockSize int, paddingType string) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPadding
	}

	switch paddingType {
	case PaddingPKCS7:
		return pkcs7Unpad(data, blockSize)
	case PaddingISO10126:
		return iso10126Unpad(data, blockSize)
	case PaddingANSIX923:
		return ansiX923Unpad(data, blockSize)
	case PaddingZero:
		return zeroUnpad(data, blockSize)
	case PaddingSpace:
		return spaceUnpad(data, blockSize)
	case PaddingNone:
		return data, nil
	default:
		return nil, ErrUnsupportedPadding
	}
}

// PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, ErrInvalidPadding
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, ErrInvalidPadding
	}

	if !bytes.HasSuffix(data, bytes.Repeat([]byte{byte(padding)}, padding)) {
		return nil, ErrInvalidPadding
	}
	return data[:len(data)-padding], nil
}

// ISO 10126填充
func iso10126Pad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		padding = blockSize
	}

	buf := make([]byte, len(data)+padding)
	copy(buf, data)

	if _, err := rand.Read(buf[len(data) : len(data)+padding-1]); err != nil {
		return nil, err
	}
	buf[len(buf)-1] = byte(padding)
	return buf, nil
}

func iso10126Unpad(data []byte, _ int) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPadding
	}

	padding := int(data[len(data)-1])
	if padding < 1 {
		return nil, ErrInvalidPadding
	}

	if len(data) < padding {
		return nil, ErrInvalidPadding
	}
	return data[:len(data)-padding], nil
}

// ANSI X.923填充
func ansiX923Pad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		padding = blockSize
	}

	buf := make([]byte, len(data)+padding)
	copy(buf, data)
	buf[len(buf)-1] = byte(padding)
	return buf, nil
}

func ansiX923Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPadding
	}

	padding := int(data[len(data)-1])
	if padding < 1 || padding > blockSize {
		return nil, ErrInvalidPadding
	}

	if len(data) < padding {
		return nil, ErrInvalidPadding
	}

	for i := len(data) - padding; i < len(data)-1; i++ {
		if data[i] != 0 {
			return nil, ErrInvalidPadding
		}
	}
	return data[:len(data)-padding], nil
}

// 零填充
func zeroPad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		return data, nil
	}
	return append(data, bytes.Repeat([]byte{0}, padding)...), nil
}

func zeroUnpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPadding
	}

	// 查找最后一个非零字节的位置
	idx := bytes.LastIndexFunc(data, func(r rune) bool {
		return r != 0
	})

	// 全零数据特殊处理
	if idx == -1 {
		return []byte{}, nil
	}

	// 计算填充长度并验证
	paddingLen := len(data) - idx - 1
	if paddingLen < 0 {
		return nil, ErrInvalidPadding
	}

	// 验证填充字节是否全零
	for _, b := range data[idx+1:] {
		if b != 0 {
			return nil, ErrInvalidPadding
		}
	}

	// 检查填充长度是否合法
	if paddingLen > blockSize {
		return nil, ErrPaddingSizeMismatch
	}

	return data[:idx+1], nil
}

// 空格填充
func spacePad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		return data, nil
	}
	return append(data, bytes.Repeat([]byte{' '}, padding)...), nil
}

func spaceUnpad(data []byte, blockSize int) ([]byte, error) {
	idx := len(data) - 1
	for ; idx >= 0; idx-- {
		if data[idx] != ' ' {
			break
		}
	}

	padding := len(data) - idx - 1
	if padding == 0 {
		return data, nil
	}

	// 验证填充字符有效性
	if !utf8.Valid(data[idx+1:]) {
		return nil, ErrInvalidPadding
	}

	if padding > blockSize {
		return nil, ErrPaddingSizeMismatch
	}
	return data[:idx+1], nil
}

// ======================== 辅助函数 ========================

func validateParams(mode string, iv []byte, blockSize int) error {
	switch mode {
	case ModeCBC, ModeCTR, ModeOFB, ModeCFB:
		if len(iv) != blockSize {
			return fmt.Errorf("%w: need %d bytes", ErrInvalidIVLength, blockSize)
		}
	case ModeECB:
		if len(iv) != 0 {
			return fmt.Errorf("%w: ECB模式需要空IV", ErrInvalidIVLength)
		}
	default:
		return ErrUnsupportedMode
	}
	return nil
}

func requiresPadding(mode string) bool {
	return mode == ModeCBC || mode == ModeECB
}
