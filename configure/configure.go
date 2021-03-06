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
// time   : 2019-11-26 11:51
// version: 1.0.0
// desc   : Web配置接口

package configure

import (
	"github.com/yhyzgn/gox/component/filter"
	"github.com/yhyzgn/gox/component/interceptor"
	"github.com/yhyzgn/gox/ctx"
)

// WebConfigure Web配置接口
type WebConfigure interface {
	// Context 配置 Context
	Context(ctx *ctx.GoXContext)

	// ConfigFilter 注册过滤器
	ConfigFilter(chain *filter.Chain)

	// ConfigInterceptor 注册拦截器
	ConfigInterceptor(register *interceptor.Register)
}
