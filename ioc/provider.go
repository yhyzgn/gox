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
// time   : 2019-10-11 2:27
// version: 1.0.0
// desc   : IOC操作

package ioc

import "reflect"

type Provider struct {
	container *Container
}

func NewProvider() *Provider {
	return &Provider{
		container: NewContainer(),
	}
}

func (p *Provider) Single(name string, bean interface{}) *Provider {
	p.container.SetSingle(name, bean)
	return p
}

func (p *Provider) Prototype(name string, factory factory) *Provider {
	p.container.SetPrototype(name, factory)
	return p
}

func (p *Provider) Get(name string) interface{} {
	return p.container.GetSingle(name)
}

func (p *Provider) GetSingle(name string) interface{} {
	return p.container.GetSingle(name)
}

func (p *Provider) GetPrototype(name string) interface{} {
	iv, _ := p.container.GetPrototype(name)
	return iv
}

func (p *Provider) GetByType(tp reflect.Type) interface{} {
	iv, _ := p.container.GetByTypeSingle(tp)
	return iv
}

func (p *Provider) GetSingleByType(tp reflect.Type) interface{} {
	iv, _ := p.container.GetByTypeSingle(tp)
	return iv
}

func (p *Provider) GetPrototypeByType(tp reflect.Type) interface{} {
	iv, _ := p.container.GetByTypePrototype(tp)
	return iv
}

func (p *Provider) Put(name string, factory factory) *Provider {
	p.container.Put(name, factory)
	return p
}

func (p *Provider) Add(factory factory) *Provider {
	p.container.Add(factory)
	return p
}

func (p *Provider) Inject(instance interface{}) error {
	return p.container.Inject(instance)
}

func (p *Provider) String() string {
	return p.container.String()
}
