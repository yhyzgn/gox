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
// time   : 2019-10-11 1:55
// version: 1.0.0
// desc   : IOC容器

package ioc

import (
	"errors"
	"fmt"
	"github.com/yhyzgn/gox/util"
	"reflect"
	"strings"
	"sync"
)

var (
	errorPrototype   = errors.New("ioc prototype factory not found")
	errorInjectPtr   = errors.New("inject instance must be struct pointer type")
	errorInjectValid = errors.New("inject instance must be struct type")
)

type factory func() (instance interface{})

type Container struct {
	sync.Mutex
	singles    map[string]interface{}
	prototypes map[string]factory
}

func NewContainer() *Container {
	return &Container{
		singles:    make(map[string]interface{}),
		prototypes: make(map[string]factory),
	}
}

// 添加实例工厂，自动从工厂获取单例保存
func (c *Container) Put(name string, factory factory) {
	instance := factory()
	// 如果名称为空，就按类型保存
	if name == "" {
		name = reflect.TypeOf(instance).String()
	}
	// 存入单例
	c.SetSingle(name, instance)
	// 存入原型
	c.SetPrototype(name, factory)
}

// 添加实例工厂，自动从工厂获取单例保存
func (c *Container) Add(factory factory) {
	c.Put("", factory)
}

func (c *Container) SetSingle(name string, bean interface{}) {
	c.Lock()
	// 如果名称为空，就按类型保存
	tp := reflect.TypeOf(bean)
	if name == "" {
		name = tp.String()
	}
	// 如果是指针类型的类，就自动依赖注入字段
	_, isStruct, isPtr := util.StructType(bean)
	if isStruct && isPtr {
		_ = c.Inject(bean)
	}
	c.singles[name] = bean
	c.Unlock()
}

func (c *Container) GetSingle(name string) interface{} {
	return c.singles[name]
}

func (c *Container) SetPrototype(name string, factory factory) {
	c.Lock()
	// 如果名称为空，就按类型保存
	if name == "" {
		instance := factory()
		name = reflect.TypeOf(instance).String()
	}
	c.prototypes[name] = factory
	c.Unlock()
}

func (c *Container) GetPrototype(name string) (interface{}, error) {
	factory, ok := c.prototypes[name]
	if ok {
		return c.factoryInject(factory)
	}
	return nil, errorPrototype
}

func (c *Container) GetByTypeSingle(tp reflect.Type) (interface{}, error) {
	// 先查找注册为空名称的bean
	if c.singles[tp.String()] != nil {
		return c.singles[tp.String()], nil
	}
	// 查找单例
	for _, item := range c.singles {
		if tp == reflect.TypeOf(item) {
			return item, nil
		}
	}
	return nil, errors.New("ioc type '" + tp.String() + "' dependency not found")
}

func (c *Container) GetByTypePrototype(tp reflect.Type) (interface{}, error) {
	// 先查找注册为空名称的bean
	factory := c.prototypes[tp.String()]
	if factory != nil {
		return c.factoryInject(factory)
	}
	// 查找原型
	for _, item := range c.prototypes {
		temp, err := c.factoryInject(item)
		if err == nil && tp == reflect.TypeOf(temp) {
			return temp, nil
		}
	}
	return nil, errors.New("ioc type '" + tp.String() + "' dependency not found")
}

func (c *Container) Inject(instance interface{}) error {
	elemType := reflect.TypeOf(instance)
	// 注入对象必须是指针类型，否则将会注入失败
	if elemType.Kind() != reflect.Ptr {
		return errorInjectPtr
	}
	elemType = elemType.Elem()
	// 依赖注入的对象必须是Struct类型
	if elemType.Kind() != reflect.Struct {
		return errorInjectValid
	}
	elemValue := reflect.ValueOf(instance).Elem()
	for i := 0; i < elemType.NumField(); i++ {
		fieldType := elemType.Field(i)
		auto, ok := fieldType.Tag.Lookup("auto")
		if !ok {
			continue
		}
		var (
			iocInstance interface{}
			err         error
		)
		scope := fieldType.Tag.Get("scope")
		if scope == "prototype" {
			// 原型模式
			if auto == "" {
				// 按类型查找
				iocInstance, err = c.GetByTypePrototype(fieldType.Type)
			} else {
				iocInstance, err = c.GetPrototype(auto)
			}
		} else {
			// 默认是单例模式
			if auto == "" {
				// 按类型查找
				iocInstance, err = c.GetByTypeSingle(fieldType.Type)
			} else {
				iocInstance = c.GetSingle(auto)
			}
		}
		if err != nil {
			return err
		}
		if iocInstance == nil {
			return errors.New("ioc field '" + auto + "' dependency not found")
		}
		// 设置字段值
		util.FieldSet(elemValue.Field(i), reflect.ValueOf(iocInstance))
	}
	return nil
}

func (c *Container) String() string {
	lines := make([]string, 0, len(c.singles)+len(c.prototypes)+2)
	lines = append(lines, "singles:")
	for name, item := range c.singles {
		addr := item
		// 如果不是指针，就取出其地址
		if !util.IsPtr(reflect.TypeOf(item)) {
			addr = &item
		}
		line := fmt.Sprintf("  %s: %p %s", name, addr, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	lines = append(lines, "prototypes:")
	for name, item := range c.prototypes {
		instance := item()
		addr := instance
		// 如果不是指针，就取出其地址
		if !util.IsPtr(reflect.TypeOf(item)) {
			addr = &instance
		}
		line := fmt.Sprintf("  %s: %p %s", name, addr, reflect.TypeOf(instance).String())
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (c *Container) factoryInject(factory factory) (interface{}, error) {
	bean := factory()
	// 如果是指针类型的类，就自动依赖注入字段
	_, isStruct, isPtr := util.StructType(bean)
	if isStruct && isPtr {
		err := c.Inject(bean)
		return bean, err
	}
	return bean, nil
}
