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
// time   : 2019-11-29 9:35
// version: 1.0.0
// desc   : HTTP 异常

package common

import "net/http"

// HTTPError HTTP 异常定义
type HTTPError struct {
	Code  int    // 状态码
	Error string // 错误信息
}

// NewHTTPError 一个新的异常
func NewHTTPError(statusCode int, error string) *HTTPError {
	return &HTTPError{
		Code:  statusCode,
		Error: error,
	}
}

// Response 响应异常
func (err *HTTPError) Response(writer http.ResponseWriter) {
	http.Error(writer, err.Error, err.Code)
}
