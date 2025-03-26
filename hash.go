package common

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash
//
//	@Description: 使用 bcrypt 对密码进行加密
//	@param password 用户密码
//	@return string 返回加密后的密码
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck
//
//	@Description: 对比明文密码和数据库的哈希值
//	@param password 用户密码
//	@param hash 数据库存放的用户数据
//	@return bool 返回是否正确
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MD5V
//
//	@Description: md5加密
//	@param str []byte
//	@param b byte
//	@return string
func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}
