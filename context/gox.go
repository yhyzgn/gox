// Copyright 2019 yhyzgn xgo
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
// time   : 2019-11-27 9:41
// version: 1.0.0
// desc   : 

package context

import (
	"github.com/yhyzgn/ghost/config"
	"sync"
)

type GoXContext struct {
	reader     *config.Reader
	components map[string]interface{}
}

var (
	once    sync.Once
	current *GoXContext
)

func init() {
	once.Do(func() {
		current = &GoXContext{
			reader:     config.NewReader(),
			components: make(map[string]interface{}),
		}
	})
}

func Current() *GoXContext {
	return current
}

func (c *GoXContext) Read(filename string) (data []byte, errs error) {
	return c.reader.Read(filename)
}

func (c *GoXContext) SetComponent(name string, component interface{}) ComponentContext {
	c.components[name] = component
	return c
}

func (c *GoXContext) GetComponent(name string) interface{} {
	return c.components[name]
}
