package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"shorterurl/admin/internal/config"
	"sync"
)

var (
	// keyHolder 用于存储 AES 加密密钥
	keyHolder []byte
	once      sync.Once
)

// Init 初始化加密密钥
// 从配置中读取 Base64 编码的密钥，解码后存储到 keyHolder
func InitAES(cfg config.Config) error {
	var err error
	once.Do(func() {
		// 解码 Base64 编码的密钥
		keyHolder, err = base64.StdEncoding.DecodeString(cfg.Crypto.AESKey)
		if err != nil {
			err = errors.New("无法解码 AES 密钥")
			return
		}

		// 检查密钥长度是否为 32 字节 (AES-256)
		if len(keyHolder) != 32 {
			err = errors.New("无效的 AES 密钥长度，必须为 32 字节 (AES-256)")
		}
	})
	return err
}

// Encrypt 对明文进行 AES-256 加密，返回 Base64 编码的密文
func Encrypt(plaintext string) (string, error) {
	if keyHolder == nil {
		return "", errors.New("加密密钥未初始化")
	}

	// 创建 AES 加密块
	block, err := aes.NewCipher(keyHolder)
	if err != nil {
		return "", err
	}

	// 使用 GCM 模式 (Galois/Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce，长度由 GCM 决定
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 使用 Seal 方法加密数据，返回密文
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// 返回 Base64 编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 对 Base64 编码的密文进行 AES-256 解密，返回明文
func Decrypt(ciphertext string) (string, error) {
	if keyHolder == nil {
		return "", errors.New("解密密钥未初始化")
	}

	// 解码 Base64 编码的密文
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建 AES 解密块
	block, err := aes.NewCipher(keyHolder)
	if err != nil {
		return "", err
	}

	// 使用 GCM 模式 (Galois/Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 检查密文长度是否足够包含 nonce
	if len(data) < aesGCM.NonceSize() {
		return "", errors.New("密文长度不足")
	}

	// 分离 nonce 和真正的密文
	nonce, ciphertextData := data[:aesGCM.NonceSize()], data[aesGCM.NonceSize():]

	// 使用 Open 方法解密数据，返回明文
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
