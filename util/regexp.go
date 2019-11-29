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
// time   : 2019-11-25 17:28
// version: 1.0.0
// desc   : 正则工具

package util

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// PatternRESTful 匹配 RESTful 格式的 pattern
	PatternRESTful = "/{([\\w-_]+?)}"
	// PatternRESTfulReal 将配置的 RESTful 格式 path 转换成 统配 path 的正则
	PatternRESTfulReal = "/[\\w-_]+?"
)

var (
	regRESTful, _ = regexp.Compile(PatternRESTful)
)

// IsRESTful path 是否是 RESTful 格式
func IsRESTful(path string) bool {
	return regRESTful.MatchString(path)
}

// GetRESTfulParams 获取 RESTful 路径中的 所有参数
func GetRESTfulParams(path string) []string {
	matches := regRESTful.FindAllStringSubmatch(path, -1)
	// str : /api/{name}/test/{age}/{user-id}/{auth_code}
	// matches : [[/{name} name] [/{age} age] [/{user-id} user-id] [/{auth_code} auth_code]]
	result := make([]string, 0)
	if matches != nil && len(matches) > 0 {
		for _, item := range matches {
			result = append(result, item[1])
		}
	}
	return result
}

// ConvertRESTfulPathToPattern 将配置的 RESTful 格式 path 转换成 统配 path 的正则
func ConvertRESTfulPathToPattern(path string) string {
	return regRESTful.ReplaceAllString(path, PatternRESTfulReal)
}

// GetPathVariableIndex 从注册的 path 中 获取 参数位置
func GetPathVariableIndex(name string, path string) int {
	nodes := strings.Split(path, "/")
	if nodes != nil {
		for i, node := range nodes {
			// 根据 RESTful 分段来匹配 path 节点
			if regRESTful.MatchString("/"+node) && node == fmt.Sprintf("{%v}", name) {
				return i
			}
		}
	}
	return -1
}
