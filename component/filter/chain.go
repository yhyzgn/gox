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
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/component/dispatcher"
	"github.com/yhyzgn/gox/util"
	"net/http"
	"regexp"
)

// Chain 过滤器链
type Chain struct {
	filters    []Filter
	pathMap    map[int]string
	dispatcher dispatcher.Dispatcher
}

// NewChain 一个新链
func NewChain() *Chain {
	return &Chain{
		filters: make([]Filter, 0),
		pathMap: make(map[int]string),
	}
}

// SetDispatcher 设置请求分发器
// 当所有过滤器执行完后，需要执行该分发器
func (fc *Chain) SetDispatcher(dispatcher dispatcher.Dispatcher) {
	fc.dispatcher = dispatcher
}

// AddFilter 向链中添加过滤器
// 添加顺序 即 执行顺序
// path 匹配方式：
// 				/		->		所有请求
//				/xx		->		严格匹配
//				/xx/*	->		前缀匹配
func (fc *Chain) AddFilter(path string, filter Filter) *Chain {
	fc.filters = append(fc.filters, filter)
	fc.pathMap[len(fc.filters)-1] = path
	gog.InfoF("The Filter [%v] registered.", path)
	return fc
}

// DoFilter 逐个执行过滤器
// 执行顺序 为 添加顺序
func (fc *Chain) DoFilter(writer http.ResponseWriter, request *http.Request) {
	// 获取到当前请求中的过滤器 索引
	index := util.GetRequestAttribute(request, common.RequestFilterIndexName).(int)
	// 实时更新 索引
	request = util.SetRequestAttribute(request, common.RequestFilterIndexName, index+1)

	// 走完所有过滤器，需要将请求交给 dispatcher
	if index == len(fc.filters) {
		fc.dispatcher.Dispatch(writer, request)
		return
	}

	filter := fc.filters[index]
	path := fc.pathMap[index]

	// 匹配 path，未匹配到的 filter 直接跳过
	if path == "/" {
		// 所有请求
		gog.InfoF("The request path [%v] matched filter path [/]", request.URL.Path)
		filter.DoFilter(writer, request, fc)
	} else if reg, e := regexp.Compile("/\\*+$"); e == nil && reg.MatchString(path) {
		// 前缀匹配
		pattern := reg.ReplaceAllString(path, "/.+?")
		if matched, err := regexp.MatchString("^"+pattern+"$", request.URL.Path); matched && err == nil {
			// 前缀匹配成功，走过滤器
			gog.InfoF("The request path [%v] matched filter path [%v]", request.URL.Path, path)
			filter.DoFilter(writer, request, fc)
		} else {
			gog.InfoF("The request path [%v] has not matched filter path [%v]", request.URL.Path, path)
		}
	} else if path == request.URL.Path {
		// 严格匹配，只有路径完全相同才走过滤器
		gog.InfoF("The request path [%v] matched filter path [%v]", request.URL.Path, path)
		filter.DoFilter(writer, request, fc)
	} else {
		gog.InfoF("The request path [%v] has not matched filter path [%v]", request.URL.Path, path)
	}
}
