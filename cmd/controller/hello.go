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
// time   : 2019-11-25 14:49
// version: 1.0.0
// desc   : 

package controller

import (
	"github.com/yhyzgn/gox/core"
	"strconv"
)

type HelloController struct {
}

func (h HelloController) Mapping(mapper *core.Mapper) {
	mapper.Request("/hello").HandlerFunc(h.Hello).Mapping()
	mapper.Request("/param").HandlerFunc(h.Param).Param("param").Mapping()
	mapper.Request("/rest/{id}").HandlerFunc(h.Rest).PathVariable("id").Mapping()
	mapper.Request("/json").HandlerFunc(h.JSON).Param("age").Mapping()
}

func (h HelloController) Hello() string {
	return "Hello GoX."
}

func (h HelloController) Param(param string) string {
	return "Hello " + param
}

func (h HelloController) Rest(id int) string {
	return "Hello " + strconv.Itoa(id)
}

func (h HelloController) JSON(age int) interface{} {
	return map[string]interface{}{
		"name": "张三",
		"age":  age,
	}
}
