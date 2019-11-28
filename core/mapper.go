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
// time   : 2019-11-24 1:38 上午
// version: 1.0.0
// desc   : 

package core

import (
	"github.com/yhyzgn/gox/common"
	"net/http"
)

type Mapper struct {
	route *Route
}

func NewMapper(route *Route) *Mapper {
	return &Mapper{route: route}
}

func (mp *Mapper) Request(path string) *Ship {
	return &Ship{
		mapper:  mp,
		path:    path,
		methods: make([]common.Method, 0),
		params:  make([]*common.Param, 0),
	}
}

func (mp *Mapper) Get(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodGet)
	return hm
}

func (mp *Mapper) Head(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodHead)
	return hm
}

func (mp *Mapper) Post(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodPost)
	return hm
}

func (mp *Mapper) Put(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodPut)
	return hm
}

func (mp *Mapper) Patch(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodPatch)
	return hm
}

func (mp *Mapper) Delete(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodDelete)
	return hm
}

func (mp *Mapper) Connect(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodConnect)
	return hm
}

func (mp *Mapper) Options(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodOptions)
	return hm
}

func (mp *Mapper) Trace(path string) *Ship {
	hm := mp.Request(path)
	hm.methods = append(hm.methods, http.MethodTrace)
	return hm
}