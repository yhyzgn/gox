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
// time   : 2019-11-26 15:32
// version: 1.0.0
// desc   : 请求分发器
//			所有过滤器执行完后，请求将到达 分发器

package dispatcher

import (
	"net/http"
)

// Dispatcher 请求分发器
type Dispatcher interface {
	// Dispatch 分发请求
	Dispatch(writer http.ResponseWriter, request *http.Request)
}
