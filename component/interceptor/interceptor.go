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
// time   : 2019-11-26 14:37
// version: 1.0.0
// desc   : 拦截器接口
//			分发器执行完后，请求将到达 拦截器
//			拦截器执行完后，请求将到达 处理器

package interceptor

import (
	"net/http"
	"reflect"

	"github.com/yhyzgn/gox/common"
)

// Interceptor 拦截器
type Interceptor interface {
	// PreHandle 请求处理前
	// 返回 true 将继续往下执行，返回 false 则截断请求
	PreHandle(writer http.ResponseWriter, request *http.Request, handler common.Handler) (bool, *http.Request, http.ResponseWriter)

	// 请求处理后
	AfterHandle(writer http.ResponseWriter, request *http.Request, handler common.Handler, result reflect.Value, err error) (*http.Request, http.ResponseWriter)
}
