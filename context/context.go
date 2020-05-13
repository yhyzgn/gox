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
// time   : 2019-11-27 9:24
// version: 1.0.0
// desc   : 一个简单的上下文环境

package context

// AppVersion 信息
type AppVersion struct {
	Name    string
	Version string
}

// ResourceContext 静态资源上下文
type ResourceContext interface {
	// Read 读取静态资源
	Read(filename string) (data []byte, errs error)
}

// WebContext mvc 上下文
type WebContext interface {
	ResourceContext
}

// WareContext 组件管理功能的上下文
type WareContext interface {
	WebContext

	// SetWare 添加组件
	SetWare(name string, component interface{}) WareContext

	// SetWareOnce 添加一次性组件，修改无效
	SetWareOnce(name string, component interface{}) WareContext

	// GetWare 获取组件
	GetWare(name string) interface{}
}

// XContext GoX 上下文
type XContext interface {
	WareContext
}
