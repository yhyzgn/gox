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
	"net/http"
)

type NormalController struct {
}

func (c NormalController) Mapping(mapper *core.Mapper) {
	mapper.
		Request("/").HandlerFunc(c.Default).Mapping().
		Get("/get").HandlerFunc(c.Get).Mapping().
		Post("/post").HandlerFunc(c.Post).Mapping().
		Delete("/delete").HandlerFunc(c.Delete).Mapping().
		Request("/multi").HandlerFunc(c.Multi).Method(http.MethodGet, http.MethodPost).Mapping()
}

func (c NormalController) Default() string {
	return c.res("Default")
}

func (c NormalController) Get() string {
	return c.res("Get")
}

func (c NormalController) Post() string {
	return c.res("Post")
}

func (c NormalController) Delete() string {
	return c.res("Delete")
}

func (c NormalController) Multi() string {
	return c.res("Multi")
}

func (c NormalController) res(str string) string {
	return "GoX Normal " + str
}
