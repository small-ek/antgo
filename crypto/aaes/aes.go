package aaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	ModeCBC = "CBC"
	ModeECB = "ECB"
	ModeCTR = "CTR"
	ModeOFB = "OFB"
	ModeCFB = "CFB"
)

const (
	PaddingPKCS7    = "PKCS7"
	PaddingPKCS5    = "PKCS5"
	PaddingISO10126 = "ISO10126"
	PaddingANSIX923 = "ANSIX923"
	PaddingZERO     = "ZERO"
	PaddingSPACE    = "SPACE"
	PaddingNONE     = "NONE"
)

var (
	ErrInvalidKeyLength    = errors.New("invalid key length (16/24/32 bytes required)")
	ErrInvalidIVLength     = errors.New("invalid IV length")
	ErrInvalidPadding      = errors.New("invalid padding detected")
	ErrUnsupportedMode     = errors.New("unsupported encryption mode")
	ErrUnsupportedPadding  = errors.New("unsupported padding type")
	ErrInvalidDataLength   = errors.New("invalid data length")
	ErrPaddingSizeMismatch = errors.New("padding size mismatch")
)

// Encrypt 加密入口
func Encrypt(plaintext, key, iv []byte, mode, padding string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}
	blockSize := block.BlockSize()

	if err := validateParams(mode, iv, blockSize); err != nil {
		return nil, err
	}

	if requiresPadding(mode) {
		plaintext, err = applyPadding(plaintext, blockSize, padding)
		if err != nil {
			return nil, err
		}
	} else if len(plaintext) == 0 {
		// 防止 CTR/OFB/CFB 空数据加密报错
		return []byte{}, nil
	}

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

// Decrypt 解密入口
func Decrypt(ciphertext, key, iv []byte, mode, padding string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}
	blockSize := block.BlockSize()

	if err := validateParams(mode, iv, blockSize); err != nil {
		return nil, err
	}

	if len(ciphertext) == 0 {
		return []byte{}, nil
	}

	var plaintext []byte
	switch mode {
	case ModeCBC:
		plaintext, err = decryptCBC(block, ciphertext, iv)
	case ModeECB:
		plaintext, err = decryptECB(block, ciphertext)
	case ModeCTR:
		plaintext, err = decryptCTR(block, ciphertext, iv)
	case ModeOFB:
		plaintext, err = decryptOFB(block, ciphertext, iv)
	case ModeCFB:
		plaintext, err = decryptCFB(block, ciphertext, iv)
	default:
		return nil, ErrUnsupportedMode
	}
	if err != nil {
		return nil, err
	}

	if requiresPadding(mode) {
		return removePadding(plaintext, blockSize, padding)
	}
	return plaintext, nil
}

// === 加解密模式 ===

func encryptCBC(block cipher.Block, plaintext, iv []byte) ([]byte, error) {
	ciphertext := make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func decryptCBC(block cipher.Block, ciphertext, iv []byte) ([]byte, error) {
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("%w: ciphertext length %d not multiple of block size %d", ErrInvalidDataLength, len(ciphertext), block.BlockSize())
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}

func encryptECB(block cipher.Block, plaintext []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(plaintext)%blockSize != 0 {
		return nil, fmt.Errorf("%w: plaintext length %d not multiple of %d", ErrInvalidDataLength, len(plaintext), blockSize)
	}
	ciphertext := make([]byte, len(plaintext))
	for start := 0; start < len(plaintext); start += blockSize {
		block.Encrypt(ciphertext[start:], plaintext[start:])
	}
	return ciphertext, nil
}

func decryptECB(block cipher.Block, ciphertext []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("%w: ciphertext length %d not multiple of %d", ErrInvalidDataLength, len(ciphertext), blockSize)
	}
	plaintext := make([]byte, len(ciphertext))
	for start := 0; start < len(ciphertext); start += blockSize {
		block.Decrypt(plaintext[start:], ciphertext[start:])
	}
	return plaintext, nil
}

// CTR、OFB 模式加解密一致，统一写法

func encryptCTR(block cipher.Block, input, iv []byte) ([]byte, error) {
	out := make([]byte, len(input))
	cipher.NewCTR(block, iv).XORKeyStream(out, input)
	return out, nil
}

func decryptCTR(block cipher.Block, input, iv []byte) ([]byte, error) {
	// CTR 对称
	return encryptCTR(block, input, iv)
}

