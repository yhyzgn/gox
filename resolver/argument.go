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
// time   : 2019-11-25 17:21
// version: 1.0.0
// desc   : 参数处理器

package resolver

import (
	"github.com/yhyzgn/gox/common"
	"net/http"
	"reflect"
)

// ArgumentResolver 参数处理器
type ArgumentResolver interface {
	// Resolve 处理参数操作
	Resolve(args []reflect.Value, writer http.ResponseWriter, request *http.Request, handler common.Handler) []reflect.Value
}

// SimpleArgumentResolver 默认的参数处理器
type SimpleArgumentResolver struct {
}

// NewSimpleArgumentResolver 创建新的参数处理器对象
func NewSimpleArgumentResolver() *SimpleArgumentResolver {
	return new(SimpleArgumentResolver)
}

// Resolve 处理参数操作
func (sar *SimpleArgumentResolver) Resolve(args []reflect.Value, writer http.ResponseWriter, request *http.Request, handler common.Handler) []reflect.Value {
	// 这里接收到的参数列表已经过初步处理
	// 默认直接返回使用
	return args
}
