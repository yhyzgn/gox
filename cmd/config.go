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
// time   : 2019-11-26 15:10
// version: 1.0.0
// desc   : 

package main

import (
	testFilter "github.com/yhyzgn/gox/cmd/filter"
	testInterceptor "github.com/yhyzgn/gox/cmd/interceptor"
	"github.com/yhyzgn/gox/component/filter"
	"github.com/yhyzgn/gox/component/interceptor"
)

type Config struct {
}

func NewConfig() *Config {
	return new(Config)
}

func (c *Config) ConfigFilter(chain *filter.FilterChain) {
	chain.
		AddFilter("/", testFilter.NewTestFilter()).
		AddFilter("/api/*", testFilter.NewLogFilter())
}

func (c *Config) ConfigInterceptor(register *interceptor.InterceptorRegister) {
	register.AddInterceptor("/", testInterceptor.NewTestInterceptor())
}
