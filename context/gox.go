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
// desc   : 

package context

import (
	"fmt"
	"github.com/yhyzgn/ghost/config"
	"net/http"
	"sync"
)

type GoXContext struct {
	reader            *config.Reader
	wares             map[string]interface{}
	onceMap           map[string]bool
	errorHandlers     map[int]http.HandlerFunc
	NotFound          http.HandlerFunc
	UnsupportedMethod http.HandlerFunc
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

func Current() *GoXContext {
	return current
}

func (c *GoXContext) Read(filename string) (data []byte, errs error) {
	return c.reader.Read(filename)
}

func (c *GoXContext) SetWare(name string, ware interface{}) WareContext {
	if !c.onceMap[name] {
		c.wares[name] = ware
	}
	return c
}

func (c *GoXContext) SetWareOnce(name string, ware interface{}) WareContext {
	if c.wares[name] == nil && !c.onceMap[name] {
		c.wares[name] = ware
		c.onceMap[name] = true
	}
	return c
}

func (c *GoXContext) GetWare(name string) interface{} {
	return c.wares[name]
}

func (c *GoXContext) AddErrorHandler(statusCode int, handler http.HandlerFunc) *GoXContext {
	c.errorHandlers[statusCode] = handler
	return c
}

func (c *GoXContext) GetErrorHandler(statusCode int) http.HandlerFunc {
	return c.errorHandlers[statusCode]
}
