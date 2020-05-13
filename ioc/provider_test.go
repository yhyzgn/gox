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
// time   : 2019-10-11 2:53
// version: 1.0.0
// desc   :

package ioc

import (
	"fmt"
	"github.com/yhyzgn/gox/util"
	"testing"
)

type A struct {
	Info string
}

type B struct {
	Info string
}

type Demo struct {
	A  *A `auto:""`
	AA *A `auto:"" scope:"single"`
	B  *B `auto:"b" scope:"prototype"`
	BB *B `auto:"b" scope:"prototype"`
}

func TestNewProvider(t *testing.T) {
	provider := NewProvider()

	provider.
		Put("a", func() (instance interface{}) {
			return &A{Info: "AA"}
		}).
		Put("b", func() (instance interface{}) {
			return &B{Info: "BB"}
		})

	fmt.Println(provider.String())

	demo := &Demo{}
	// 注册时自动注入
	provider.Single("", demo)

	fmt.Println(provider.String())

	// 打印指针，确保单例和实例的指针地址
	fmt.Printf("a: %p\naa: %p\nb: %p\nbb: %p\n", demo.A, demo.AA, demo.B, demo.BB)

	fmt.Println(util.StructType(demo))
}
