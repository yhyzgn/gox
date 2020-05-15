// Copyright 2020 yhyzgn gox
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
// time   : 2020-05-13 10:12 下午
// version: 1.0.0
// desc   : 拦截器适配器

package adapter

import (
	"github.com/yhyzgn/gox/util"
	"net/http"
	"reflect"

	"github.com/yhyzgn/gox/common"
)

type Interceptor struct{}

// PreHandle 请求处理前
// 返回 true 将继续往下执行，返回 false 则截断请求
func (*Interceptor) PreHandle(writer http.ResponseWriter, request *http.Request, handler common.Handler) (bool, *http.Request, http.ResponseWriter) {
	return true, request, writer
}

// 请求处理后
func (*Interceptor) AfterHandle(writer http.ResponseWriter, request *http.Request, handler common.Handler, result reflect.Value, err error) (*http.Request, http.ResponseWriter) {
	return request, writer
}

// SetReqAttr 设置request字段
func (*Interceptor) SetReqAttr(req *http.Request, key common.AttributeKey, value interface{}) *http.Request {
	return util.SetRequestAttribute(req, key, value)
}

// GetReqAttr 获取request字段
func (*Interceptor) GetReqAttr(req *http.Request, key common.AttributeKey) interface{} {
	return util.GetRequestAttribute(req, key)
}
