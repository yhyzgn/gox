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
// time   : 2019-11-24 1:05 上午
// version: 1.0.0
// desc   : 

package gox

import (
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/core"
	"net/http"
	"regexp"
	"testing"
)

type A struct {
}

func (a A) Mapping(mapper *core.Mapper) {
	mapper.Request("/hello").HandlerFunc(a.Hello).Method(http.MethodGet, http.MethodPost).Param("name").Param("age").Mapping()
}

func (A) Hello(name string, age int, writer http.ResponseWriter, request *http.Request) string {
	return "hello"
}

func TestRouter_Add(t *testing.T) {
	reg, _ := regexp.Compile("/\\*+?$")
	gog.Debug(reg.MatchString("/api/**"))

	gog.Debug(reg.ReplaceAllString("/api/**", "/.+?"))

	//fmt.Println(util.StringToInt("100", 0))
	//
	//pattern := "/{([\\w-_]+?)}"
	//test := "/api/{name}/test/{age}/{user-id}/{auth_code}"
	//
	//fmt.Println(regexp.MatchString(pattern, test))
	//
	//reg, _ := regexp.Compile(pattern)
	//matches := reg.FindAllStringSubmatch(test, -1)
	//fmt.Println(matches)
	//fmt.Println(reg.FindAllStringSubmatchIndex(test, -1))
	//
	//for _, i := range matches {
	//	fmt.Println(i[1])
	//}
	//
	//real := reg.ReplaceAllString(test, "/[\\w-_]+?")
	//fmt.Println(real)
	//
	//fmt.Println(regexp.MatchString("^" + real + "$", "/api/Jason/test/23/1234/asdfasfasfasf"))


	//NewRouter().Mapping("/api", &A{})

	//a := &A{}

	//value := reflect.ValueOf(a.Hello)
	//
	//x := value.Type()
	//fmt.Println(x.String())

	//numIn := x.NumIn() //Count inbound parameters
	//numOut := x.NumOut() //Count outbounding parameters

	//fmt.Println("Method:", x.String())
	//fmt.Println("Variadic:", x.IsVariadic()) // Used (<type> ...) ?
	//fmt.Println("Package:", x.PkgPath())

	//for i := 0; i < numIn; i++ {
	//
	//	inV := x.In(i)
	//	if inV.Kind() == reflect.Ptr {
	//		inV = inV.Elem()
	//	}
	//	fmt.Println(inV.PkgPath())
	//	fmt.Println(inV.Kind())
	//	fmt.Println(inV.Name())
	//
	//	//if inV.Kind() == reflect.Ptr {
	//	//	inV = inV.Elem()
	//	//}
	//	//in_Kind := inV.Kind() //func
	//	//fmt.Printf("\nParameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\n-----------", in_Kind, inV.Name())
	//}
	//for o := 0; o < numOut; o++ {
	//
	//	returnV := x.Out(0)
	//	return_Kind := returnV.Kind()
	//	fmt.Printf("\nParameter OUT: "+strconv.Itoa(o)+"\nKind: %v\nName: %v\n", return_Kind, returnV.Name())
	//}
}
