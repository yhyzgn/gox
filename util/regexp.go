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
	"regexp"
	"strings"
)

const (
	// PatternRESTFul 匹配 RESTFul 格式的 pattern
	PatternRESTFul = "/{([\\w-_]+?)}"
	// PatternRESTFulReal 将配置的 RESTFul 格式 path 转换成 统配 path 的正则
	PatternRESTFulReal = "/[\\w-_]+?"
)

var (
	regRESTFul, _ = regexp.Compile(PatternRESTFul)
)

// IsRESTFul path 是否是 RESTFul 格式
func IsRESTFul(path string) bool {
	return regRESTFul.MatchString(path)
}

// GetRESTFulParams 获取 RESTFul 路径中的 所有参数
func GetRESTFulParams(path string) []string {
	matches := regRESTFul.FindAllStringSubmatch(path, -1)
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

// ConvertRESTFulPathToPattern 将配置的 RESTFul 格式 path 转换成 统配 path 的正则
func ConvertRESTFulPathToPattern(path string) string {
	return regRESTFul.ReplaceAllString(path, PatternRESTFulReal)
}

// GetPathVariableIndex 从注册的 path 中 获取 参数位置
func GetPathVariableIndex(name string, path string) int {
	nodes := strings.Split(path, "/")
	if nodes != nil {
		for i, node := range nodes {
			if regRESTFul.MatchString("/" + node) {
				return i
			}
		}
	}
	return -1
}
