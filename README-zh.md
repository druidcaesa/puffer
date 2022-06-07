# Puffer

### 简体中文 [English](./README.md)

#### 一、介绍

puffer 是一款快速、简洁的Go Web框架，使用极少的代码来开发高性能的Web应用程序，框架遵循Go的编程思量，零入侵。

框架专注于web服务处理，其他工具使用组件或者中间件进行集成

#### 二、安装教程

##### 1、go get 安装

```shell
go get -u -v github.com/druidcaesa/puffer
```

#### 2、go mod 安装

```shell
require github.com/druidcaesa/puffer
```

#### 3、go版本显示

```shell
go >=1.15
```

#### 三、使用说明

##### 1、快速开始

- 编写代码

```go
package main

import (
	"github.com/druidcaesa/puffer"
	"net/http"
)

func main() {
	//快速创建启动服务
	server := puffer.New()
	//注册路由,第二个参数是一个函数
	server.GET("/hello", func(c *puffer.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	//监听端口
	server.Run(":8080")
}
```

- 请求结果

```shell
curl "http://127.0.0.1:8080/hello"
hello , you're at /hello
```

- 就这么简单一个简单基础服务搭建完毕

##### 2、路由

###### 1、静态路由

```go

package main

import (
	"github.com/druidcaesa/puffer"
	"net/http"
)

func main() {
	r := puffer.New()
	r.GET("/", func(c *puffer.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Puffer</h1>")
	})

	r.GET("/hello", func(c *puffer.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.Run(":8080")
}
```

###### 2、动态路由

```go
package main

import (
	"github.com/druidcaesa/puffer"
	"net/http"
)

func main() {
	//快速创建启动服务
	server := puffer.New()
	//一级路由
	server.GET("/hello/:name", func(c *puffer.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	//多级路由测试
	server.GET("/assets/*filepath", func(c *puffer.Context) {
		c.JSON(http.StatusOK, puffer.H{"filepath": c.Param("filepath")})
	})
	//监听端口
	server.Run(":8080")
}
```

###### 3、分组路由

```go
package main

import (
	"github.com/druidcaesa/puffer"
)

func main() {
	//快速创建启动服务
	server := puffer.New()
	//正常路由
	server.GET("/login", login)
	//分组路由
	v1 := server.Group("/v1")
	{
		v1.GET("/getInfo", getInfo)
	}
}

func getInfo(context *puffer.Context) {

}

func login(context *puffer.Context) {

}
```

```shell
2022/06/07 16:01:55 Route  GET - /login
2022/06/07 16:01:55 Route  GET - /v1/getInfo
```

##### 3、参数获取

###### 1、GET请求参数

```go

package main

import (
	"fmt"
	"github.com/druidcaesa/puffer"
)

//用户获取get参数，必须使用tag form
type Query struct {
	Name string `form:"name"`
}

func main() {
	//快速创建启动服务
	server := puffer.New()
	//正常路由
	server.GET("/login", login)
	//分组路由
	v1 := server.Group("/v1")
	{
		v1.GET("/getInfo", getInfo)
	}
	server.Run(":8080")
}

func getInfo(c *puffer.Context) {
	q := new(Query)
	_, err := c.BindQuery(q)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("获取前端传的参数%s", q.Name)
}

func login(context *puffer.Context) {

}
```

###### 2、POST请求的body 体的JSON数据

```go
package main

import (
	"fmt"
	"github.com/druidcaesa/puffer"
	"net/http"
)

//用户获取get参数，必须使用tag form
type Query struct {
	Name string `form:"name"`
}

type JSON struct {
	UserName string `json:"userName"`
}

func main() {
	//快速创建启动服务
	server := puffer.New()
	//正常路由
	server.POST("/login", login)
	//分组路由
	v1 := server.Group("/v1")
	{
		v1.GET("/getInfo", getInfo)
	}
	server.Run(":8080")
}

func getInfo(c *puffer.Context) {
	q := new(Query)
	_, err := c.BindQuery(q)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("获取前端传的参数%s", q.Name)
}

func login(c *puffer.Context) {
	j := new(JSON)
	_, err := c.BindJsonBody(j)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("获取前端传的参数%s", j.UserName)
	c.JSON(http.StatusOK, puffer.H{
		"data": j,
	})
}

```
![请求结果](http://39.105.57.46:9000/cloud-disk/WX20220607-164541@2x.png)


##### 4、中间件

##### 1、创建中间件

```go
/**
  @author: fanyanan
  @date: 2022/6/7
  @note: //日志中间件
**/
package middlewares

import (
	"github.com/druidcaesa/puffer"
	"log"
	"time"
)

func Logger() puffer.HandlerFunc {
	return func(c *puffer.Context) {
		//start time
		t := time.Now()
		//Middleware processing logic
		log.Printf("[Status:%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
		//Pass to the next middleware
		c.Next()
	}
}
```

##### 2、中间件注册

```go
package main

import (
	"fmt"
	"github.com/druidcaesa/puffer"
	"net/http"
	"web-demo/middlewares"
)

//用户获取get参数，必须使用tag form
type Query struct {
	Name string `form:"name"`
}

type JSON struct {
	UserName string `json:"userName"`
}

func main() {
	//快速创建启动服务
	server := puffer.New()
	//注册全局中间件
	server.Use(middlewares.Logger())
	//分组路由
	v1 := server.Group("/v1")
	{
		v1.GET("/getInfo", getInfo)
	}
	server.Run(":8080")
}

func getInfo(c *puffer.Context) {
	q := new(Query)
	_, err := c.BindQuery(q)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("获取前端传的参数%s", q.Name)
}
```
![请求结果](http://39.105.57.46:9000/cloud-disk/WX20220607-170843@2x.png)

