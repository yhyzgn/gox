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
	"fmt"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gox/wire"
	"github.com/yhyzgn/gog"
	"net/http"
	"reflect"
	"strings"
)

type Mapper struct {
	route *Route
}

type mapping struct {
	mapper      *Mapper
	path        string
	handlerFunc common.HandlerFunc
	methods     []common.Method
	params      []*common.Param
}

func NewMapper(route *Route) *Mapper {
	return &Mapper{route: route}
}

func (mp *mapping) Mapping() *Mapper {
	if mp.methods == nil {
		mp.methods = make([]common.Method, 0)
	}
	if len(mp.methods) == 0 {
		mp.methods = append(mp.methods, http.MethodGet)
	}

	v := reflect.ValueOf(mp.handlerFunc)
	if v.Kind() != reflect.Func {
		gog.FatalF("The handlerFunc %v must be function.", mp.handlerFunc)
	}

	// 检查参数有效性
	// 每个方法最多 2 个参数可以是 http.ResponseWriter 和 *http.Request
	// 其他均是自定义参数，需要注册
	x := v.Type()
	paramCount := x.NumIn()
	if paramCount > len(mp.params)+2 {
		// 有些参数未注册
		gog.Fatal("Maybe some params have not been registered.")
	}

	pos := 0
	for i := 0; i < paramCount; i++ {
		tp := x.In(i)
		if tp.Kind() == reflect.Ptr {
			tp = tp.Elem()
		}
		pkg := tp.PkgPath()
		kind := tp.Kind()
		name := tp.Name()

		if pkg == "net/http" {
			// 可能是 http.ResponseWriter 或者 *http.Request
			if kind == reflect.Interface && name == "ResponseWriter" {
				// http.ResponseWriter
				continue
			}

			if kind == reflect.Struct && name == "Request" {
				// http.Request
				continue
			}
			gog.FatalF("Unsupported argument [%v] of function [%v]", tp, v)
		}

		if pos >= len(mp.params) {
			gog.Fatal("Maybe some params have not been registered.")
		}

		// 映射 Type
		mp.params[pos].Type = tp
		pos++
	}

	wire.Instance.Mapping(mp.resolvePath(), common.Handler(v), mp.methods, mp.params)
	return mp.mapper
}

func (mp *mapping) resolvePath() string {
	pref := mp.mapper.route.path
	path := mp.path

	if !strings.HasPrefix(pref, "/") {
		pref = "/" + pref
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return util.ReplaceAll(fmt.Sprintf("%s%s", pref, path), "//", "/")
}

func (mp *Mapper) Path(path string) *mapping {
	mpg := new(mapping)
	mpg.mapper = mp
	mpg.path = path
	mpg.methods = make([]common.Method, 0)
	mpg.params = make([]*common.Param, 0)
	return mpg
}

func (mp *mapping) HandlerFunc(handlerFunc common.HandlerFunc) *mapping {
	mp.handlerFunc = handlerFunc
	return mp
}

func (mp *mapping) Method(methods ...common.Method) *mapping {
	if methods != nil {
		mp.methods = append(mp.methods, methods...)
	}
	return mp
}

func (mp *mapping) Header(name string) *mapping {
	mp.params = append(mp.params, common.NewParam(name, true, true, false, false))
	return mp
}
func (mp *mapping) Param(name string) *mapping {
	mp.params = append(mp.params, common.NewParam(name, true, false, false, false))
	return mp
}

func (mp *mapping) ParamNil(name string) *mapping {
	mp.params = append(mp.params, common.NewParam(name, false, false, false, false))
	return mp
}

func (mp *mapping) PathVariable(name string) *mapping {
	mp.params = append(mp.params, common.NewParam(name, true, false, true, false))
	return mp
}

func (mp *mapping) Body(name string) *mapping {
	mp.params = append(mp.params, common.NewParam(name, true, false, false, true))
	return mp
}
