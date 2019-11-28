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
// time   : 2019-11-24 1:54 上午
// version: 1.0.0
// desc   : http 请求方法集合，保证唯一性

package common

// MethodSet http 请求方法集合，保证唯一性
type MethodSet struct {
	flags   map[Method]bool
	methods []Method
}

// NewMethodSet 一个新的集合对象
func NewMethodSet() *MethodSet {
	return &MethodSet{
		flags:   make(map[Method]bool),
		methods: make([]Method, 0),
	}
}

// Has 是否已经存在 方法 method
func (ms *MethodSet) Has(method Method) bool {
	return ms.flags[method]
}

// Add 向集合中添加方法
func (ms *MethodSet) Add(method Method) *MethodSet {
	if !ms.Has(method) {
		ms.methods = append(ms.methods, method)
		ms.flags[method] = true
	}
	return ms
}

// Methods 获取所有方法
func (ms *MethodSet) Methods() []Method {
	return ms.methods
}
