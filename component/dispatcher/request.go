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
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/component/interceptor"
	"github.com/yhyzgn/gox/ctx"
	"github.com/yhyzgn/gox/resolver"
	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gox/wire"
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

	// 如果路由未找到，可能是静态资源
	if reqPath == "" {
		// 默认首页
		reqPath = "index.html"
	}
	filename := strings.ReplaceAll(ctx.C().GetStaticDir()+"/"+reqPath, "//", "/")
	if util.FileExist(filename) {
		http.ServeFile(writer, request, filename)
		return
	}

	// 匹配不到，就只能 404 啦~
	ctx.C().GetNotFoundHandler()(writer, request)
}

// doDispatch 具体的请求分发操作
func (rd *RequestDispatcher) doDispatch(hw *wire.HandlerWire, writer http.ResponseWriter, request *http.Request, isRESTful bool) {
	md := VerifyMethod(hw, request.Method)
	if !md {
		// 不支持的 http 方法
		ctx.C().GetUnSupportMethodHandler()(writer, request)
		return
	}

	// 处理器
	handler := hw.Handler

	// 参数处理器
	argumentResolver := ctx.GetWare(common.ArgumentResolverName, resolver.NewSimpleArgumentResolver()).(resolver.ArgumentResolver)
	// 结果处理器
	resultResolver := ctx.GetWare(common.ResultResolverName, resolver.NewSimpleResultResolver()).(resolver.ResultResolver)

	// 先处理一遍参数
	args, ex := rd.resolve(hw, writer, request, isRESTful)
	if ex != nil {
		gog.Error(ex.Error)
		ex.Response(writer)
		return
	}

	// 再调用参数处理器处理
	argumentResolver.Resolve(args, writer, request, handler)

	gog.DebugF("Params of request path [{}] are [{}], matched router [{}] of params [{}]", request.URL.Path, util.FormatRealArgsValue(args), hw.Path, util.FormatHandlerArgs(hw.Params))

	// 处理前，执行拦截器 PreHandle() 方法
	if rd.register != nil && !util.IsExcludedRequest(request, rd.register.GetExcludes()) {
		passed, path := rd.register.Iterate(func(index int, path string, interceptor interceptor.Interceptor) (skipped, passed bool) {
			// 匹配 path，未匹配到的直接跳过
			if path == "/" {
				// 所有请求
				skipped = false
				passed, request, writer = interceptor.PreHandle(writer, request, handler)
			} else if path == request.URL.Path {
				// 严格匹配，只有路径完全相同才走过滤器
				skipped = false
				passed, request, writer = interceptor.PreHandle(writer, request, handler)
			} else if util.MatchedRequestByPathPattern(request, path) {
				// 前缀匹配成功，执行拦截器
				skipped = false
				passed, request, writer = interceptor.PreHandle(writer, request, handler)
			} else {
				// 跳过
				skipped = true
				passed = false
			}

			if skipped {
				gog.DebugF("The request [%v] has skipped by interceptor [%v].", request.URL.Path, path)
			} else if passed {
				gog.DebugF("The request [%v] has passed by interceptor [%v].", request.URL.Path, path)
			} else {
				gog.DebugF("The request [%v] has been intercepted by interceptor [%v].", request.URL.Path, path)
			}

			return
		})

		// 拦截器不通过
		if !passed {
			gog.InfoF("The request [%v] has been intercepted by interceptor [%v].", request.URL.Path, path)
			return
		}
	}

	// 将request和writer设置回请求中
	for i, arg := range args {
		// ----------------------------------------------------------------------------------------------    net/http    ----------------------------------------------------------------------------------------------
		// http.ResponseWriter || *http.Request
		if arg.Type().Elem() != nil && arg.Type().Elem().PkgPath() == "net/http" {
			if arg.Type().Elem().Kind() == reflect.Interface && arg.Type().Elem().Name() == "ResponseWriter" {
				// http.ResponseWriter
				args[i] = reflect.ValueOf(writer)
			} else if arg.Type().Kind() == reflect.Ptr && arg.Type().Elem().Kind() == reflect.Struct && arg.Type().Elem().Name() == "Request" {
				// http.Request
				args[i] = reflect.ValueOf(request)
			}
		}
	}

	// 拦截器通过后，将请求交由 处理器 处理
	// 已经获取到参数列表，执行方法即可
	results := handler.Call(args)
	noResult := results == nil || len(results) == 0
	var (
		res reflect.Value
		err error
	)
	if noResult {
		// 无返回值
		gog.InfoF("The request [{}] responded, and the handler needn't return any value.", request.URL.Path)
	} else {
		// 响应结果交由 结果处理器 处理
		res, err = resultResolver.Resolve(hw, results, writer, request)
		// 如果有错误，就响应错误信息
		if err != nil {
			res = reflect.ValueOf(err)
		}
	}

	// 处理完成后，执行拦截器的 AfterHandle() 方法
	if rd.register != nil && !util.IsExcludedRequest(request, rd.register.GetExcludes()) {
		rd.register.ReverseIterate(func(index int, path string, interceptor interceptor.Interceptor) {
			// 匹配 path，未匹配到的直接跳过
			if path == "/" {
				// 所有请求
				request, writer = interceptor.AfterHandle(writer, request, handler, res, err)
			} else if reg, e := regexp.Compile("/\\*+$"); e == nil && reg.MatchString(path) {
				// 前缀匹配
				pattern := reg.ReplaceAllString(path, "/.+?")
				if matched, err := regexp.MatchString("^"+pattern+"$", request.URL.Path); matched && err == nil {
					// 前缀匹配成功，执行拦截器
					request, writer = interceptor.AfterHandle(writer, request, handler, res, err)
				}
			} else if path == request.URL.Path {
				// 严格匹配，只有路径完全相同才走过滤器
				request, writer = interceptor.AfterHandle(writer, request, handler, res, err)
			}
		})
	}

	if !noResult {
		// 拦截器通过后，响应处理结果
		resultResolver.Response(res, writer)
	}
}

