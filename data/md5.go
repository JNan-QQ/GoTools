package data

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5 代价密字符串
type MD5 string

// Encrypt 默认加密 32位小写
func (m MD5) Encrypt() MD5 {
	h := md5.New()
	h.Write([]byte(m))
	return MD5(hex.EncodeToString(h.Sum(nil)))
}

func (m MD5) IsBig() MD5 {
	return MD5(strings.ToUpper(string(m)))
}

func (m MD5) IsShort() MD5 {
	if len(m) != 32 {
		return m
	}
	return MD5(string(m)[8:24])
}

func (m MD5) ToString() string {
	return string(m)
}
