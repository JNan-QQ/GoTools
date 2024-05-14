/**
  Copyright (c) [2024] [JiangNan]
  [go-tools] is licensed under Mulan PSL v2.
  You can use this software according to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
           http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
*/

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
