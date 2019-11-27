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
// desc   : 

package util

import (
	"regexp"
	"strings"
)

const (
	PatternRESTFul     = "/{([\\w-_]+?)}"
	PatternRESTFulReal = "/[\\w-_]+?"
)

var (
	regRESTFul, _ = regexp.Compile(PatternRESTFul)
)

func IsRESTFul(path string) bool {
	return regRESTFul.MatchString(path)
}

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

func ConvertRESTFulPathToPattern(path string) string {
	return regRESTFul.ReplaceAllString(path, PatternRESTFulReal)
}

// 从注册的 path 中 获取 参数位置
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
