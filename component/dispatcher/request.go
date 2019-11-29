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
// time   : 2019-11-24 3:00 上午
// version: 1.0.0
// desc   : 请求分发器-实现类

package dispatcher

import (
	"encoding/json"
	"fmt"
	"github.com/yhyzgn/ghost/utils"
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/component/interceptor"
	"github.com/yhyzgn/gox/context"
	"github.com/yhyzgn/gox/resolver"
	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gox/wire"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// RequestDispatcher 请求分发器-实现类
type RequestDispatcher struct {
	register *interceptor.Register
}

// NewRequestDispatcher 创建新的分发器
func NewRequestDispatcher() *RequestDispatcher {
	return new(RequestDispatcher)
}

// SetInterceptorRegister 配置拦截器注册器
// 用于请求分发后的拦截操作
func (rd *RequestDispatcher) SetInterceptorRegister(register *interceptor.Register) {
	rd.register = register
}

// Dispatch 分发具体请求
func (rd *RequestDispatcher) Dispatch(writer http.ResponseWriter, request *http.Request) {
	reqPath := request.URL.Path

	// 如果请求路径以 / 结尾，则自动去除
	if strings.HasSuffix(reqPath, "/") {
		reqPath = reqPath[0 : len(reqPath)-1]
	}

	// 匹配路由
	for _, h := range wire.Instance.All() {
		// 如果直接完全匹配，说明不是 RESTful 模式
		if reqPath == h.Path {
			rd.doDispatch(h, writer, request, false)
			return
		} else if util.IsRESTful(h.Path) {
			// 否则 正则匹配
			// 将 路由注册的路径 转换为 正则匹配模板，再看是否与真实路径匹配
			realPathPattern := util.ConvertRESTfulPathToPattern(h.Path)
			matched, err := regexp.MatchString(realPathPattern, reqPath)
			if err == nil && matched {
				// RESTful 匹配上了
				rd.doDispatch(h, writer, request, true)
				return
			}
		}
	}

	// 匹配不到，就只能 404 啦~
	context.Current().NotFound(writer, request)
}

// doDispatch 具体的请求分发操作
func (rd *RequestDispatcher) doDispatch(hw *wire.HandlerWire, writer http.ResponseWriter, request *http.Request, isRESTful bool) {
	md := VerifyMethod(hw, request.Method)
	if !md {
		// 不支持的 http 方法
		context.Current().UnsupportedMethod(writer, request)
		return
	}

	// 处理器
	handler := hw.Handler

	// 参数处理器
	argumentResolver := util.GetWare(common.ArgumentResolverName, resolver.NewSimpleArgumentResolver()).(resolver.ArgumentResolver)
	// 结果处理器
	resultResolver := util.GetWare(common.ResultResolverName, resolver.NewSimpleResultResolver()).(resolver.ResultResolver)

	// 先处理一遍参数
	args, ex := rd.resolve(hw, writer, request, isRESTful)
	if ex != nil {
		gog.Error(ex.Error)
		ex.Response(writer)
		return
	}

	// 再调用参数处理器处理
	argumentResolver.Resolve(args, writer, request, handler)

	gog.InfoF("Params of request path [{}] are [{}], matched router [{}] of params [{}]", request.URL.Path, util.FormatRealArgsValue(args), hw.Path, util.FormatHandlerArgs(hw.Params))

	// 处理前，执行拦截器 PreHandle() 方法
	if rd.register != nil {
		pass, path := rd.register.Iterate(func(index int, path string, interceptor interceptor.Interceptor) (skip, passed bool) {
			// 匹配 path，未匹配到的直接跳过
			defer func() {
				if skip {
					gog.InfoF("The request [%v] has skipped by interceptor [%v].", request.URL.Path, path)
				} else if passed {
					gog.InfoF("The request [%v] has passed by interceptor [%v].", request.URL.Path, path)
				} else {
					gog.InfoF("The request [%v] has been intercepted by interceptor [%v].", request.URL.Path, path)
				}
			}()
			if path == "/" {
				// 所有请求
				return false, interceptor.PreHandle(writer, request, handler)
			} else if reg, e := regexp.Compile("/\\*+$"); e == nil && reg.MatchString(path) {
				// 前缀匹配
				pattern := reg.ReplaceAllString(path, "/.+?")
				if matched, err := regexp.MatchString("^"+pattern+"$", request.URL.Path); matched && err == nil {
					// 前缀匹配成功，执行拦截器
					return false, interceptor.PreHandle(writer, request, handler)
				}
				// 匹配不成功的直接跳过
				return true, true
			} else if path == request.URL.Path {
				// 严格匹配，只有路径完全相同才走过滤器
				return false, interceptor.PreHandle(writer, request, handler)
			} else {
				// 跳过
				return true, true
			}
		})

		// 拦截器不通过
		if !pass {
			gog.InfoF("The request [%v] has been intercepted by interceptor [%v].", request.URL.Path, path)
			return
		}
	}

	// 拦截器通过后，将请求交由 处理器 处理
	// 已经获取到参数列表，执行方法即可
	results := handler.Call(args)
	if results == nil || len(results) == 0 {
		// 无返回值
		gog.InfoF("The request [{}] responded, and the handler needn't return any value.", request.URL.Path)
		return
	}
	// 响应结果交由 结果处理器 处理
	res, err := resultResolver.Resolve(hw, results, writer, request)
	// 如果有错误，就响应错误信息
	if err != nil {
		res = reflect.ValueOf(err)
	}

	// 处理完成后，执行拦截器的 AfterHandle() 方法
	if rd.register != nil {
		rd.register.ReverseIterate(func(index int, path string, interceptor interceptor.Interceptor) {
			// 匹配 path，未匹配到的直接跳过
			if path == "/" {
				// 所有请求
				interceptor.AfterHandle(writer, request, handler, res, err)
			} else if reg, e := regexp.Compile("/\\*+$"); e == nil && reg.MatchString(path) {
				// 前缀匹配
				pattern := reg.ReplaceAllString(path, "/.+?")
				if matched, err := regexp.MatchString("^"+pattern+"$", request.URL.Path); matched && err == nil {
					// 前缀匹配成功，执行拦截器
					interceptor.AfterHandle(writer, request, handler, res, err)
				}
			} else if path == request.URL.Path {
				// 严格匹配，只有路径完全相同才走过滤器
				interceptor.AfterHandle(writer, request, handler, res, err)
			}
		})
	}

	// 拦截器通过后，响应处理结果
	resultResolver.Response(res, writer)
}

// resolve 初步处理参数
func (rd *RequestDispatcher) resolve(hw *wire.HandlerWire, writer http.ResponseWriter, request *http.Request, isRESTful bool) ([]reflect.Value, *common.HTTPError) {
	path := request.URL.Path
	handler := reflect.Value(hw.Handler)
	handlerName := util.ReplaceAll(runtime.FuncForPC(handler.Pointer()).Name(), "-fm", util.FormatHandlerArgs(hw.Params))

	var pathVariables []string
	if isRESTful {
		pathVariables = util.GetRESTfulParams(hw.Path)
	}
	args := make([]reflect.Value, 0)

	for _, param := range hw.Params {
		// ----------------------------------------------------------------------------------------------    net/http    ----------------------------------------------------------------------------------------------
		// http.ResponseWriter || *http.Request
		if param.ElemType.PkgPath() == "net/http" {
			if param.ElemType.Kind() == reflect.Interface && param.ElemType.Name() == "ResponseWriter" {
				// http.ResponseWriter
				args = append(args, reflect.ValueOf(writer))
				continue
			} else if param.RealType.Kind() == reflect.Ptr && param.ElemType.Kind() == reflect.Struct && param.ElemType.Name() == "Request" {
				// http.Request
				args = append(args, reflect.ValueOf(request))
				continue
			}
			return nil, common.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The net/http param only support http.ResponseWriter and *http.Request, but handler [%v] contains [%v].", handlerName, param.ElemType.Name()))
		}

		// ----------------------------------------------------------------------------------------------     RESTful    ----------------------------------------------------------------------------------------------
		// RESTful 格式传参
		if param.InPath && pathVariables != nil {
			index := util.GetPathVariableIndex(param.Name, hw.Path)
			if index > -1 {
				// 找到啦~
				// 找到位置后，去 url path 的对应位置上获取参数值即可
				temp := GetPathVariableValue(index, path)
				// 添加到参数列表
				args = append(args, util.StringToValue(param.RealType.Kind(), temp))
				continue
			}
			return nil, common.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The path [%v] does not contains path variable [%v].", hw.Path, param.Name))
		}

		// ----------------------------------------------------------------------------------------------     Header     ----------------------------------------------------------------------------------------------
		// 从请求头获取
		if param.InHeader {
			temp := getHeaderParam(request, param.Name)
			if temp == "" && param.Required {
				return nil, common.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The param [%v] is nested, but received value is empty.", param.Name))
			}
			// 添加到参数列表
			args = append(args, util.StringToValue(param.RealType.Kind(), temp))
			continue
		}

		// ----------------------------------------------------------------------------------------------   RequestBody  ----------------------------------------------------------------------------------------------
		// 有 requestBody 参数
		// 仅支持 POST 和 PUT 方法
		if param.IsBody {
			if request.Method != http.MethodPost && request.Method != http.MethodPut {
				return nil, common.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("RequestBody only support 'POST' and 'PUT' method, but now is [%v].", request.Method))
			}

			if !VerifyMethod(hw, http.MethodPost) && !VerifyMethod(hw, http.MethodPut) {
				return nil, common.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Maybe the handler [%v] should be register as 'POST' or 'PUT' method, now is %v.", handlerName, hw.Methods))
			}

			// 获取到 requestBody
			bs := util.RecycleRequestBody(request)
			if bs != nil {
				var arg interface{}
				// switch 即使两个 case 条件相同，也不能合并，必须分开
				switch param.RealType.Kind() {
				case reflect.Map:
					arg = reflect.MakeMap(param.ElemType).Interface()
					break
				case reflect.Slice:
					arg = reflect.MakeSlice(param.ElemType, 0, 0).Interface()
					break
				case reflect.Array:
					arg = reflect.New(reflect.ArrayOf(0, param.ElemType))
					break
				case reflect.Struct:
					arg = reflect.New(param.ElemType).Interface()
					break
				case reflect.Ptr:
					arg = reflect.New(param.ElemType).Interface()
					break
				}

				// json 解码
				err := json.Unmarshal(bs, arg)
				if err != nil {
					gog.Error(err)
					continue
				}

				val := reflect.ValueOf(arg)
				// 如果接收的是 struct 类型，需要从指针中获取到 struct
				if param.RealType.Kind() == reflect.Struct {
					val = reflect.ValueOf(arg).Elem()
				}
				// 添加到参数列表
				args = append(args, val)
				continue
			}
			return nil, common.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("There is no request body of request [%v], but it's nested.", hw.Path))
		}

		// ----------------------------------------------------------------------------------------------       VO       -----------------------------------------------------------------------------------------------
		// 判断参数类型是否是 VO 类
		if param.ElemType.Kind() == reflect.Struct {
			temp := reflect.New(param.ElemType)

			// 装配VO模型
			temp, ex := getVOParam(request, temp.Interface())
			if ex != nil {
				return nil, ex
			}
			// 添加到参数列表
			// 如果接收的是 struct 类型，需要从指针中获取到 struct
			if param.RealType.Kind() == reflect.Struct {
				temp = temp.Elem()
			}
			args = append(args, temp)
			continue
		}

		// ----------------------------------------------------------------------------------------------  MultipartFile ----------------------------------------------------------------------------------------------
		// 文件上传
		// 兼容 MultipartFile 和 *MultipartFile 两种类型
		if t := reflect.TypeOf(common.MultipartFile{}); param.ElemType.Name() == t.Name() || param.IsPtr && param.ElemType.Name() == t.Name() {
			file, header, err := request.FormFile(param.Name)
			if err != nil {
				return nil, common.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			mf := &common.MultipartFile{
				Header: header,
				File:   file,
			}

			// 根据具体接收类型（对象|指针）装配参数
			var val reflect.Value
			if param.IsPtr {
				val = reflect.ValueOf(mf)
			} else {
				val = reflect.ValueOf(*mf)
			}
			args = append(args, val)
			continue
		}

		// ----------------------------------------------------------------------------------------------    Normal    ----------------------------------------------------------------------------------------------
		// 普通参数
		// Query / Form / PostForm
		temp := getNormalParam(request, param.Name)
		// 如果没获取到参数但又必须，则直接报错
		if temp == "" && param.Required {
			return nil, common.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The param [%v] is nested, but no value received.", param.Name))
		}
		// 添加到参数列表
		args = append(args, util.StringToValue(param.RealType.Kind(), temp))
	}
	return args, nil
}

