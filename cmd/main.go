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
// time   : 2019-11-25 14:47
// version: 1.0.0
// desc   : 

package main

import (
	"fmt"
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gox"
	"github.com/yhyzgn/gox/cmd/controller"
	"net/http"
)

const (
	port = 8080
)

func main() {
	handler := gox.NewGoX()

	initWeb(handler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	gog.InfoF("GoX 启动于端口 【%d】", port)
	err := server.ListenAndServe()
	if err != nil {
		gog.Error(err)
	}
}

func initWeb(r *gox.GoX) {
	r.Configure(NewConfig())

	r.Mapping("/api/hello", new(controller.HelloController))
}
