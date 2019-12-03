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
// time   : 2019-11-24 2:11 上午
// version: 1.0.0
// desc   : 处理器映射缓存

package wire

import (
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/util"
)

// HandlerWire 处理器映射缓存
type HandlerWire struct {
	Path    string          // 配置的 path
	Handler common.Handler  // 处理器
	Methods []common.Method // 请求方法
	Params  []*common.Param // 参数列表
}

// Wires 处理器映射缓存
type Wires struct {
	wires  map[string]*HandlerWire // path 处理器映射
	sorted []*HandlerWire          // 从长到短 排序后的映射
}

var (
	once sync.Once
	// Instance 一个全局实例
	Instance *Wires
)

func init() {
	once.Do(func() {
		Instance = &Wires{
			wires:  make(map[string]*HandlerWire),
			sorted: make([]*HandlerWire, 0),
		}
	})
}

// Mapping 注册映射关系
func (w *Wires) Mapping(path string, handler common.Handler, methods []common.Method, params []*common.Param) {
	wire := &HandlerWire{
		Path:    path,
		Handler: handler,
		Methods: methods,
		Params:  params,
	}
	// Request 节点  或者  路径长度 从长到端排序
	w.sorted = appendSorted(w.sorted, wire)
	w.wires[path] = wire

	pc := reflect.Value(handler).Pointer()
	name := strings.ReplaceAll(runtime.FuncForPC(pc).Name(), "-fm", util.FormatHandlerArgs(wire.Params))
	gog.InfoF("Mapped [%v-->\t%v] with http methods %v", util.FillSuffix(path, " ", 40), name, methods)
}

// Get 获取一条映射关系
func (w *Wires) Get(path string) *HandlerWire {
	return w.wires[path]
}

// All 获取所有映射关系
func (w *Wires) All() []*HandlerWire {
	return w.sorted
}

// appendSorted 按 path 从长到短 插入数组
func appendSorted(wires []*HandlerWire, wire *HandlerWire) []*HandlerWire {
	length := len(wires)
	index := sort.Search(length, func(i int) bool {
		tempNodeCount := strings.Split(wires[i].Path, "/")
		wireNodeCount := strings.Split(wire.Path, "/")
		// Request 节点  或者  路径长度 从长到端排序
		return len(tempNodeCount) < len(wireNodeCount) || len(wires[i].Path) < len(wire.Path)
	})
	if index == length {
		return append(wires, wire)
	}

	wires = append(wires, &HandlerWire{})
	copy(wires[index+1:], wires[index:])
	wires[index] = wire
	return wires
}
