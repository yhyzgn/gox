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
	excludes   map[string]bool
	dispatcher dispatcher.Dispatcher
}

// NewChain 一个新链
func NewChain() *Chain {
	return &Chain{
		filters:  make([]item, 0),
		excludes: make(map[string]bool),
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
	gog.InfoF("The Filters [%v] registered.", path)
	return fc
}

// Exclude 添加排除路径
//
// 支持 前缀匹配 & 严格匹配
func (fc *Chain) Exclude(path string) *Chain {
	if !fc.excludes[path] {
		fc.excludes[path] = true
	}
	return fc
}

// DoFilter 逐个执行过滤器
// 执行顺序 为 添加顺序
func (fc *Chain) DoFilter(writer http.ResponseWriter, request *http.Request) {
	// 先判断这些请求是否已经被排除在 过滤器 外
	if util.IsExcludedRequest(request, fc.excludes) {
		gog.DebugF("The request [%v] has been excluded", request.URL.Path)
		fc.dispatcher.Dispatch(writer, request)
		return
	}

	// 获取到当前请求中的过滤器 索引
	index := util.GetRequestAttribute(request, common.RequestFilterIndexName).(int)
	// 实时更新 索引
	request = util.SetRequestAttribute(request, common.RequestFilterIndexName, index+1)

	// 走完所有过滤器，需要将请求交给 dispatcher
	if index == len(fc.filters) {
		fc.dispatcher.Dispatch(writer, request)
		return
	}

	item := fc.filters[index]

	// 匹配 path，未匹配到的 filter 直接跳过
	if item.path == "/" {
		// 所有请求
		gog.DebugF("The request [%v] has passed by filter [/]", request.URL.Path)
		item.filter.DoFilter(writer, request, fc)
	} else if item.path == request.URL.Path {
		// 严格匹配，只有路径完全相同才走过滤器
		gog.DebugF("The request [%v] has passed by filter [%v]", request.URL.Path, item.path)
		item.filter.DoFilter(writer, request, fc)
	} else if util.MatchedRequestByPathPattern(request, item.path) {
		// 前缀匹配成功，走过滤器
		gog.DebugF("The request [%v] has passed by filter [%v]", request.URL.Path, item.path)
		item.filter.DoFilter(writer, request, fc)
	} else {
		gog.DebugF("The request [%v] has skipped by filter [%v]", request.URL.Path, item.path)
		// 匹配不到过滤器，则递归回当前链，继续下一次匹配
		fc.DoFilter(writer, request)
	}
}
