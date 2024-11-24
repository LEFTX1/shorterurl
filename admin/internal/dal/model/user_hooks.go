package model

import (
	"fmt"
	"go-zero-shorterurl/admin/internal/common"
	"gorm.io/gorm"
)

// BeforeSave GORM 钩子，在保存前加密手机号
func (u *TUser) BeforeSave(tx *gorm.DB) (err error) {
	if u.Phone != "" {
		// 打印原始手机号
		fmt.Printf("BeforeSave - 原始手机号: %s\n", u.Phone)

		encrypted, err := common.Encrypt(u.Phone)
		if err != nil {
			fmt.Printf("BeforeSave - 加密失败: %v\n", err)
			return err
		}

		// 打印加密后的手机号
		fmt.Printf("BeforeSave - 加密后手机号: %s\n", encrypted)
		u.Phone = encrypted
	}
	return nil
}

// AfterFind GORM 钩子，在查询后解密手机号
func (u *TUser) AfterFind(tx *gorm.DB) (err error) {
	if u.Phone != "" {
		// 打印从数据库中查询到的加密手机号
		fmt.Printf("AfterFind - 数据库中存储的加密手机号: %s\n", u.Phone)

		decrypted, err := common.Decrypt(u.Phone)
		if err != nil {
			fmt.Printf("AfterFind - 解密失败: %v\n", err)
			return err
		}

		// 打印解密后的手机号
		fmt.Printf("AfterFind - 解密后手机号: %s\n", decrypted)
		u.Phone = decrypted
	}
	return nil
}
