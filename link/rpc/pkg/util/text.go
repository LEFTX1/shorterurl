package util

import (
	"strings"
)

// StrBuilder 字符串构建器，用于高效拼接字符串
type StrBuilder struct {
	builder strings.Builder
}

// Create 创建一个新的StrBuilder实例，可以选择性地提供初始字符串
func Create(initialStr ...string) *StrBuilder {
	sb := &StrBuilder{}
	if len(initialStr) > 0 {
		sb.builder.WriteString(initialStr[0])
	}
	return sb
}

// Append 追加字符串到构建器
func (sb *StrBuilder) Append(str string) *StrBuilder {
	sb.builder.WriteString(str)
	return sb
}

// String 返回构建的字符串
func (sb *StrBuilder) String() string {
	return sb.builder.String()
}

// Reset 重置构建器
func (sb *StrBuilder) Reset() *StrBuilder {
	sb.builder.Reset()
	return sb
}

// Len 返回当前构建的字符串长度
func (sb *StrBuilder) Len() int {
	return sb.builder.Len()
}
