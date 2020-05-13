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
// time   : 2019-10-11 3:25
// version: 1.0.0
// desc   : 配置初始化

package resource

import "github.com/yhyzgn/gox/ioc"

func init() {
	// 无名称注入，使用时按类型获取
	ioc.C().Add(func() (instance interface{}) {
		instance = NewReader()
		return
	})
}
