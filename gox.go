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
// time   : 2019-11-24 12:47 上午
// version: 1.0.0
// desc   : 

package gox

import (
	"github.com/yhyzgn/ghost/ioc"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/component/dispatcher"
	"github.com/yhyzgn/gox/component/filter"
	"github.com/yhyzgn/gox/component/interceptor"
	"github.com/yhyzgn/gox/configure"
	"github.com/yhyzgn/gox/context"
	"github.com/yhyzgn/gox/core"
	"github.com/yhyzgn/gox/resolver"
	"github.com/yhyzgn/gox/util"
	"net/http"
	"sync"
)

type GoX struct {
	mu sync.RWMutex
}

// 做一些初始化配置
func init() {
	context.Current().
		SetWareOnce(common.FilterChainName, filter.NewFilterChain()). // 过滤器链
		SetWareOnce(common.RequestDispatcherName, dispatcher.NewRequestDispatcher()). // 请求分发器
		SetWareOnce(common.InterceptorRegisterName, interceptor.NewInterceptorRegister()). // 拦截器
		SetWare(common.ArgumentResolverName, resolver.NewSimpleArgumentResolver()). // 参数处理器
		SetWare(common.ResultResolverName, resolver.NewSimpleResultResolver()) // 结果处理器
}

func (r *GoX) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 每个请求 过滤器 开始标记
	request = util.SetRequestAttribute(request, common.RequestFilterIndexName, 0)

	// -----------------------------------------------------------------------
	// 过滤器
	//   ↓
	// 分发器
	//   ↓
	// 拦截器
	//   ↓
	// 处理器
	// -----------------------------------------------------------------------

	// 过滤器链
	filterChain := util.GetWare(common.FilterChainName, filter.NewFilterChain()).(*filter.FilterChain)
	// 分发器
	dispatch := util.GetWare(common.RequestDispatcherName, dispatcher.NewRequestDispatcher()).(*dispatcher.RequestDispatcher)
	// 拦截器
	interceptorRegister := util.GetWare(common.InterceptorRegisterName, interceptor.NewInterceptorRegister()).(*interceptor.InterceptorRegister)

	// 将拦截器设置到分发器
	dispatch.SetInterceptorRegister(interceptorRegister)
	// 将分发器设置到过滤器链
	filterChain.SetDispatcher(dispatch)

	// 开始啦~
	filterChain.DoFilter(writer, request)
}

func NewGoX() *GoX {
	return new(GoX)
}

func (r *GoX) Configure(configure configure.WebConfigure) *GoX {
	r.config(configure)
	return r
}

func (r *GoX) NotFoundHandler(handler http.HandlerFunc) *GoX {
	context.Current().NotFound = handler
	return r
}

func (r *GoX) UnsupportedMethodHandler(handler http.HandlerFunc) *GoX {
	context.Current().UnsupportedMethod = handler
	return r
}

func (r *GoX) Mapping(path string, ctrl core.Controller) *GoX {
	r.mu.Lock()
	defer r.mu.Unlock()
	route := core.NewRoute()
	route.Path(path).Controller(ctrl)
	mapper := core.NewMapper(route)
	ctrl.Mapping(mapper)

	// 注册到 IOC
	ioc.GetProvider().Single("", ctrl)
	return r
}

func (r *GoX) config(configure configure.WebConfigure) {
	if configure != nil {
		// 注册过滤器
		configure.ConfigFilter(util.GetWare(common.FilterChainName, filter.NewFilterChain()).(*filter.FilterChain))

		// 注册拦截器
		configure.ConfigInterceptor(util.GetWare(common.InterceptorRegisterName, interceptor.NewInterceptorRegister()).(*interceptor.InterceptorRegister))
	}
}
