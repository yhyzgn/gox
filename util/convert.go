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
// time   : 2019-11-24 11:17 下午
// version: 1.0.0
// desc   : 类型转换工具

package util

import (
	"strconv"
)

// StringToInt String 转为不同长度 int
func StringToInt(value string, size int) int64 {
	it, err := strconv.ParseInt(value, 10, size)
	if err == nil {
		return it
	}
	return 0
}

// StringToUInt String 转为不同长度的 uint
func StringToUInt(value string, size int) uint64 {
	it, err := strconv.ParseUint(value, 10, size)
	if err == nil {
		return it
	}
	return 0
}