// getVOParam 自动填充 VO 参数
func getVOParam(request *http.Request, value interface{}) (reflect.Value, *common.HTTPError) {
	realValue := reflect.ValueOf(value)
	// 临时变量，便于从指针转换为 struct
	elemValue := realValue
	if elemValue.Kind() == reflect.Ptr {
		elemValue = elemValue.Elem()
	}
	elemType := elemValue.Type()
	count := elemType.NumField()

	for i := 0; i < count; i++ {
		field := elemType.Field(i)
		val, ex := getVOParamValueByField(request, field)
		if ex != nil {
			return val, ex
		}
		if realValue.Kind() == reflect.Ptr {
			// 如果 对象 是 指针，则可以直接设置字段值
			utils.FieldSet(elemValue.Field(i), val)
		} else {
			// 否则需要到指针指向的对象，再设置值
			utils.FieldSet(elemValue.Elem().Field(i), val)
		}
	}
	// 返回原始 数据
	return realValue, nil
}

// getVOParamValueByField 填充 VO 字段
func getVOParamValueByField(request *http.Request, field reflect.StructField) (reflect.Value, *common.HTTPError) {
	name := field.Tag.Get("param")
	if name == "" {
		name = util.FirstToLower(field.Name)
	}
	_, inHeader := field.Tag.Lookup("header")
	_, required := field.Tag.Lookup("required")

	var temp string
	if inHeader {
		temp = getHeaderParam(request, name)
		if temp == "" {
			temp = getHeaderParam(request, field.Name)
		}
	} else {
		temp = getNormalParam(request, name)
		if temp == "" {
			temp = getNormalParam(request, field.Name)
		}
	}
	if temp == "" && required {
		return reflect.ValueOf(nil), common.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("The param [%v] is nested, but no value received.", name))
	}
	return util.StringToValue(field.Type.Kind(), temp), nil
}

// getHeaderParam 从请求头中获取参数
func getHeaderParam(request *http.Request, name string) string {
	return request.Header.Get(name)
}

// getNormalParam 普通参数获取
func getNormalParam(request *http.Request, name string) string {
	// 先从 URL 中获取
	temp := request.URL.Query().Get(name)
	if temp == "" {
		// 从 form 中获取
		temp = request.FormValue(name)
		if temp == "" && request.Method == http.MethodPost {
			temp = request.PostFormValue(name)
		}
	}
	return temp
}

// GetPathVariableValue 用 参数位置 从 实际 path 中获取到参数值
func GetPathVariableValue(index int, path string) string {
	nodes := strings.Split(path, "/")
	if nodes != nil && index < len(nodes) {
		return nodes[index]
	}
	return ""
}

// VerifyMethod 校验请求方法
func VerifyMethod(hw *wire.HandlerWire, method string) bool {
	for _, md := range hw.Methods {
		if string(md) == method {
			return true
		}
	}
	return false
}
