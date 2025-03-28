package hash

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
	"strings"
)

const (
	// Base62字符集
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// 短链接长度
	shortLinkLength = 6
)

// HashToBase62 将原始URL哈希为Base62编码的短链接
// 返回长度为6的短链接字符串
func HashToBase62(originUrl string) string {
	// 计算MD5哈希值
	hasher := md5.New()
	hasher.Write([]byte(originUrl))
	md5Hex := hex.EncodeToString(hasher.Sum(nil))

	// 取前8个字符转为int64
	hexInt := new(big.Int)
	hexInt.SetString(md5Hex[:8], 16)

	// 转为Base62编码
	var result strings.Builder
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for hexInt.Cmp(zero) > 0 && result.Len() < shortLinkLength {
		hexInt.DivMod(hexInt, base, mod)
		result.WriteByte(base62Chars[mod.Int64()])
	}

	// 如果长度不足，则用0补齐
	for result.Len() < shortLinkLength {
		result.WriteByte(base62Chars[0])
	}

	return result.String()
}
