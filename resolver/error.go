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
// time   : 2020-05-15 21:34
// version: 1.0.0
// desc   : 异常处理器

package resolver

import (
	"github.com/yhyzgn/gox/util"
	"net/http"
)

// ErrorResolver 异常处理器
type ErrorResolver interface {

	// Resolve 处理异常
	Resolve(status int, err error, writer http.ResponseWriter) interface{}
}

// SimpleErrorResolver 默认的异常处理器
type SimpleErrorResolver struct{}

func NewSimpleErrorResolver() *SimpleErrorResolver {
	return new(SimpleErrorResolver)
}

func (ser *SimpleErrorResolver) Resolve(status int, err error, writer http.ResponseWriter) interface{} {
	util.ResponseJSONStatus(status, writer, err.Error())
	return nil
}
