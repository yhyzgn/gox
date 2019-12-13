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
// desc   : MVC 入口

package gox

import (
	"net/http"
	"sync"

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
)

// GoX MVC 服务处理器
type GoX struct {
	mu sync.RWMutex
}

// 做一些初始化配置
func init() {
	context.Current().
		SetWareOnce(common.FilterChainName, filter.NewChain()).                       // 过滤器链
		SetWareOnce(common.RequestDispatcherName, dispatcher.NewRequestDispatcher()). // 请求分发器
		SetWareOnce(common.InterceptorRegisterName, interceptor.NewRegister()).       // 拦截器
		SetWare(common.ArgumentResolverName, resolver.NewSimpleArgumentResolver()).   // 参数处理器
		SetWare(common.ResultResolverName, resolver.NewSimpleResultResolver())        // 结果处理器
}

// ServeHTTP 接收处理请求
func (gx *GoX) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.RequestURI == "*" {
		if request.ProtoAtLeast(1, 1) {
			util.SetResponseWriterHeader(writer, "Connection", "closed")
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

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
	filterChain := util.GetWare(common.FilterChainName, filter.NewChain()).(*filter.Chain)
	// 分发器
	dispatch := util.GetWare(common.RequestDispatcherName, dispatcher.NewRequestDispatcher()).(*dispatcher.RequestDispatcher)
	// 拦截器
	interceptorRegister := util.GetWare(common.InterceptorRegisterName, interceptor.NewRegister()).(*interceptor.Register)

	// 将拦截器设置到分发器
	dispatch.SetInterceptorRegister(interceptorRegister)
	// 将分发器设置到过滤器链
	filterChain.SetDispatcher(dispatch)

	// 开始啦~
	filterChain.DoFilter(writer, request)
}

// NewGoX 创建新服务
func NewGoX() *GoX {
	return new(GoX)
}

// Configure 配置 Web
func (gx *GoX) Configure(configure configure.WebConfigure) *GoX {
	gx.config(configure)
	return gx
}

// StaticDir 静态资源文件夹
func (gx *GoX) StaticDir(dir string) *GoX {
	context.Current().SetStaticDir(dir)
	return gx
}

// NotFoundHandler 配置 404 处理器
func (gx *GoX) NotFoundHandler(handler http.HandlerFunc) *GoX {
	context.Current().SetNotFoundHandler(handler)
	return gx
}

// UnsupportedMethodHandler 配置 方法不支持 处理器
func (gx *GoX) UnsupportedMethodHandler(handler http.HandlerFunc) *GoX {
	context.Current().SetUnSupportMethodHandler(handler)
	return gx
}

// ErrorCodeHandler 为错误码添加处理器
func (gx *GoX) ErrorCodeHandler(statusCode int, handler http.HandlerFunc) *GoX {
	context.Current().AddErrorHandler(statusCode, handler)
	return gx
}

// Mapping 添加 控制器 映射
func (gx *GoX) Mapping(path string, ctrls ...core.Controller) *GoX {
	if ctrls == nil || len(ctrls) == 0 {
		return gx
	}
	gx.mu.Lock()
	defer gx.mu.Unlock()
	// 逐个添加
	for _, ctrl := range ctrls {
		// 创建一个 处理器映射器对象
		mapper := core.NewMapper(path, ctrl)
		// 执行每个控制器的 Mapping() 方法，完成 处理器的注册
		ctrl.Mapping(mapper)
		// 注册到 IOC
		ioc.GetProvider().Single("", ctrl)
	}
	return gx
}

// config 触发配置装载
func (gx *GoX) config(configure configure.WebConfigure) {
	if configure != nil {
		// 配置 Context
		configure.Context(context.Current())

		// 注册过滤器
		configure.ConfigFilter(util.GetWare(common.FilterChainName, filter.NewChain()).(*filter.Chain))

		// 注册拦截器
		configure.ConfigInterceptor(util.GetWare(common.InterceptorRegisterName, interceptor.NewRegister()).(*interceptor.Register))
	}
}
