// Copyright 2020 yhyzgn gox
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
// time   : 2020-05-13 9:53 下午
// version: 1.0.0
// desc   : WebConfigure适配器

package adapter

import (
	"github.com/yhyzgn/gox/component/filter"
	"github.com/yhyzgn/gox/component/interceptor"
	"github.com/yhyzgn/gox/ctx"
)

type WebConfig struct{}

// Context 配置 Context
func (wc *WebConfig) Context(ctx *ctx.GoXContext) {}

// ConfigFilter 注册过滤器
func (wc *WebConfig) ConfigFilter(chain *filter.Chain) {}

// ConfigInterceptor 注册拦截器
func (wc *WebConfig) ConfigInterceptor(register *interceptor.Register) {}
