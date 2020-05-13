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
	"net/http"

	"github.com/yhyzgn/gog"

	"github.com/yhyzgn/gox/common"
)

// Mapper 处理器映射器
type Mapper struct {
	contextPath string
	path        string
	ctrl        Controller
}

// NewMapper 创建映射器
func NewMapper(contextPath, path string, ctrl Controller) *Mapper {
	return &Mapper{
		contextPath: contextPath,
		path:        path,
		ctrl:        ctrl,
	}
}

// Request 注册一个新的处理器
func (mp *Mapper) Request(paths ...string) *Ship {
	if paths == nil || len(paths) == 0 {
		gog.Fatal("The param 'paths' can not be nil, must be '' at least.")
		return nil
	}
	return &Ship{
		contextPath: mp.contextPath,
		mapper:      mp,
		paths:       paths,
		methods:     make([]common.Method, 0),
		params:      make([]*common.Param, 0),
	}
}

// Get 注册一个 GET 请求的处理器
func (mp *Mapper) Get(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodGet)
}

// Head 注册一个 HEAD 请求的处理器
func (mp *Mapper) Head(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodHead)
}

// Post 注册一个 POST 请求的处理器
func (mp *Mapper) Post(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodPost)
}

// Put 注册一个 PUT 请求的处理器
func (mp *Mapper) Put(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodPut)
}

// Patch 注册一个 PATCH 请求的处理器
func (mp *Mapper) Patch(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodPatch)
}

// Delete 注册一个 DELETE 请求的处理器
func (mp *Mapper) Delete(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodDelete)
}

// Connect 注册一个 CONNECT 请求的处理器
func (mp *Mapper) Connect(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodConnect)
}

// Options 注册一个 OPTIONS 请求的处理器
func (mp *Mapper) Options(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodOptions)
}

// Trace 注册一个 TRACE 请求的处理器
func (mp *Mapper) Trace(paths ...string) *Ship {
	return mp.Request(paths...).Method(http.MethodTrace)
}
