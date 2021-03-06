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
// time   : 2019-11-26 14:47
// version: 1.0.0
// desc   : 过滤器链-责任链模式

package filter

import (
	"net/http"
	"strings"
	"sync"

	"github.com/yhyzgn/gox/ctx"

	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/component/dispatcher"
	"github.com/yhyzgn/gox/util"
)

type item struct {
	path   string
	filter Filter
}

// Chain 过滤器链
type Chain struct {
	filters    []item
	excludes   sync.Map
	dispatcher dispatcher.Dispatcher
}

// NewChain 一个新链
func NewChain() *Chain {
	return &Chain{
		filters: make([]item, 0),
	}
}

// SetDispatcher 设置请求分发器
// 当所有过滤器执行完后，需要执行该分发器
func (fc *Chain) SetDispatcher(dispatcher dispatcher.Dispatcher) {
	fc.dispatcher = dispatcher
}

// AddFilters 向链中添加过滤器
// 添加顺序 即 执行顺序
// path 匹配方式：
// 				/		->		所有请求
//				/xx		->		严格匹配
//				/xx/*	->		前缀匹配
func (fc *Chain) AddFilters(path string, filters ...Filter) *Chain {
	if path == "" || filters == nil || len(filters) == 0 {
		return fc
	}
	for _, flt := range filters {
		fc.filters = append(fc.filters, item{
			path:   path,
			filter: flt,
		})
	}
	gog.DebugF("The Filters [%v] registered.", path)
	return fc
}

// Exclude 添加排除路径
//
// 支持 前缀匹配 & 严格匹配
func (fc *Chain) Exclude(path string) *Chain {
	if _, ok := fc.excludes.Load(path); !ok {
		fc.excludes.Store(path, true)
	}
	return fc
}

// GetExcludes 获取那些被排除的路径
func (fc *Chain) GetExcludes() map[string]bool {
	excludes := make(map[string]bool)
	fc.excludes.Range(func(key, value interface{}) bool {
		excludes[key.(string)] = value.(bool)
		return true
	})
	return excludes
}

// DoFilter 逐个执行过滤器
// 执行顺序 为 添加顺序
func (fc *Chain) DoFilter(writer http.ResponseWriter, request *http.Request) {
	// 匹配时忽略ContextPath
	reqPath := strings.ReplaceAll(request.URL.Path, ctx.C().GetContextPath(), "")
	// 先判断这些请求是否已经被排除在 过滤器 外
	if util.IsExcludedRequest(reqPath, fc.GetExcludes()) {
		gog.TraceF("The request [%v] has been excluded", request.URL.Path)
		fc.dispatcher.Dispatch(writer, request)
		return
	}

	// 计算索引
	request, index := fc.getIndexAndIncrement(request)

	// 走完所有过滤器，需要将请求交给 dispatcher
	if index == len(fc.filters) {
		fc.dispatcher.Dispatch(writer, request)
		return
	}

	item := fc.filters[index]

	// 匹配 path，未匹配到的 filter 直接跳过
	if item.path == "/" {
		// 所有请求
		gog.TraceF("The request [%v] has passed by filter [/]", request.URL.Path)
		item.filter.DoFilter(writer, request, fc)
	} else if item.path == reqPath {
		// 严格匹配，只有路径完全相同才走过滤器
		gog.TraceF("The request [%v] has passed by filter [%v]", request.URL.Path, item.path)
		item.filter.DoFilter(writer, request, fc)
	} else if util.MatchedRequestByPathPattern(reqPath, item.path) {
		// 前缀|后缀匹配成功，走过滤器
		gog.TraceF("The request [%v] has passed by filter [%v]", request.URL.Path, item.path)
		item.filter.DoFilter(writer, request, fc)
	} else {
		gog.TraceF("The request [%v] has skipped by filter [%v]", request.URL.Path, item.path)
		// 匹配不到过滤器，则递归回当前链，继续下一次匹配
		fc.DoFilter(writer, request)
	}
}

// getIndexAndIncrement 获取当前索引，并递增
func (fc *Chain) getIndexAndIncrement(req *http.Request) (*http.Request, int) {
	// 获取到当前请求中的过滤器 索引
	index := util.GetRequestAttribute(req, common.RequestFilterIndexName).(int)
	// 实时更新 索引
	return util.SetRequestAttribute(req, common.RequestFilterIndexName, index+1), index
}
