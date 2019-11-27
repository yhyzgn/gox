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
// desc   : 

package wire

import (
	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/util"
	"github.com/yhyzgn/gog"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type HandlerWire struct {
	Path    string
	Handler common.Handler
	Methods []common.Method
	Params  []*common.Param
}

type Wires struct {
	wires  map[string]*HandlerWire
	sorted []*HandlerWire
}

var (
	once     sync.Once
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

func (w *Wires) Mapping(path string, handler common.Handler, methods []common.Method, params []*common.Param) {
	wire := &HandlerWire{
		Path:    path,
		Handler: handler,
		Methods: methods,
		Params:  params,
	}
	// Path 节点  或者  路径长度 从长到端排序
	w.sorted = appendSorted(w.sorted, wire)
	w.wires[path] = wire

	pc := reflect.Value(handler).Pointer()
	name := util.ReplaceAll(runtime.FuncForPC(pc).Name(), "-fm", "(...)")
	gog.InfoF("Mapped [%v --> %v] with http method %v", path, name, methods)
}

func (w *Wires) Get(path string) *HandlerWire {
	return w.wires[path]
}

func (w *Wires) All() []*HandlerWire {
	return w.sorted
}

func appendSorted(wires []*HandlerWire, wire *HandlerWire) []*HandlerWire {
	length := len(wires)
	index := sort.Search(length, func(i int) bool {
		tempNodeCount := strings.Split(wires[i].Path, "/")
		wireNodeCount := strings.Split(wire.Path, "/")
		// Path 节点  或者  路径长度 从长到端排序
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