func encryptOFB(block cipher.Block, input, iv []byte) ([]byte, error) {
	out := make([]byte, len(input))
	cipher.NewOFB(block, iv).XORKeyStream(out, input)
	return out, nil
}

func decryptOFB(block cipher.Block, input, iv []byte) ([]byte, error) {
	return encryptOFB(block, input, iv)
}

func encryptCFB(block cipher.Block, input, iv []byte) ([]byte, error) {
	out := make([]byte, len(input))
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(out, input)
	return out, nil
}

func decryptCFB(block cipher.Block, input, iv []byte) ([]byte, error) {
	out := make([]byte, len(input))
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(out, input)
	return out, nil
}

// === 填充 ===

func applyPadding(data []byte, blockSize int, paddingType string) ([]byte, error) {
	switch paddingType {
	case PaddingPKCS7, PaddingPKCS5:
		return pkcs7Pad(data, blockSize), nil
	case PaddingISO10126:
		return iso10126Pad(data, blockSize)
	case PaddingANSIX923:
		return ansiX923Pad(data, blockSize)
	case PaddingZERO:
		return zeroPad(data, blockSize)
	case PaddingSPACE:
		return spacePad(data, blockSize)
	case PaddingNONE:
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
	case PaddingPKCS7, PaddingPKCS5:
		return pkcs7Unpad(data, blockSize)
	case PaddingISO10126:
		return iso10126Unpad(data)
	case PaddingANSIX923:
		return ansiX923Unpad(data, blockSize)
	case PaddingZERO:
		return zeroUnpad(data, blockSize)
	case PaddingSPACE:
		return spaceUnpad(data, blockSize)
	case PaddingNONE:
		return data, nil
	default:
		return nil, ErrUnsupportedPadding
	}
}

// PKCS7 (PKCS5相同)
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
	for _, v := range data[len(data)-padding:] {
		if int(v) != padding {
			return nil, ErrInvalidPadding
		}
	}
	return data[:len(data)-padding], nil
}

// ISO10126 填充实现优化
func iso10126Pad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		padding = blockSize
	}
	buf := make([]byte, len(data)+padding)
	copy(buf, data)
	if _, err := rand.Read(buf[len(data) : len(buf)-1]); err != nil {
		return nil, err
	}
	buf[len(buf)-1] = byte(padding)
	return buf, nil
}

func iso10126Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPadding
	}
	padding := int(data[len(data)-1])
	if padding < 1 || padding > len(data) {
		return nil, ErrInvalidPadding
	}
	return data[:len(data)-padding], nil
}

// ANSI X.923
func ansiX923Pad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		padding = blockSize
	}
	buf := make([]byte, len(data)+padding)
	copy(buf, data)
	// 中间填充0，最后一个字节填充长度
	// buf[len(data):len(buf)-1]默认为0
	buf[len(buf)-1] = byte(padding)
	return buf, nil
}

func ansiX923Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPadding
	}
	padding := int(data[len(data)-1])
	if padding < 1 || padding > blockSize || padding > len(data) {
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
	idx := bytes.LastIndexFunc(data, func(r rune) bool { return r != 0 })
	if idx == -1 {
		return []byte{}, nil
	}
	paddingLen := len(data) - idx - 1
	if paddingLen > blockSize {
		return nil, ErrPaddingSizeMismatch
	}
	for _, b := range data[idx+1:] {
		if b != 0 {
			return nil, ErrInvalidPadding
		}
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
	for ; idx >= 0 && data[idx] == ' '; idx-- {
	}
	padding := len(data) - idx - 1
	if padding > blockSize {
		return nil, ErrPaddingSizeMismatch
	}
	return data[:idx+1], nil
}

// === 辅助 ===

func validateParams(mode string, iv []byte, blockSize int) error {
	switch mode {
	case ModeCBC, ModeCTR, ModeOFB, ModeCFB:
		if len(iv) != blockSize {
			return fmt.Errorf("%w: expected %d bytes, got %d", ErrInvalidIVLength, blockSize, len(iv))
		}
	case ModeECB:
		if len(iv) != 0 {
			return fmt.Errorf("%w: ECB mode requires empty IV", ErrInvalidIVLength)
		}
	default:
		return ErrUnsupportedMode
	}
	return nil
}

func requiresPadding(mode string) bool {
	return mode == ModeCBC || mode == ModeECB
}
