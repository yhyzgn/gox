# gox

![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/yhyzgn/gox?color=orange&label=size) ![](https://img.shields.io/github/v/release/yhyzgn/gox?color=brightgreen)

> 基于`Golang`实现的`MVC`后台开发框架。

​		与`gin`、`beego`等框架定位不同，`gox`借鉴了`Spring MVC`框架设计思想，将`http`请求的请求地址、请求参数、参数处理、逻辑处理、响应参数等操作分离。

​		一般情况下，开发者只需要关注自己的“逻辑处理”模块，也就是具体`controller`下的方法实现，无需手动从`request`中获取参数；逻辑处理完成后，将处理结果直接返回即可，也无需手动`write`到`response`中。

​		已支持普通`http`请求，也支持`单/多文件上传`等操作

​		不过，由于 `golang`没有注解能力，反射功能也远不及`java`强大，所以`handler`注册只能手动注册完成



[TOC]

# 0、功能提点

* 普通`http`请求处理
* 自动提取参数并处理
* 结果处理器处理响应结果
* 支持单/多文件上传
* 支持`RESTful`请求
* 支持过滤器，按`path`配置`Filter`
* 支持拦截器，按`path`配置`Interceptor`
* 可配置`Context Path`
* 支持静态资源请求
* 可自由配置参数处理器、结果处理器和异常处理器
* s



# 1、基本用法

> **简单示例：[`gox-app`](https://github.com/yhyzgn/gox-app)**
>
> **项目样例：[`gox-simple`](https://github.com/yhyzgn/gox-simple)**

## 1.1、引入`gox`

* 方式一 `go mod`

  ```shell
  go mod init # 如果是go mod项目，则忽略此命令
  go get github.com/yhyzgn/gox
  ```

* 方式二 `go get`

  ```shell
  go get -u github.com/yhyzgn/gox
  ```

  

##  1.2、基本用法

### 1.2.1、服务入口配置

> 项目唯一入口，即主程序`main`方法中

```go
x := gox.NewGoX()

server := &http.Server{
    Addr: fmt.Sprintf(":%d", 8888),
}

x.Run(server)
```

这样就启动了运行于`8888`端口的`http`服务。不过只是运行了服务而已，无法处理任何请求，因为未配置任何`controller`。



### 1.2.2、`controller`实现

> 由于`golang`没有注解可用，`gox`只能靠手动注册`controller`实现，故**所有`controller`都需要实现`core.Controller`接口**。

`core.Controller`源码：

```go
// Controller 控制器接口
type Controller interface {
    // Mapping 在该方法中完成当前控制器中的所有处理器注册
    //
    // mapper.Request("/rest/{id}").HandlerFunc(h.Rest).PathVariable("id").Mapping()
    Mapping(*Mapper)
}
```

`NormalController`实现：

```go
package controller

import (
    "github.com/yhyzgn/gox/core"
    "net/http"
)

type NormalController struct {
}

// 实现 core.Controller.Mapping 方法，在此完成 请求处理器的注册
func (c NormalController) Mapping(mapper *core.Mapper) {
    mapper.
    Request("/").HandlerFunc(c.Default).Mapping().
    Get("/get").HandlerFunc(c.Get).Mapping().
    Post("/post").HandlerFunc(c.Post).Mapping().
    Delete("/delete").HandlerFunc(c.Delete).Mapping().
    Request("/multi").HandlerFunc(c.Multi).Method(http.MethodGet, http.MethodPost).Mapping()
}

func (c NormalController) Default() string {
    return c.res("Default")
}

func (c NormalController) Get() string {
    return c.res("Get")
}

func (c NormalController) Post() string {
    return c.res("Post")
}

func (c NormalController) Delete() string {
    return c.res("Delete")
}

func (c NormalController) Multi() string {
    return c.res("Multi")
}

func (c NormalController) res(str string) string {
    return "GoX Normal " + str
}
```



### 1.2.3、配置`controller`到`gox`

> 重新配置服务入口，在`x.Run(server)`前配置`controller`即可

```go
// ...

// 配置 controller 映射
x.
Mapping("/api/normal", new(controller.NormalController)).
Mapping("/api/user", new(controller.UserController))

// 启动服务
x.Run(server)
```

重新启动后，访问`http://localhost:8888/api/normal`，如成功响应`GoX Normal Default`，则说明服务已经运行成功并可用。



## 1.3、高级用法

### 1.3.1、参数姿势

> `http`请求参数接收，包括文件上传

需要在实现`core.Controller.Mapping(mapper *core.Mapper)`方法中完成参数注册

```go
package controller

import (
    "fmt"
    "github.com/yhyzgn/gog"
    "github.com/yhyzgn/gox/common"
    "github.com/yhyzgn/gox/core"
    "github.com/yhyzgn/gox/util"
    "io/ioutil"
    "net/http"
)

type ParamController struct {
}

func (c ParamController) Mapping(mapper *core.Mapper) {
    // 原生方式获取参数
    mapper.Get("/native").HandlerFunc(c.Native).Mapping()
    mapper.Get("/noReturn").HandlerFunc(c.NoReturn).Mapping()
    // name 和 age 参数都是必传的，缺失则报错
    mapper.Get("/required").HandlerFunc(c.Required).Required("name").Required("age").Mapping()
    // name 为必传，age 为可选
    mapper.Get("/param").HandlerFunc(c.Param).Required("name").Param("age").Mapping()
    // 从 request header 中获取名为 Token 的参数
    mapper.Get("/header").HandlerFunc(c.Header).Header("Token").Param("rand").Mapping()
    // RESTful 格式传参
 mapper.Get("/rest/{name}/{age}/test").HandlerFunc(c.Param).PathVariable("name").PathVariable("age").Mapping()
    // request body 传参
    mapper.Post("/body").HandlerFunc(c.Body).Body("user").Method(http.MethodPut).Mapping()
    // 普通 表单提交参数，用对象接收
    mapper.Post("/vo").HandlerFunc(c.VO).Param("std").Mapping()
    // 单文件上传
    mapper.Post("/singleFile").HandlerFunc(c.SingleFile).Param("file").Mapping()
    // 多文件上传
    mapper.Post("/multiFiles").HandlerFunc(c.MultiFiles).Param("files").Mapping()
}

// 如果你坚持要从 request 中手动获取
func (c ParamController) Native(writer http.ResponseWriter, request *http.Request) string {
    user := request.URL.Query().Get("param")
    request.Header.Set("X-Env", "dev")
    writer.Header().Add("User", user)
    return c.res("Native " + user)
}

// 无返回值
func (c ParamController) NoReturn(writer http.ResponseWriter, request *http.Request) {
    util.ResponseJSON(writer, c.res("NoReturn "+request.URL.Query().Get("param")))
}

// 必须参数
func (c ParamController) Required(name string, age int) string {
    return c.res(fmt.Sprintf("Required name = %v, age = %d", name, age))
}

// 可空参数
func (c ParamController) Param(name string, age int) string {
    return c.res(fmt.Sprintf("Param name = %v, age = %d", name, age))
}

// 从Header取参数
func (c ParamController) Header(token string, rand int) string {
    return c.res(fmt.Sprintf("Header token = %v, rand = %d", token, rand))
}

// RESTful 参数
func (c ParamController) Rest(name string, age int) string {
    return c.res(fmt.Sprintf("Rest name = %v, age = %d", name, age))
}

// 接收 RequestBody，如 json
func (c ParamController) Body(user *User) *User {
    return user
}

// 普通表单参数转为 对象
func (c ParamController) VO(sdt *Student) *Student {
    return sdt
}

// 单文件上传
func (c ParamController) SingleFile(file *common.MultipartFile) string {
    bs, err := ioutil.ReadAll(file.File)
    if err != nil {
        gog.Error(err)
        return "文件上传失败"
    }
    gog.Info(string(bs))
    return c.res("SingleFile filename = " + file.Header.Filename)
}

// 多文件上传
func (c ParamController) MultiFiles(files []*common.MultipartFile) string {
    if files == nil || len(files) == 0 {
        return "未接收到任何文件"
    }
    gog.DebugF("接收到【{}】个文件", len(files))
    for i, file := range files {
        gog.DebugF("文件【{}】名称为【{}】", i, file.Header.Filename)
    }
    return c.res("MultipartFile " + fmt.Sprintf("接收到【%d】个文件", len(files)))
}

func (c ParamController) res(str string) string {
    return "GoX Param " + str
}

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type Student struct {
    ID    int `param:"id"`
    Name  string
    Age   int    `param:"age" required:""`
    Token string `header:""`
}
```



### 1.3.2、自定义配置

> 如`ContextPath`、`Filter`、`Interceptor`等，需要通过`configure.WebConfigure`接口来实现

#### 1.3.2.1、 实现`configure.WebConfigure`接口

也提供了默认实现类`adapter.WebConfig`适配器，继承并重写具体方法即可。

以下是实现样例

```go
package web

type WebConfig struct {
    adapter.WebConfig // 继承 adapter.WebConfig 类
}

func (wc *WebConfig) Context(ctx *ctx.GoXContext) {
    ctx.
    // 404 处理器
    SetNotFoundHandler(func(writer http.ResponseWriter, request *http.Request) {
        http.Error(writer, "服务器度假去啦~", http.StatusNotFound)
    })
}

// 配置一些 过滤器
func (wc *WebConfig) ConfigFilter(chain *filter.Chain) {
    chain.AddFilters("/**", filters.NewBuiltFilter(), filters.NewLogFilter())
}

// 配置一些 拦截器
func (wc *WebConfig) ConfigInterceptor(register *interceptor.Register) {
    register.AddInterceptors("/api/**", interceptors.NewAuthInterceptor())
}
```



#### 1.3.2.2、将配置应用到`gox`

> 在 `x.Run(server)`前配置

```go
// ...

// 配置 Context Path
x.SetContextPath("/backend")

// 应用配置
x.Configure(new(web.WebConfig))

x.Run(server)
```



# 2、定制化配置

## 2.1、统一异常处理器

> 实现`resolver.ErrorResolver`接口，并在`configure.WebConfigure.Context()`实现中配置

* 内置实现

  ```go
  // SimpleErrorResolver 默认的异常处理器
  type SimpleErrorResolver struct{}
  
  func NewSimpleErrorResolver() *SimpleErrorResolver {
      return new(SimpleErrorResolver)
  }
  
  func (ser *SimpleErrorResolver) Resolve(err error, writer http.ResponseWriter) interface{} {
      return err
  }
  ```

* 具体处理器

  ```go
  type ErrorResolver struct{}
  
  func NewErrorResolver() *ErrorResolver {
      return new(ErrorResolver)
  }
  
  func (er *ErrorResolver) Resolve(value reflect.Value, writer http.ResponseWriter) interface{} {
      err := value.Interface()
      if err == nil {
          return res.Failed("未知错误")
      }
      // 默认为error类型
      return res.Failed(err.(error).Error())
  }
  ```

* 配置应用

  ```go
  func (wc *WebConfig) Context(ctx *ctx.GoXContext) {
      ctx.
      // 统一异常处理器
      SetWare(common.ErrorResolverName, resolver.NewErrorResolver())
  }
  ```

  

## 2.2、参数处理器

> 实现`resolver.ArgumentResolver`接口

* 内置实现

  ```go
  // SimpleArgumentResolver 默认的参数处理器
  type SimpleArgumentResolver struct {
  }
  
  // NewSimpleArgumentResolver 创建新的参数处理器对象
  func NewSimpleArgumentResolver() *SimpleArgumentResolver {
      return new(SimpleArgumentResolver)
  }
  
  // Resolve 处理参数操作
  func (sar *SimpleArgumentResolver) Resolve(args []reflect.Value, writer http.ResponseWriter, request *http.Request, handler common.Handler) []reflect.Value {
      // 这里接收到的参数列表已经过初步处理
      // 默认直接返回使用
      return args
  }
  ```



## 2.3、结果处理器

> 实现`resolver.ResultResolver`接口

* 内置实现

  ```go
  // SimpleResultResolver 默认的结果处理器
  type SimpleResultResolver struct {
  }
  
  // NewSimpleResultResolver 创建新的结果处理器对象
  func NewSimpleResultResolver() *SimpleResultResolver {
      return new(SimpleResultResolver)
  }
  
  // Resolve 处理结果集
  // 只接受最多两个返回值的结果集
  // 如果不满足需求，可自定义
  func (srr *SimpleResultResolver) Resolve(hw *wire.HandlerWire, values []reflect.Value, writer http.ResponseWriter, request *http.Request) (value reflect.Value, err error) {
      path := request.URL.Path
      handler := reflect.Value(hw.Handler)
      pc := handler.Pointer()
      handlerName := strings.ReplaceAll(runtime.FuncForPC(pc).Name(), "-fm", "(...)")
  
      if values == nil || len(values) == 0 {
          // 没有返回值，无需处理
          return
      }
  
      ln := len(values)
      // 只有1个返回值，必定是 请求响应结果
      if ln == 1 {
          //srr.Response(values[0], writer)
          value = values[0]
          return
      }
  
      // 结果1：请求响应结果
      // 结果2：错误信息
      if ln == 2 {
          if e := values[1]; e.Interface() != nil {
              err = e.Interface().(error)
              return
          }
          value = values[0]
          return
      }
  
      // 结果不能超过2个
      err = fmt.Errorf("the path [%v] handled [%v] support 2 results at most, but now is [%d]", path, handlerName, ln)
      return
  }
  
  // Response 响应结果
  func (srr *SimpleResultResolver) Response(value reflect.Value, writer http.ResponseWriter) {
      util.ResponseJSON(writer, value.Interface())
  }
  ```

  