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
// time   : 2019-11-29 15:32
// version: 1.0.0
// desc   : 

package controller

import (
	"fmt"
	"github.com/yhyzgn/gox/core"
	"github.com/yhyzgn/gox/util"
	"net/http"
)

type ParamController struct {
}

func (c ParamController) Mapping(mapper *core.Mapper) {
	mapper.Get("/native").HandlerFunc(c.Native).Mapping()
	mapper.Get("/noReturn").HandlerFunc(c.NoReturn).Mapping()
	mapper.Get("/required").HandlerFunc(c.Required).Required("name").Required("age").Mapping()
	mapper.Get("/param").HandlerFunc(c.Param).Required("name").Param("age").Mapping()
	mapper.Get("/header").HandlerFunc(c.Header).Header("token").Param("rand").Mapping()
	mapper.Get("/rest/{name}/{age}/test").HandlerFunc(c.Param).PathVariable("name").PathVariable("age").Mapping()
	mapper.Post("/body").HandlerFunc(c.Body).Body("user").Method(http.MethodPut).Mapping()
	mapper.Post("/vo").HandlerFunc(c.VO).Param("std").Mapping()
}

func (c ParamController) Native(writer http.ResponseWriter, request *http.Request) string {
	return c.res("Native " + request.URL.Query().Get("param"))
}

func (c ParamController) NoReturn(writer http.ResponseWriter, request *http.Request) {
	util.ResponseJSON(writer, c.res("NoReturn "+request.URL.Query().Get("param")))
}

func (c ParamController) Required(name string, age int) string {
	return c.res(fmt.Sprintf("Required name = %v, age = %d", name, age))
}

func (c ParamController) Param(name string, age int) string {
	return c.res(fmt.Sprintf("Param name = %v, age = %d", name, age))
}

func (c ParamController) Header(token string, rand int) string {
	return c.res(fmt.Sprintf("Header token = %v, rand = %d", token, rand))
}

func (c ParamController) Rest(name string, age int) string {
	return c.res(fmt.Sprintf("Rest name = %v, age = %d", name, age))
}

func (c ParamController) Body(user *User) *User {
	return user
}

func (c ParamController) VO(sdt *Student) *Student {
	return sdt
}

func (c ParamController) res(str string) string {
	return "GoX Param " + str
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Student struct {
	ID    int `param:"id"`
	Name  string
	Age   int    `param:"age" required:""`
	Token string `header:""`
}
