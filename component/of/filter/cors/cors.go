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
// time   : 2019-12-03 11:14
// version: 1.0.0
// desc   : 跨域拦截器

package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/yhyzgn/gox/util"

	"github.com/yhyzgn/gox/component/filter"
)

// XCorsFilter 跨域拦截器
type XCorsFilter struct {
	origins    []string // 授权的源控制
	methods    []string // 允许请求的 HTTP Method
	headers    []string // 控制哪些 header 能发送真正的请求
	exposes    []string // 那些允许暴露的 header
	credential bool     // 控制是否开启与 ajax 的 Cookie 提交方式
	maxAge     int64    // 授权的时间
}

// NewXCorsFilter 创建新拦截器
func NewXCorsFilter() *XCorsFilter {
	return &XCorsFilter{
		origins:    make([]string, 0),
		methods:    make([]string, 0),
		headers:    make([]string, 0),
		exposes:    make([]string, 0),
		credential: false,
		maxAge:     3600,
	}
}

// DoFilter 执行跨域拦截器
func (c *XCorsFilter) DoFilter(writer http.ResponseWriter, request *http.Request, chain *filter.Chain) {
	// 支持跨域
	if len(c.origins) == 0 {
		c.origins = append(c.origins, "*")
	}
	if len(c.methods) == 0 {
		c.methods = append(c.methods, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead)
	}

	// 设置到响应头
	util.SetResponseWriterHeader(writer, "Access-Control-Allow-Origin", strings.Join(c.origins, ", "))
	util.SetResponseWriterHeader(writer, "Access-Control-Allow-Methods", strings.Join(c.methods, ", "))

	if len(c.headers) > 0 {
		util.SetResponseWriterHeader(writer, "Access-Control-Allow-Headers", strings.Join(c.headers, ", "))
	}
	if len(c.exposes) > 0 {
		util.SetResponseWriterHeader(writer, "Access-Control-Expose-Headers", strings.Join(c.exposes, ", "))
	}
	util.SetResponseWriterHeader(writer, "Access-Control-Allow-Credentials", strconv.FormatBool(c.credential))
	util.SetResponseWriterHeader(writer, "Access-Control-Max-Age", strconv.FormatInt(c.maxAge, 10))

	// 屏蔽 OPTIONS 请求
	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	// 继续往下执行
	chain.DoFilter(writer, request)
}

// AllowedOrigins 配置授权源控制
func (c *XCorsFilter) AllowedOrigins(origins ...string) *XCorsFilter {
	c.origins = append(c.origins, origins...)
	return c
}

// AllowedMethods 允许的请求方法
//
// 允许请求头中携带
func (c *XCorsFilter) AllowedMethods(methods ...string) *XCorsFilter {
	c.methods = append(c.methods, methods...)
	return c
}

// AllowedHeaders 允许的自定义 Header
//
// 允许访问携带的该属性
func (c *XCorsFilter) AllowedHeaders(headers ...string) *XCorsFilter {
	c.headers = append(c.headers, headers...)
	return c
}

// ExposedHeaders 允许暴露的 Header
func (c *XCorsFilter) ExposedHeaders(exposes ...string) *XCorsFilter {
	c.exposes = append(c.exposes, exposes...)
	return c
}

// AllowCredential 是否支持 cookie 上传
func (c *XCorsFilter) AllowCredential(credential bool) *XCorsFilter {
	c.credential = credential
	return c
}

// MaxAge 授权时间
func (c *XCorsFilter) MaxAge(maxAge int64) *XCorsFilter {
	c.maxAge = maxAge
	return c
}
