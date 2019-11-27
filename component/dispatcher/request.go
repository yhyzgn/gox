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
// desc   : 

package dispatcher

import (
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/component/interceptor"
	"github.com/yhyzgn/gox/resolver"
	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gox/wire"
	"net/http"
	"reflect"
	"regexp"
)

type RequestDispatcher struct {
	register *interceptor.InterceptorRegister
}

func NewRequestDispatcher() *RequestDispatcher {
	return new(RequestDispatcher)
}

func (rd *RequestDispatcher) SetInterceptorRegister(register *interceptor.InterceptorRegister) {
	rd.register = register
}

func (rd *RequestDispatcher) Dispatch(writer http.ResponseWriter, request *http.Request) {
	for _, h := range wire.Instance.All() {
		// 如果直接完全匹配，说明不是 RESTFul 模式
		if request.URL.Path == h.Path {
			rd.doDispatch(h, writer, request, false)
			return
		} else if util.IsRESTFul(h.Path) {
			// 否则 正则匹配
			// 将 路由注册的路径 转换为 正则匹配模板，再看是否与真实路径匹配
			realPathPattern := util.ConvertRESTFulPathToPattern(h.Path)
			matched, err := regexp.MatchString(realPathPattern, request.URL.Path)
			if err == nil && matched {
				// RESTFul 匹配上了
				rd.doDispatch(h, writer, request, true)
				return
			}
		}
	}

	// TODO 匹配不到，就只能 404 啦~🌶🌶
	http.NotFound(writer, request)
}

func (rd *RequestDispatcher) doDispatch(hw *wire.HandlerWire, writer http.ResponseWriter, request *http.Request, isRESTFul bool) {
	md := resolver.VerifyMethod(hw, request.Method)
	if !md {
		// TODO 不支持的 http 方法
		return
	}

	// 处理器
	handler := reflect.Value(hw.Handler)

	// 参数处理器
	argumentResolver := util.GetComponent(common.ArgumentResolverName, resolver.NewSimpleArgumentResolver()).(resolver.ArgumentResolver)
	// 结果处理器
	resultResolver := util.GetComponent(common.ResultResolverName, resolver.NewSimpleResultResolver()).(resolver.ResultResolver)

	// 获取到处理后的参数
	args := argumentResolver.Resolve(hw, writer, request, isRESTFul)
	gog.TraceF("Params of request path [{}] are {}, matched router [{}] of params {}", request.URL.Path, args, hw.Path, hw.Params)

	// 处理前，执行拦截器 PreHandle() 方法
	if rd.register != nil {
		pass, path := rd.register.Iterate(func(index int, path string, interceptor interceptor.Interceptor) (skip, passed bool) {
			// 匹配 path，未匹配到的直接跳过
			if path == "/" {
				// 所有请求
				return false, interceptor.PreHandle(writer, request, handler)
			} else if reg, e := regexp.Compile("/\\*+$"); e == nil && reg.MatchString(path) {
				// 前缀匹配
				pattern := reg.ReplaceAllString(path, "/.+?")
				if matched, err := regexp.MatchString("^"+pattern+"$", request.URL.Path); matched && err == nil {
					// 前缀匹配成功，执行拦截器
					return false, interceptor.PreHandle(writer, request, handler)
				} else {
					// 匹配不成功的直接跳过
					return true, true
				}
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
			gog.TraceF("Current request [%v] has been intercepted by interceptor [%v].", request.URL.Path, path)
			return
		}
	}

	// 拦截器通过后，将请求交由 处理器 处理
	// 已经获取到参数列表，执行方法即可
	results := handler.Call(args)
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
