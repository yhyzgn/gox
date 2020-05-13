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
// time   : 2019-10-11 1:49
// version: 1.0.0
// desc   :

package util

import (
	"fmt"
	"reflect"
	"testing"
)

type RfTest struct {
	Info string
	desc string
}

func Test(t *testing.T) {
	str := "测试看看"
	fmt.Println(IsPtr(reflect.TypeOf(str)))
	fmt.Println(IsPtr(reflect.TypeOf(&str)))

	rf := RfTest{}
	fmt.Println(IsStruct(reflect.TypeOf(rf)))
	fmt.Println(IsStruct(reflect.TypeOf(&rf)))

	fmt.Println(ElemType(str))
	fmt.Println(ElemType(&str))
	fmt.Println(ElemType(rf))
	fmt.Println(ElemType(&rf))

	fmt.Println(StructType(str))
	fmt.Println(StructType(&str))
	fmt.Println(StructType(rf))
	fmt.Println(StructType(&rf))
}

func TestFieldSet(t *testing.T) {
	rf := RfTest{Info: "aa", desc: "desc"}
	fmt.Println(StructType(rf))

	set(&rf)
}

func set(bean interface{}) {
	fmt.Println(IsPtr(reflect.TypeOf(bean)))
	val := reflect.ValueOf(bean).Elem()
	fmt.Println(val.NumField())

	field := val.Field(1)
	FieldSet(field, reflect.ValueOf("呵呵"))
	fmt.Println(bean)
}
