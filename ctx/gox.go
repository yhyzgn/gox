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
// time   : 2019-11-27 9:41
// version: 1.0.0
// desc   : GoX 上下文

package ctx

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/yhyzgn/gox/resolver"

	"github.com/yhyzgn/gox/resource"
)

// GoXContext GoX 上下文
type GoXContext struct {
	contextPath       string                    // 根路径
	reader            *resource.Reader          // 资源读取器
	errorHandlers     sync.Map                  // 错误处理器，每个错误码对应一个处理器
	staticDir         string                    // 静态资源文件夹路径
	notFound          http.HandlerFunc          // 404错误处理器
	unSupportedMethod http.HandlerFunc          // 方法不支持错误处理器
	argumentResolver  resolver.ArgumentResolver // 参数处理器
	resultResolver    resolver.ResultResolver   // 结果处理器
	errorResolver     resolver.ErrorResolver    // 全局异常处理器
}

var (
	once    sync.Once
	current *GoXContext
)

func init() {
	once.Do(func() {
		current = &GoXContext{
			reader:   resource.NewReader(),
			notFound: http.NotFound,
			unSupportedMethod: func(writer http.ResponseWriter, request *http.Request) {
				http.Error(writer, fmt.Sprintf("Unsupported http method [%v].", request.Method), http.StatusMethodNotAllowed)
			},
			argumentResolver: resolver.NewSimpleArgumentResolver(),
			resultResolver:   resolver.NewSimpleResultResolver(),
			errorResolver:    resolver.NewSimpleErrorResolver(),
		}
	})
}

// C 获取当前上下文对象
func C() *GoXContext {
	return current
}

// Read 读取资源文件
func (c *GoXContext) Read(filename string) (data []byte, errs error) {
	return c.reader.Read(filename)
}

// Load 加载资源文件到实例
func (c *GoXContext) Load(filename string, bean interface{}) (err error) {
	return c.reader.Load(filename, bean)
}

// SetContextPath 设置根路径
func (c *GoXContext) SetContextPath(contextPath string) *GoXContext {
	c.contextPath = contextPath
	return c
}

// GetContextPath 获取根路径
func (c *GoXContext) GetContextPath() string {
	return c.contextPath
}

// SetStaticDir 设置静态资源文件夹
func (c *GoXContext) SetStaticDir(dir string) *GoXContext {
	c.staticDir = dir
	return c
}

// SetNotFoundHandler 设置 404 错误处理器
func (c *GoXContext) SetNotFoundHandler(handler http.HandlerFunc) *GoXContext {
	c.notFound = handler
	return c
}

// SetUnSupportMethodHandler 设置 请求方法 错误处理器
func (c *GoXContext) SetUnSupportMethodHandler(handler http.HandlerFunc) *GoXContext {
	c.unSupportedMethod = handler
	return c
}

// AddErrorHandler 添加错误码处理器
func (c *GoXContext) AddErrorHandler(statusCode int, handler http.HandlerFunc) *GoXContext {
	c.errorHandlers.Store(statusCode, handler)
	return c
}

// SetArgumentResolver 设置参数处理器
func (c *GoXContext) SetArgumentResolver(resolver resolver.ArgumentResolver) *GoXContext {
	c.argumentResolver = resolver
	return c
}

// SetResultResolver 设置结果处理器
func (c *GoXContext) SetResultResolver(resolver resolver.ResultResolver) *GoXContext {
	c.resultResolver = resolver
	return c
}

// SetErrorResolver 设置全局异常处理器
func (c *GoXContext) SetErrorResolver(resolver resolver.ErrorResolver) *GoXContext {
	c.errorResolver = resolver
	return c
}

// GetArgumentResolver 获取参数处理器
func (c *GoXContext) GetArgumentResolver() resolver.ArgumentResolver {
	return c.argumentResolver
}

// GetResultResolver 获取结果处理器
func (c *GoXContext) GetResultResolver() resolver.ResultResolver {
	return c.resultResolver
}

// GetErrorResolver 获取全局异常处理器
func (c *GoXContext) GetErrorResolver() resolver.ErrorResolver {
	return c.errorResolver
}

// GetStaticDir 获取静态资源文件夹
func (c *GoXContext) GetStaticDir() string {
	return c.staticDir
}

// GetNotFoundHandler 获取 404 错误处理器
func (c *GoXContext) GetNotFoundHandler() http.HandlerFunc {
	return c.notFound
}

// GetUnSupportMethodHandler 获取 请求方法 错误处理器
func (c *GoXContext) GetUnSupportMethodHandler() http.HandlerFunc {
	return c.notFound
}

// GetErrorHandler 获取错误码处理器
func (c *GoXContext) GetErrorHandler(statusCode int) http.HandlerFunc {
	handler, ok := c.errorHandlers.Load(statusCode)
	if !ok {
		return nil
	}
	return handler.(http.HandlerFunc)
}
