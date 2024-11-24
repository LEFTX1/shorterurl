package common

import "strings"

// MaskPhone 对手机号进行脱敏处理，保留前3位和后4位
func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone // 如果手机号长度不足，直接返回原始数据
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskEmail 对邮箱进行脱敏处理，仅保留邮箱名的首字母和域名
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // 如果邮箱格式无效，直接返回原始数据
	}

	name := parts[0]
	domain := parts[1]

	// 脱敏邮箱名
	maskedName := name[:1] + "****"
	return maskedName + "@" + domain
}

// MaskIDCard 对身份证号进行脱敏处理，保留前4位和后4位
func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard // 如果身份证号长度不足，直接返回原始数据
	}
	return idCard[:4] + "********" + idCard[len(idCard)-4:]
}

// MaskGeneric 对任意字符串进行脱敏处理，保留指定的前后位数，其余用 * 替代
func MaskGeneric(data string, prefixLen, suffixLen int) string {
	if len(data) <= prefixLen+suffixLen {
		return data // 如果字符串长度不足以脱敏，直接返回原始数据
	}

	// 构造脱敏字符串
	middle := strings.Repeat("*", len(data)-prefixLen-suffixLen)
	return data[:prefixLen] + middle + data[len(data)-suffixLen:]
}
