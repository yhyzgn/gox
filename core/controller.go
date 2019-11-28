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
// time   : 2019-11-24 3:52 上午
// version: 1.0.0
// desc   : 控制器接口
//
// 			所有控制器必须实现该接口，并在 Mapping() 方法中完成 处理器 注册

package core

// Controller 控制器接口
type Controller interface {
	// Mapping 在该方法中完成当前控制器中的所有处理器注册
	//
	// mapper.Request("/rest/{id}").HandlerFunc(h.Rest).PathVariable("id").Mapping()
	Mapping(*Mapper)
}
