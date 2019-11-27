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
// time   : 2019-11-25 17:21
// version: 1.0.0
// desc   : 参数处理器

package resolver

import (
	"encoding/json"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gox/wire"
	"github.com/yhyzgn/gog"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type ArgumentResolver interface {
	Resolve(hw *wire.HandlerWire, writer http.ResponseWriter, request *http.Request, isRESTFul bool) []reflect.Value
}

type SimpleArgumentResolver struct {
}

func NewSimpleArgumentResolver() *SimpleArgumentResolver {
	return new(SimpleArgumentResolver)
}

func (sar SimpleArgumentResolver) Resolve(hw *wire.HandlerWire, writer http.ResponseWriter, request *http.Request, isRESTFul bool) []reflect.Value {
	path := request.URL.Path
	handler := reflect.Value(hw.Handler)
	// 每个方法最多 2 个参数可以是 http.ResponseWriter 和 *http.Request
	// 其他均是自定义参数，需要注册
	x := handler.Type()
	paramCount := x.NumIn()
	pc := handler.Pointer()
	handlerName := util.ReplaceAll(runtime.FuncForPC(pc).Name(), "-fm", "(...)")

	var pathVariables []string
	if isRESTFul {
		pathVariables = util.GetRESTFulParams(hw.Path)
	}
	args := make([]reflect.Value, 0)

	if hw.Params != nil && len(hw.Params) > 0 {
		// 已经注册过参数，这里就需要获取参数
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
					args = append(args, reflect.ValueOf(writer))
					continue
				}

				if kind == reflect.Struct && name == "Request" {
					// http.Request
					args = append(args, reflect.ValueOf(request))
					continue
				}
			}

			param := hw.Params[pos]
			pos++

			// RESTful 格式传参
			if param.InPath && pathVariables != nil {
				index := util.GetPathVariableIndex(param.Name, hw.Path)
				if index > -1 {
					// 找到啦~
					// 找到位置后，去 url path 的对应位置上获取参数值即可
					temp := getPathVariableValue(index, path)
					// 添加到参数列表
					args = append(args, stringToValue(param, temp))
				} else {
					gog.ErrorF("The path [%v] has not contains path variable [%v].", hw.Path, param.Name)
				}
				continue
			}

			// 从请求头获取
			if param.InHeader {
				temp := request.Header.Get(param.Name)
				if temp == "" && param.Required {
					gog.ErrorF("Maybe the param [%v] in request header defected.", param.Name)
					continue
				}
				// 添加到参数列表
				args = append(args, stringToValue(param, temp))
				continue
			}

			// 有 requestBody 参数
			// 仅支持 POST 方法
			if param.IsBody {
				if request.Method != http.MethodPost && request.Method != http.MethodPut {
					gog.ErrorF("RequestBody only support 'POST' and 'PUT' method, but now is [%v].", request.Method)
					continue
				}

				if !VerifyMethod(hw, http.MethodPost) && !VerifyMethod(hw, http.MethodPut) {
					gog.ErrorF("Maybe the handler [%v] should be register as 'POST' or 'PUT' method, now is %v.", handlerName, hw.Methods)
					continue
				}

				// 获取到 requestBody
				bs := util.RecycleRequestBody(request)
				if bs != nil {
					var arg interface{}
					switch kind {
					case reflect.Map:
						arg = reflect.MakeMap(tp).Interface()
						break
					case reflect.Slice:
						arg = reflect.MakeSlice(tp, 0, 0).Interface()
						break
					case reflect.Struct:
						arg = reflect.New(tp).Interface()
						break
					}
					err := json.Unmarshal(bs, &arg)
					if err != nil {
						gog.Error(err)
						continue
					}
					// 添加到参数列表
					args = append(args, reflect.ValueOf(arg))
				}
				continue
			}

			// 普通参数
			// 先从 URL 中获取
			temp := request.URL.Query().Get(param.Name)
			if temp == "" {
				// 从 form 中获取
				temp = request.FormValue(param.Name)
				if temp == "" && request.Method == http.MethodPost {
					temp = request.PostFormValue(param.Name)
				}
			}
			// 添加到参数列表
			args = append(args, stringToValue(param, temp))
		}
	} else {
		// 未注册，可能会有 http.ResponseWriter 和 *http.Request
		// 其他参数先不管啦~
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
					args = append(args, reflect.ValueOf(writer))
					continue
				}

				if kind == reflect.Struct && name == "Request" {
					// http.Request
					args = append(args, reflect.ValueOf(request))
					continue
				}
			}
		}
	}

	return args
}

// 用 参数位置 从 实际 path 中获取到参数值
func getPathVariableValue(index int, path string) string {
	nodes := strings.Split(path, "/")
	if nodes != nil && index < len(nodes) {
		return nodes[index]
	}
	return ""
}

func stringToValue(param *common.Param, value string) reflect.Value {
	var arg reflect.Value
	switch param.Type.Kind() {
	case reflect.String:
		arg = reflect.ValueOf(value)
		break
	case reflect.Int:
		it, err := strconv.Atoi(value)
		if err == nil {
			arg = reflect.ValueOf(it)
		} else {
			arg = reflect.ValueOf(0)
		}
		break
	case reflect.Int8:
		arg = reflect.ValueOf(int8(util.StringToInt(value, 8)))
		break
	case reflect.Int16:
		arg = reflect.ValueOf(int16(util.StringToInt(value, 16)))
		break
	case reflect.Int32:
		arg = reflect.ValueOf(int32(util.StringToInt(value, 32)))
		break
	case reflect.Int64:
		arg = reflect.ValueOf(util.StringToInt(value, 64))
		break
	case reflect.Uint:
		arg = reflect.ValueOf(uint(util.StringToUInt(value, 0)))
		break
	case reflect.Uint8:
		arg = reflect.ValueOf(uint8(util.StringToUInt(value, 8)))
		break
	case reflect.Uint16:
		arg = reflect.ValueOf(uint16(util.StringToUInt(value, 16)))
		break
	case reflect.Uint32:
		arg = reflect.ValueOf(uint32(util.StringToUInt(value, 32)))
		break
	case reflect.Uint64:
		arg = reflect.ValueOf(util.StringToUInt(value, 64))
		break
	case reflect.Bool:
		bl, err := strconv.ParseBool(value)
		if err == nil {
			arg = reflect.ValueOf(bl)
		} else {
			arg = reflect.ValueOf(false)
		}
		break
	}
	return arg
}

func VerifyMethod(hw *wire.HandlerWire, method string) bool {
	for _, md := range hw.Methods {
		if string(md) == method {
			return true
		}
	}
	return false
}
