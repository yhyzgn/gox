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
// desc   : 

package common

type MethodSet struct {
	flags   map[Method]bool
	methods []Method
}

func NewMethodSet() *MethodSet {
	return &MethodSet{
		flags:   make(map[Method]bool),
		methods: make([]Method, 0),
	}
}

func (ms *MethodSet) Has(method Method) bool {
	return ms.flags[method]
}

func (ms *MethodSet) Add(method Method) *MethodSet {
	if !ms.Has(method) {
		ms.methods = append(ms.methods, method)
		ms.flags[method] = true
	}
	return ms
}

func (ms *MethodSet) Methods() []Method {
	return ms.methods
}
