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
	"net/http"
	"reflect"
	"strings"

	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/wire"
)

// Ship 路由关系映射器
type Ship struct {
	mapper      *Mapper            // 所属的 处理器映射器
	paths       []string           // 配置的 path 路径
	handlerFunc common.HandlerFunc // 配置的 处理器
	methods     []common.Method    // http 请求方法列表
	params      []*common.Param    // 配置的参数列表
}

// Mapping 完成一条 处理器关系 映射
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

	// 判断是否有 http.ResponseWriter 和 *http.Request 参数
	// 便于参数有效性的判断
	hasWriter, hasRequest := sp.hasResponseWriterAndRequest(x)
	delta := 0
	if hasWriter {
		delta++
	}
	if hasRequest {
		delta++
	}

	if paramCount > len(sp.params)+delta {
		// 有些参数未注册
		gog.Fatal("Maybe some params have not been registered.")
	} else if paramCount < len(sp.params)+delta {
		// 注册了一些不必要参数
		gog.Fatal("Maybe some params is needn't.")
	}

	// 记录所有参数
	tempParams := make([]*common.Param, paramCount)

	// pos 用来取 已注册过的 参数
	pos := 0
	for i := 0; i < paramCount; i++ {
		// 原始类型
		realType := x.In(i)
		// 具体类型，如果是指针，则变换为具体类型
		elemType := realType
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		elemKind := elemType.Kind()
		elemName := elemType.Name()

		var param *common.Param
		// http.ResponseWriter 或者 *http.Request
		if elemType.PkgPath() == "net/http" && (elemKind == reflect.Interface && elemName == "ResponseWriter" || realType.Kind() == reflect.Ptr && elemKind == reflect.Struct && elemName == "Request") {
			// http.ResponseWriter || *http.Request
			param = new(common.Param)
		} else {
			// 已注册过的参数 映射 Type
			param = sp.params[pos]
		}
		param.RealType = realType
		param.IsPtr = realType != elemType
		param.ElemType = elemType
		tempParams[i] = param
		pos++
	}
	sp.params = tempParams

	// 注册 每一条映射关系
	for _, path := range sp.resolvePath() {
		wire.Instance.Mapping(path, common.Handler(v), sp.methods, sp.params)
	}
	return sp.mapper
}

func (sp *Ship) hasResponseWriterAndRequest(tp reflect.Type) (hasWriter bool, hasRequest bool) {
	paramCount := tp.NumIn()
	for i := 0; i < paramCount; i++ {
		realType := tp.In(i)
		elemType := realType
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}

		if elemType.PkgPath() == "net/http" {
			if elemType.Kind() == reflect.Interface && elemType.Name() == "ResponseWriter" {
				hasWriter = true
				continue
			}

			if realType.Kind() == reflect.Ptr && elemType.Kind() == reflect.Struct && elemType.Name() == "Request" {
				hasRequest = true
			}
		}

		if hasWriter && hasRequest {
			return
		}
	}
	return
}

// resolvePath 用 / 处理 path，构建标准 url path
func (sp *Ship) resolvePath() []string {
	pref := sp.mapper.path

	if !strings.HasPrefix(pref, "/") {
		pref = "/" + pref
	}

	paths := make([]string, 0)
	for _, path := range sp.paths {
		// 如果 handler 的 path 为 / ，则表示 controller 中的默认路径
		if path == "/" {
			path = ""
		} else if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		paths = append(paths, strings.ReplaceAll(fmt.Sprintf("%s%s", pref, path), "//", "/"))
	}

	return paths
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

// PathVariable 注册RESTful格式参数
func (sp *Ship) PathVariable(name string) *Ship {
	sp.params = append(sp.params, common.NewParam(name, true, false, true, false))
	return sp
}

// Body 注册 RequestBody 参数
func (sp *Ship) Body(name string) *Ship {
	sp.params = append(sp.params, common.NewParam(name, true, false, false, true))
	return sp
}
