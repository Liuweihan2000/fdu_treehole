package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5WithSalt(str string) string {
	target := str + viper.GetString("email_hash") + "@fudan.edu.cn"
	h := md5.New()
	h.Write([]byte(target))
	return hex.EncodeToString(h.Sum(nil))
}
