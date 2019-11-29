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
// time   : 2019-11-24 9:01 下午
// version: 1.0.0
// desc   : 

package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/yhyzgn/gox/common"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// RecycleRequestBody 复用 request.Body
//
// 获取到本来的 request.Body
// 再把获取到的设置回去
func RecycleRequestBody(req *http.Request) []byte {
	if req != nil && req.Body != nil {
		bs, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
		return bs
	}
	return nil
}

// SetRequestAttribute 给 request 添加属性
func SetRequestAttribute(request *http.Request, key common.AttributeKey, value interface{}) *http.Request {
	return request.WithContext(context.WithValue(request.Context(), key, value))
}

// GetRequestAttribute 从 request 获取属性
func GetRequestAttribute(request *http.Request, key common.AttributeKey) interface{} {
	return request.Context().Value(key)
}

// SetRequestHeader 设置请求头
//
// 指定大小写
func SetRequestHeader(req *http.Request, key, value string) {
	req.Header[key] = []string{value}
}

// SetResponseHeader 设置响应头
//
// 指定大小写
func SetResponseHeader(res *http.Response, key, value string) {
	res.Header[key] = []string{value}
}

// SetResponseWriterHeader 设置响应头
//
// 指定大小写
func SetResponseWriterHeader(res http.ResponseWriter, key, value string) {
	res.Header()[key] = []string{value}
}

// AddURLQuery 向 URL 中添加 query 参数
//
// 添加 URL 参数
func AddURLQuery(url, key, value string) string {
	var sb strings.Builder
	sb.WriteString(url)

	if strings.Contains(url, "?") {
		// 如果不以 ? 结尾，也不以 & 结尾，就加上 & 连接符
		if !strings.HasSuffix(url, "?") && !strings.HasSuffix(url, "&") {
			sb.WriteString("&")
		}
	} else {
		sb.WriteString("?")
	}
	sb.WriteString(key)
	sb.WriteString("=")
	sb.WriteString(value)
	return sb.String()
}

// ResponseJSON 响应 json 数据
func ResponseJSON(writer http.ResponseWriter, value interface{}) {
	SetResponseWriterHeader(writer, "Content-Type", "application/json;charset=utf-8")
	if value != nil {
		bs, err := json.Marshal(value)
		if err == nil {
			err = ResponseBytes(writer, bs)
			if err == nil {
				return
			}
		}
		_ = ResponseBytes(writer, []byte(err.Error()))
		return
	}
	_ = ResponseBytes(writer, []byte("nil response"))
}

// ResponseBytes http 响应
func ResponseBytes(writer http.ResponseWriter, bytes []byte) (err error) {
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(bytes)
	return
}

// FormatHandlerArgs 格式化方法参数
func FormatRealArgsValue(args []reflect.Value) string {
	var sb strings.Builder
	sb.WriteString("(")
	if args != nil && len(args) > 0 {
		for i, arg := range args {
			if i > 0 {
				sb.WriteString(", ")
			}

			// net/http 下的类型，只打印 类型 即可
			if tp := arg.Type(); tp.PkgPath() == "net/http" || tp.Kind() == reflect.Ptr && tp.Elem().PkgPath() == "net/http" {
				sb.WriteString(arg.String())
				continue
			}
			sb.WriteString(fmt.Sprint(arg.Interface()))
		}
	}
	sb.WriteString(")")
	return sb.String()
}

// FormatHandlerArgs 格式化方法参数
func FormatHandlerArgs(params []*common.Param) string {
	var sb strings.Builder
	sb.WriteString("(")
	if params != nil && len(params) > 0 {
		for i, param := range params {
			name := param.Name
			if name == "" {
				name = fmt.Sprintf("arg%d", i)
			}
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(name)
			sb.WriteString(" ")
			if param.IsPtr {
				sb.WriteString("*")
				sb.WriteString(param.ElemType.Name())
			} else {
				sb.WriteString(param.RealType.Name())
			}
		}
	}
	sb.WriteString(")")
	return sb.String()
}
