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
// time   : 2019-10-11 3:20
// version: 1.0.0
// desc   : 配置文件读取

package resource

import (
	"bytes"
	"errors"
	"github.com/yhyzgn/gox/util"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Reader 配置文件读取器
type Reader struct {
}

// NewReader 读取器实例
func NewReader() *Reader {
	return &Reader{}
}

// Read 读取配置文件
func (cp *Reader) Read(filename string) (data []byte, errs error) {
	if filename == "" {
		errs = errors.New("filename can not be empty")
		return
	}
	if !util.FileExist(filename) {
		errs = errors.New("no such config file '" + filename + "'")
		return
	}
	data, errs = ioutil.ReadFile(filename)
	return
}

// Decode 解析配置文件
func (cp *Reader) Load(filename string, bean interface{}) error {
	bs, err := cp.Read(filename)
	if err != nil {
		return err
	}
	if isYaml(filename) {
		return yaml.NewDecoder(bytes.NewBuffer(bs)).Decode(bean)
	} else if isToml(filename) {
		_, err := toml.Decode(string(bs), bean)
		return err
	} else {
		return errors.New("unknown config file '" + filename + "'")
	}
}

func isYaml(filename string) bool {
	return util.EndWith(filename, ".yml") || util.EndWith(filename, ".yaml")
}

func isToml(filename string) bool {
	return util.EndWith(filename, ".toml")
}
