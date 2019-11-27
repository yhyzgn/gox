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
// time   : 2019-11-24 12:51 上午
// version: 1.0.0
// desc   : 

package core

import (
	"github.com/yhyzgn/gox/common"
)

type Route struct {
	path    string
	methods *common.MethodSet
	ctrl    Controller
	handler common.Handler
}

func NewRoute() *Route {
	r := new(Route)
	r.methods = common.NewMethodSet()
	return r
}

func (r *Route) Path(path string) *Route {
	r.path = path
	return r
}

func (r *Route) Method(methods ...common.Method) *Route {
	if methods != nil && len(methods) > 0 {
		for _, method := range methods {
			r.methods.Add(method)
		}
	}
	return r
}

func (r *Route) Controller(ctrl Controller) *Route {
	r.ctrl = ctrl
	return r
}

func (r *Route) Handler(handler common.Handler) *Route {
	r.handler = handler
	return r
}
