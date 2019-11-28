// Copyright 2019 yhyzgn xgo
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

package context

import (
	"fmt"
	"github.com/yhyzgn/ghost/config"
	"net/http"
	"sync"
)

// GoXContext GoX 上下文
type GoXContext struct {
	reader            *config.Reader           // 资源读取器
	wares             map[string]interface{}   // 一些组件
	onceMap           map[string]bool          // 一次性组件
	errorHandlers     map[int]http.HandlerFunc // 错误处理器，每个错误码对应一个处理器
	NotFound          http.HandlerFunc         // 404错误处理器
	UnsupportedMethod http.HandlerFunc         // 方法不支持错误处理器
}

var (
	once    sync.Once
	current *GoXContext
)

func init() {
	once.Do(func() {
		current = &GoXContext{
			reader:        config.NewReader(),
			wares:         make(map[string]interface{}),
			onceMap:       make(map[string]bool),
			errorHandlers: make(map[int]http.HandlerFunc),
			NotFound:      http.NotFound,
			UnsupportedMethod: func(writer http.ResponseWriter, request *http.Request) {
				http.Error(writer, fmt.Sprintf("Unsupported http method [%v].", request.Method), http.StatusMethodNotAllowed)
			},
		}
	})
}

// Current 获取当前上下文对象
func Current() *GoXContext {
	return current
}

// Read 读取资源文件
func (c *GoXContext) Read(filename string) (data []byte, errs error) {
	return c.reader.Read(filename)
}

// SetWare 设置组件
func (c *GoXContext) SetWare(name string, ware interface{}) WareContext {
	if !c.onceMap[name] {
		c.wares[name] = ware
	}
	return c
}

// SetWareOnce 设置一次性组件，修改无效
func (c *GoXContext) SetWareOnce(name string, ware interface{}) WareContext {
	if c.wares[name] == nil && !c.onceMap[name] {
		c.wares[name] = ware
		c.onceMap[name] = true
	}
	return c
}

// GetWare 获取组件
func (c *GoXContext) GetWare(name string) interface{} {
	return c.wares[name]
}

// AddErrorHandler 添加错误码处理器
func (c *GoXContext) AddErrorHandler(statusCode int, handler http.HandlerFunc) *GoXContext {
	c.errorHandlers[statusCode] = handler
	return c
}

// GetErrorHandler 获取错误码处理器
func (c *GoXContext) GetErrorHandler(statusCode int) http.HandlerFunc {
	return c.errorHandlers[statusCode]
}