// resolve 初步处理参数
func (rd *RequestDispatcher) resolve(hw *wire.HandlerWire, writer http.ResponseWriter, request *http.Request, isRESTful bool) ([]reflect.Value, *common.HTTPError) {
	path := request.URL.Path
	handler := reflect.Value(hw.Handler)
	handlerName := strings.ReplaceAll(runtime.FuncForPC(handler.Pointer()).Name(), "-fm", util.FormatHandlerArgs(hw.Params))

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

		// ----------------------------------------------------------------------------------------------  MultipartFile ----------------------------------------------------------------------------------------------
		// 文件上传
		// 兼容 MultipartFile 和 *MultipartFile && []MultipartFile 和 []*MultipartFile 四种类型
		if isMultipart, multi, isPtr := isMultipartFile(param); isMultipart {
			files, headers, err := util.FormFiles(request, param.Name)
			if err != nil {
				return nil, common.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			var val reflect.Value
			if multi {
				// 多个文件
				if isPtr {
					// 指针类型
					temp := make([]*common.MultipartFile, 0)
					for i, file := range files {
						temp = append(temp, &common.MultipartFile{
							Header: headers[i],
							File:   file,
						})
					}
					val = reflect.ValueOf(temp)
				} else {
					// 普通类型
					temp := make([]common.MultipartFile, 0)
					for i, file := range files {
						temp = append(temp, common.MultipartFile{
							Header: headers[i],
							File:   file,
						})
					}
					val = reflect.ValueOf(temp)
				}
			} else {
				// 单个文件
				mf := &common.MultipartFile{
					Header: headers[0],
					File:   files[0],
				}
				// 根据具体接收类型（对象|指针）装配参数
				if param.IsPtr {
					val = reflect.ValueOf(mf)
				} else {
					val = reflect.ValueOf(*mf)
				}
			}

			// 添加参数到参数列表
			args = append(args, val)
			continue
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

// 是否是文件上传
//
// 返回值： 是否是文件上传，是否有多个文件，是否是文件指针
func isMultipartFile(param *common.Param) (isMultipart, multi, isPtr bool) {
	tp := reflect.TypeOf(new(common.MultipartFile))

	// (file MultipartFile) || (file *MultipartFile)
	if param.RealType.Name() == tp.Elem().Name() || param.IsPtr && param.ElemType.Name() == tp.Elem().Name() {
		// 接收单个文件
		isMultipart = true
		isPtr = param.IsPtr
		return
	}

	// (files []MultipartFile) || (files []*MultipartFile)
	if param.RealType.Kind() == reflect.Slice {
		// 多个文件
		itemType := param.RealType.Elem()
		elemType := itemType
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
			isPtr = true
		}

		// 匹配类型
		if elemType.Kind() == tp.Elem().Kind() {
			isMultipart = true
			multi = true
		}
	}
	return
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
			util.FieldSet(elemValue.Field(i), val)
		} else {
			// 否则需要到指针指向的对象，再设置值
			util.FieldSet(elemValue.Elem().Field(i), val)
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
