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
	"errors"
	"fmt"
	"github.com/yhyzgn/gox/core"
	"github.com/yhyzgn/gox/util"
	"net/http"
	"testing"
)

type A struct {
}

func (a A) Mapping(mapper *core.Mapper) {
	mapper.Request("/normal").HandlerFunc(a.Normal).Required("id").Method(http.MethodGet, http.MethodPost).Mapping()
	mapper.Request("/hello").HandlerFunc(a.Hello).Method(http.MethodGet, http.MethodPost).Required("name").Required("age").Mapping()
}

func (A) Normal(id int) error {
	return errors.New("Normal")
}

func (A) Hello(name string, age int) string {
	return fmt.Sprintf("hello %s %d", name, age)
}

func TestRouter_Add(t *testing.T) {
	server := NewGoX()

	server.Mapping("/api", new(A))

	server.NotFoundHandler(func(writer http.ResponseWriter, request *http.Request) {
		util.ResponseJSONStatus(http.StatusNotFound, writer, "你的请求被绑架了")
	})

	server.Run(&http.Server{
		Addr: fmt.Sprintf(":%d", 8888),
	})

	//wtr := http.ResponseWriter(&of.ResponseWriter{})
	//req := &http.Request{}
	//
	//tpReq := reflect.ValueOf(req)
	//tpWtr := reflect.ValueOf(wtr)
	//
	//fmt.Println(tpReq.Type().Elem().PkgPath())
	//fmt.Println(tpReq.Type().Kind())
	//fmt.Println(tpReq.Type().Elem().Kind())
	//fmt.Println(tpReq.Type().Elem().Name())
	//
	//fmt.Println(tpWtr.Type().Elem().PkgPath())
	//fmt.Println(tpWtr.Type().Kind())
	//fmt.Println(tpWtr.Type().Elem().Kind())
	//fmt.Println(tpWtr.Type().Elem().Name())

	//var val []*A
	//tp := reflect.TypeOf(val)
	//gog.Debug(tp)
	//gog.Debug(tp.Kind() == reflect.Slice)
	//gog.Debug(tp.Elem())

	//gog.Debug(util.FirstToLower("AA"))

	//reg, _ := regexp.Compile("/\\*+?$")
	//gog.Debug(reg.MatchString("/api/**"))
	//
	//gog.Debug(reg.ReplaceAllString("/api/**", "/.+?"))

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
