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
// time   : 2019-11-26 12:43
// version: 1.0.0
// desc   : 过滤器接口
//			所有请求将先到达 过滤器

package filter

import (
	"net/http"
)

// Filter 过滤器
type Filter interface {
	// DoFilter 执行过滤器操作
	DoFilter(writer http.ResponseWriter, request *http.Request, chain *Chain)
}
