// Copyright 2019 yhyzgn gox
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2019-11-24 7:14 下午
// version: 1.0.0
// desc   : 字符串工具

package util

import (
	"bytes"
	"strings"
)

// ReplaceAll 替换
func ReplaceAll(str, old, new string) string {
	return string(bytes.ReplaceAll([]byte(str), []byte(old), []byte(new)))
}

// Replace 替换
func Replace(str, old, new string, n int) string {
	return string(bytes.Replace([]byte(str), []byte(old), []byte(new), n))
}

// FillPrefix 前缀填充补齐
func FillPrefix(src, ch string, ln int) string {
	if len(src) >= ln {
		return src
	}
	delta := ln - len(src)

	var sb strings.Builder
	for i := 0; i < delta; i++ {
		sb.WriteString(ch)
	}
	sb.WriteString(src)
	return sb.String()
}

// FillSuffix 后缀填充补齐
func FillSuffix(src, ch string, ln int) string {
	if len(src) >= ln {
		return src
	}
	delta := ln - len(src)

	var sb strings.Builder
	sb.WriteString(src)
	for i := 0; i < delta; i++ {
		sb.WriteString(ch)
	}
	return sb.String()
}
