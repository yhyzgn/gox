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
// time   : 2019-11-24 1:38 上午
// version: 1.0.0
// desc   : 处理器映射器

package core

import (
	"github.com/yhyzgn/gox/common"
	"net/http"
)

// Mapper 处理器映射器
type Mapper struct {
	path string
	ctrl Controller
}

// NewMapper 创建映射器
func NewMapper(path string, ctrl Controller) *Mapper {
	return &Mapper{
		path: path,
		ctrl: ctrl,
	}
}

// Request 注册一个新的处理器
func (mp *Mapper) Request(path string) *Ship {
	return &Ship{
		mapper:  mp,
		path:    path,
		methods: make([]common.Method, 0),
		params:  make([]*common.Param, 0),
	}
}

// Get 注册一个 GET 请求的处理器
func (mp *Mapper) Get(path string) *Ship {
	return mp.Request(path).Method(http.MethodGet)
}

// Head 注册一个 HEAD 请求的处理器
func (mp *Mapper) Head(path string) *Ship {
	return mp.Request(path).Method(http.MethodHead)
}

// Post 注册一个 POST 请求的处理器
func (mp *Mapper) Post(path string) *Ship {
	return mp.Request(path).Method(http.MethodPost)
}

// Put 注册一个 PUT 请求的处理器
func (mp *Mapper) Put(path string) *Ship {
	return mp.Request(path).Method(http.MethodPut)
}

// Patch 注册一个 PATCH 请求的处理器
func (mp *Mapper) Patch(path string) *Ship {
	return mp.Request(path).Method(http.MethodPatch)
}

// Delete 注册一个 DELETE 请求的处理器
func (mp *Mapper) Delete(path string) *Ship {
	return mp.Request(path).Method(http.MethodDelete)
}

// Connect 注册一个 CONNECT 请求的处理器
func (mp *Mapper) Connect(path string) *Ship {
	return mp.Request(path).Method(http.MethodConnect)
}

// Options 注册一个 OPTIONS 请求的处理器
func (mp *Mapper) Options(path string) *Ship {
	return mp.Request(path).Method(http.MethodOptions)
}

// Trace 注册一个 TRACE 请求的处理器
func (mp *Mapper) Trace(path string) *Ship {
	return mp.Request(path).Method(http.MethodTrace)
}
