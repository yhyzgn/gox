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
// time   : 2019-11-26 16:57
// version: 1.0.0
// desc   : 拦截器注册器

package interceptor

import (
	"github.com/yhyzgn/gog"
	"sync"
)

type item struct {
	path        string
	interceptor Interceptor
}

// Register 拦截器注册器
type Register struct {
	interceptors []item
	excludes     sync.Map
}

// NewRegister 新的注册器
func NewRegister() *Register {
	return &Register{
		interceptors: make([]item, 0),
	}
}

// AddInterceptors 添加拦截器
// 添加顺序 即 执行顺序
// path 匹配方式：
// 				/		->		所有请求
//				/xx		->		严格匹配
//				/xx/*	->		前缀匹配
func (ir *Register) AddInterceptors(path string, interceptors ...Interceptor) *Register {
	if path == "" || interceptors == nil || len(interceptors) == 0 {
		return ir
	}
	for _, ipt := range interceptors {
		ir.interceptors = append(ir.interceptors, item{
			path:        path,
			interceptor: ipt,
		})
	}
	gog.InfoF("The Interceptor [%v] registered.", path)
	return ir
}

// Exclude 添加排除路径
//
// 支持 前缀匹配 & 严格匹配
func (ir *Register) Exclude(path string) *Register {
	if _, ok := ir.excludes.Load(path); !ok {
		ir.excludes.Store(path, true)
	}
	return ir
}

// GetExcludes 获取那些被排除的路径
func (ir *Register) GetExcludes() map[string]bool {
	excludes := make(map[string]bool)
	ir.excludes.Range(func(key, value interface{}) bool {
		excludes[key.(string)] = value.(bool)
		return true
	})
	return excludes
}

// Iterate 遍历所有拦截器，并执行相应回到操作
func (ir *Register) Iterate(iterator func(index int, path string, interceptor Interceptor) (skip, passed bool)) (bool, string) {
	if iterator != nil {
		var (
			skip   bool
			passed bool
		)
		for i, item := range ir.interceptors {
			skip, passed = iterator(i, item.path, item.interceptor)
			if skip {
				// 可能 path 不匹配，跳过当前拦截器
				continue
			}
			if !passed {
				// 拦截器不通过
				return false, item.path
			}
		}
	}
	return true, ""
}

// ReverseIterate 逆序遍历所有拦截器，并执行相应回到操作
func (ir *Register) ReverseIterate(iterator func(index int, path string, interceptor Interceptor)) {
	if iterator != nil {
		for i := len(ir.interceptors) - 1; i >= 0; i-- {
			iterator(i, ir.interceptors[i].path, ir.interceptors[i].interceptor)
		}
	}
}
