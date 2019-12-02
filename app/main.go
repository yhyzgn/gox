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
	"context"
	"fmt"
	"github.com/yhyzgn/gog"
	"github.com/yhyzgn/gog/level"
	"github.com/yhyzgn/gog/writer"
	"github.com/yhyzgn/gox"
	"github.com/yhyzgn/gox/app/api/controller"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	port = 8080
)

func init() {
	env := os.Getenv("ENV")
	if env == "prod" {
		// 生产环境
		gog.AddWriter(writer.NewJSONWriter())
		gog.Level(level.INFO)
	} else {
		// 开发环境
		gog.AddWriter(writer.NewConsoleWriter())
	}
	gog.Async(true)
}

func main() {
	handler := gox.NewGoX()

	initWeb(handler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit

		gog.Info("收到关闭信号，正在关闭 GoX 服务 ...")
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			gog.ErrorF("GoX 关闭错误 【{}】", err)
		}
	}()

	gog.InfoF("GoX 启动于端口 【%d】", port)
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			gog.Info("GoX 服务已安全退出！")
			return
		}
		gog.Error(err)
	}
}

func initWeb(r *gox.GoX) {
	r.Configure(NewConfig())

	r.Mapping("/api/normal", new(controller.NormalController))
	r.Mapping("/api/param", new(controller.ParamController))
}
