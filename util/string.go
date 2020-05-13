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
	"reflect"
	"strconv"
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

// StringToValue 将字符串转换为其他类型
func StringToValue(kind reflect.Kind, value string) reflect.Value {
	var arg reflect.Value
	switch kind {
	case reflect.String:
		arg = reflect.ValueOf(value)
		break
	case reflect.Int:
		it, err := strconv.Atoi(value)
		if err == nil {
			arg = reflect.ValueOf(it)
		} else {
			arg = reflect.ValueOf(0)
		}
		break
	case reflect.Int8:
		arg = reflect.ValueOf(int8(StringToInt(value, 8)))
		break
	case reflect.Int16:
		arg = reflect.ValueOf(int16(StringToInt(value, 16)))
		break
	case reflect.Int32:
		arg = reflect.ValueOf(int32(StringToInt(value, 32)))
		break
	case reflect.Int64:
		arg = reflect.ValueOf(StringToInt(value, 64))
		break
	case reflect.Uint:
		arg = reflect.ValueOf(uint(StringToUInt(value, 0)))
		break
	case reflect.Uint8:
		arg = reflect.ValueOf(uint8(StringToUInt(value, 8)))
		break
	case reflect.Uint16:
		arg = reflect.ValueOf(uint16(StringToUInt(value, 16)))
		break
	case reflect.Uint32:
		arg = reflect.ValueOf(uint32(StringToUInt(value, 32)))
		break
	case reflect.Uint64:
		arg = reflect.ValueOf(StringToUInt(value, 64))
		break
	case reflect.Float32:
		arg = reflect.ValueOf(StringToFloat(value, 32))
		break
	case reflect.Float64:
		arg = reflect.ValueOf(StringToFloat(value, 64))
		break
	case reflect.Bool:
		bl, err := strconv.ParseBool(value)
		if err == nil {
			arg = reflect.ValueOf(bl)
		} else {
			arg = reflect.ValueOf(false)
		}
		break
	}
	return arg
}

// FirstToUpper 首字母大写
func FirstToUpper(src string) string {
	if src == "" {
		return src
	}
	chars := []rune(src)
	if chars[0] >= 97 && chars[0] <= 122 {
		chars[0] -= 32
	}
	return string(chars)
}

// FirstToLower 首字母小写
func FirstToLower(src string) string {
	if src == "" {
		return src
	}
	chars := []rune(src)
	if chars[0] >= 65 && chars[0] <= 90 {
		chars[0] += 32
	}
	return string(chars)
}

// StartWith 字符串开头匹配
func StartWith(text string, start string) bool {
	return strings.HasPrefix(text, start)
}

// EndWith 字符串结尾匹配
func EndWith(text string, end string) bool {
	return strings.HasSuffix(text, end)
}

// Trim 字符串去除两端空白
func Trim(str string) string {
	return strings.TrimSpace(str)
}
