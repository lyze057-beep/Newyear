package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func MD5(str string) string {

	//方法一
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return md5str1
}
func Salt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", errors.New("生成盐失败:" + err.Error())
	}
	return hex.EncodeToString(salt), nil
}

//	func EncryptPassword(rawPassword string) (string, error) {
//		// 空密码直接返回错误
//		if rawPassword == "" {
//			return "", errors.New("密码不能为空")
//		}
//
//		// GenerateFromPassword：生成哈希密码
//		// bcrypt.DefaultCost：默认加密强度（可调整，值越大越安全但越慢）
//		hashBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
//		if err != nil {
//			return "", errors.New("密码加密失败：" + err.Error())
//		}
//
//		// 字节转字符串，方便存储到数据库
//		return string(hashBytes), nil
//	}
//
// // VerifyPassword 验证密码是否正确
// // 参数：原始密码 / 数据库中存储的哈希密码
// // 返回：是否验证通过
//
//	func VerifyPassword(rawPassword, hashedPassword string) bool {
//		// CompareHashAndPassword：对比哈希密码和原始密码
//		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
//		return err == nil // 无错误则密码正确
//	}
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 密码比对
func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil

}
