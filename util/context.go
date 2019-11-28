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
// time   : 2019-11-26 14:46
// version: 1.0.0
// desc   : 

package util

import (
	"github.com/yhyzgn/gox/context"
)

func GetWare(name string, defValue interface{}) interface{} {
	instance := context.Current().GetWare(name)
	if instance == nil {
		instance = defValue
	}
	context.Current().SetWare(name, instance)
	return instance
}
