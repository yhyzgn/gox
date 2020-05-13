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
// time   : 2020-05-13 9:30 下午
// version: 1.0.0
// desc   : 自定义响应器

package internal

import (
	"bytes"
	"net/http"
)

// 自定义响应器
type ResponseWriter struct {
	wt     http.ResponseWriter
	buf    *bytes.Buffer
	status int
}

// 创建自定义响应器
func NewResponseWriter(writer http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		wt:  writer,
		buf: bytes.NewBuffer(nil),
	}
}

// 设置实际响应器
func (rw *ResponseWriter) SetWriter(writer http.ResponseWriter) {
	rw.wt = writer
}

// 响应数据长度
func (rw *ResponseWriter) ContentLength() int {
	return len(rw.ResponseBody())
}

// 响应体
func (rw *ResponseWriter) ResponseBody() []byte {
	return rw.buf.Bytes()
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.wt.Header()
}

func (rw *ResponseWriter) Write(bs []byte) (int, error) {
	_, _ = rw.buf.Write(bs)
	return rw.wt.Write(bs)
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.wt.WriteHeader(statusCode)
}
