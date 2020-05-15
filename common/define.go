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
// time   : 2019-11-24 12:47 上午
// version: 1.0.0
// desc   : 定义一些东西

package common

import (
	"mime/multipart"
	"reflect"
)

// 定义一些常量
const (
	FilterChainName         = "gox-filter-chain"         // 过滤器注册链名字
	InterceptorRegisterName = "gox-interceptor-register" // 拦截器注册器名字
	RequestDispatcherName   = "gox-request-dispatcher"   // 请求分发器名字
	ArgumentResolverName    = "gox-argument-resolver"    // 参数处理器名字
	ResultResolverName      = "gox-result-resolver"      // 结果处理器名字
	ErrorResolverName       = "gox-error-resolver"       // 全局异常处理器名字
	RequestFilterIndexName  = "gox-filter-index"         // 每个请求过滤器索引名字
)

// AttributeKey request 属性的键类型
type AttributeKey string

// Method http 请求方法
type Method string

// HandlerFunc 控制器中的方法
// 只能是 Func 类型
type HandlerFunc interface{}

// Handler 请求处理器，指向 HandlerFunc
type Handler reflect.Value

// MultipartFile 表单上传文件模型
type MultipartFile struct {
	Header *multipart.FileHeader
	File   multipart.File
}

// Call 处理器实际处理操作，调用具体方法，并返回值
func (h Handler) Call(args []reflect.Value) []reflect.Value {
	return h.Get().Call(args)
}

// Get 处理器实际类型
func (h Handler) Get() reflect.Value {
	return reflect.Value(h)
}
