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
// desc   : 结果处理器

package resolver

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gox/wire"
)

// ResultResolver 结果处理器
type ResultResolver interface {
	// Resolve 处理结果集
	Resolve(hw *wire.HandlerWire, values []reflect.Value, writer http.ResponseWriter, request *http.Request) (reflect.Value, error)

	// Response 响应结果
	Response(value reflect.Value, writer http.ResponseWriter)
}

// SimpleResultResolver 默认的结果处理器
type SimpleResultResolver struct {
}

// NewSimpleResultResolver 创建新的结果处理器对象
func NewSimpleResultResolver() *SimpleResultResolver {
	return new(SimpleResultResolver)
}

// Resolve 处理结果集
// 只接受最多两个返回值的结果集
// 如果不满足需求，可自定义
func (srr *SimpleResultResolver) Resolve(hw *wire.HandlerWire, values []reflect.Value, writer http.ResponseWriter, request *http.Request) (value reflect.Value, err error) {
	path := request.URL.Path
	handler := reflect.Value(hw.Handler)
	pc := handler.Pointer()
	handlerName := strings.ReplaceAll(runtime.FuncForPC(pc).Name(), "-fm", "(...)")

	if values == nil || len(values) == 0 {
		// 没有返回值，无需处理
		return
	}

	ln := len(values)
	// 只有1个返回值，必定是 请求响应结果
	if ln == 1 {
		//srr.Response(values[0], writer)
		value = values[0]
		return
	}

	// 结果1：请求响应结果
	// 结果2：错误信息
	if ln == 2 {
		if e := values[1]; e.Interface() != nil {
			err = e.Interface().(error)
			return
		}
		value = values[0]
		return
	}

	// 结果不能超过2个
	err = fmt.Errorf("the path [%v] handled [%v] support 2 results at most, but now is [%d]", path, handlerName, ln)
	return
}

// Response 响应结果
func (srr *SimpleResultResolver) Response(value reflect.Value, writer http.ResponseWriter) {
	util.ResponseJSON(writer, value.Interface())
}
