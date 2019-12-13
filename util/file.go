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
// time   : 2019-12-13 11:02
// version: 1.0.0
// desc   : 文件工具类

package util

import "os"

// FileExist 判断文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}
