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
// time   : 2019-10-11 1:41
// version: 1.0.0
// desc   : 工具库

package util

import (
	"reflect"
	"unsafe"
)

// IsStruct 是否是结构体
func IsStruct(tp reflect.Type) bool {
	return tp.Kind() == reflect.Struct
}

// IsPtr 是否是指针类型
func IsPtr(tp reflect.Type) bool {
	return tp.Kind() == reflect.Ptr
}

// ElemType 获取某个变量的实际类型，自动忽略指针
func ElemType(bean interface{}) reflect.Type {
	tp := reflect.TypeOf(bean)
	if IsPtr(tp) {
		tp = tp.Elem()
	}
	return tp
}

// StructType 获取某个变量的实际类型，并判断是否是结构体，是否是指针
func StructType(bean interface{}) (tp reflect.Type, isStruct bool, isPtr bool) {
	tp = reflect.TypeOf(bean)
	if IsPtr(tp) {
		tp = tp.Elem()
		isPtr = true
	}
	if tp.Kind() == reflect.Struct {
		isStruct = true
	}
	return
}

// FieldSet 通过反射给结构体字段设置值
func FieldSet(field reflect.Value, value reflect.Value) {
	if !field.CanSet() {
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	}
	field.Set(value)
}
