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
// time   : 2020-05-13 10:04 下午
// version: 1.0.0
// desc   : Controller基类

package of

import (
	"net/http"

	"github.com/yhyzgn/gox/common"
	"github.com/yhyzgn/gox/util"
)

// Controller 基类
type Controller struct{}

// SetReqAttr 设置request字段
func (c Controller) SetReqAttr(req *http.Request, key string, value interface{}) {
	util.SetRequestAttribute(req, common.AttributeKey(key), value)
}

// GetReqAttr 获取request字段
func (c Controller) GetReqAttr(req *http.Request, key string) interface{} {
	return util.GetRequestAttribute(req, common.AttributeKey(key))
}
