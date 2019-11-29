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
// time   : 2019-11-27 14:29
// version: 1.0.0
// desc   : 一个路由关系映射器

package core

import (
	"fmt"
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gox/wire"
	"net/http"
	"reflect"
	"strings"
)

// Ship 路由关系映射器
type Ship struct {
	mapper      *Mapper            // 所属的 处理器映射器
	path        string             // 配置的 path 路径
	handlerFunc common.HandlerFunc // 配置的 处理器
	methods     []common.Method    // http 请求方法列表
	params      []*common.Param    // 配置的参数列表
}

// 完成一条 处理器关系 映射
func (sp *Ship) Mapping() *Mapper {
	if sp.methods == nil {
		sp.methods = make([]common.Method, 0)
	}
	if len(sp.methods) == 0 {
		sp.methods = append(sp.methods, http.MethodGet)
	}

	v := reflect.ValueOf(sp.handlerFunc)
	if v.Kind() != reflect.Func {
		gog.FatalF("The handlerFunc %v must be function.", sp.handlerFunc)
	}

	// 检查参数有效性
	// 每个方法最多 2 个参数可以是 http.ResponseWriter 和 *http.Request
	// 其他均是自定义参数，需要注册
	x := v.Type()
	paramCount := x.NumIn()
	if paramCount > len(sp.params)+2 {
		// 有些参数未注册
		gog.Fatal("Maybe some params have not been registered.")
	}

	pos := 0
	for i := 0; i < paramCount; i++ {
		// 原始类型
		typ := x.In(i)
		// 具体类型，如果是指针，则变换为具体类型
		tp := typ
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
			gog.FatalF("Unsupported argument [%v] of function [%v]", typ, v)
		}

		if pos >= len(sp.params) {
			gog.Fatal("Maybe some params have not been registered.")
		}

		// 映射 Type
		sp.params[pos].Type = typ
		pos++
	}

	// 注册 每一条映射关系
	wire.Instance.Mapping(sp.resolvePath(), common.Handler(v), sp.methods, sp.params)
	return sp.mapper
}

// resolvePath 用 / 处理 path，构建标准 url path
func (sp *Ship) resolvePath() string {
	pref := sp.mapper.path
	path := sp.path

	if !strings.HasPrefix(pref, "/") {
		pref = "/" + pref
	}

	// 如果 handler 的 path 为 / ，则表示 controller 中的默认路径
	if path == "/" {
		path = ""
	} else if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return util.ReplaceAll(fmt.Sprintf("%s%s", pref, path), "//", "/")
}

// HandlerFunc 配置处理器
func (sp *Ship) HandlerFunc(handlerFunc common.HandlerFunc) *Ship {
	sp.handlerFunc = handlerFunc
	return sp
}

// Method 配置 http 请求方法
func (sp *Ship) Method(methods ...common.Method) *Ship {
	if methods != nil {
		sp.methods = append(sp.methods, methods...)
	}
	return sp
}

// Header 注册请求头参数
func (sp *Ship) Header(name string) *Ship {
	sp.params = append(sp.params, common.NewParam(name, true, true, false, false))
	return sp
}

// Required 注册普通参数，必需参数
func (sp *Ship) Required(name string) *Ship {
	sp.params = append(sp.params, common.NewParam(name, true, false, false, false))
	return sp
}

// Param 注册普通参数，可空
func (sp *Ship) Param(name string) *Ship {
	sp.params = append(sp.params, common.NewParam(name, false, false, false, false))
	return sp
}

// PathVariable 注册RESTFul格式参数
func (sp *Ship) PathVariable(name string) *Ship {
	sp.params = append(sp.params, common.NewParam(name, true, false, true, false))
	return sp
}

// Body 注册 RequestBody 参数
func (sp *Ship) Body(name string) *Ship {
	sp.params = append(sp.params, common.NewParam(name, true, false, false, true))
	return sp
}
