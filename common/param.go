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
// time   : 2019-11-24 4:15 下午
// version: 1.0.0
// desc   : 参数信息定义

package common

import "reflect"

// Param 参数信息定义
type Param struct {
	Name     string       // 参数名
	Required bool         // 是否必须
	InHeader bool         // 是否在 header 中，普通 header 参数
	InPath   bool         // 是否在 path 中，RESTful 参数
	IsBody   bool         // 是否在 body 中，RequestBody 参数
	RealType reflect.Type // 参数的实际类型
	IsPtr    bool         // 参数是否是指针
	ElemType reflect.Type // 如果实际类型是指针，这里记录指针所指向的类型
}

// NewParam 一个新的参数
func NewParam(name string, required, inHeader, inPath, isBody bool) *Param {
	return &Param{
		Name:     name,
		Required: required,
		InHeader: inHeader,
		InPath:   inPath,
		IsBody:   isBody,
	}
}
