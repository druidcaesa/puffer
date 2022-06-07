# Puffer

### English [简体中文](./README-zh.md)

#### 一、introduce

puffer It is a fast and concise Go web framework that uses very little code to develop high-performance web applications. The framework follows Go's programming thinking and has zero intrusion.

The framework focuses on web service processing, other tools use components or middleware for integration

#### 二、Installation tutorial

##### 1、go get Install

```shell
go get -u -v github.com/druidcaesa/puffer
```

#### 2、go mod Install

```shell
require github.com/druidcaesa/puffer
```

#### 3、go version display

```shell
go >=1.15
```

#### 三、go version display

##### 1、go version display

- go version display

```go
package main

import (
	"github.com/druidcaesa/puffer"
	"net/http"
)

func main() {
	//go version display
	server := puffer.New()
	//Register the route, the second parameter is a function
	server.GET("/hello", func(c *puffer.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	//listening port
	server.Run(":8080")
}
```

- request result

```shell
curl "http://127.0.0.1:8080/hello"
hello , you're at /hello
```

- It's that simple, a simple basic service is built

##### 2、routing

###### 1、static routing

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

###### 2、dynamic routing

```go
package main

import (
	"github.com/druidcaesa/puffer"
	"net/http"
)

func main() {
	//Quickly create a startup service
	server := puffer.New()
	//first-level routing
	server.GET("/hello/:name", func(c *puffer.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	//Multi-level routing test
	server.GET("/assets/*filepath", func(c *puffer.Context) {
		c.JSON(http.StatusOK, puffer.H{"filepath": c.Param("filepath")})
	})
	//listening port
	server.Run(":8080")
}
```

###### 3、packet routing

```go
package main

import (
	"github.com/druidcaesa/puffer"
)

func main() {
	//Quickly create a startup service
	server := puffer.New()
	//normal route
	server.GET("/login", login)
	//packet routing
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

##### 3、parameter acquisition

###### 1、GET request parameters

```go

package main

import (
	"fmt"
	"github.com/druidcaesa/puffer"
)

//The user gets the get parameter and must use the tag form
type Query struct {
	Name string `form:"name"`
}

func main() {
	//Quickly create a startup service
	server := puffer.New()
	//normal route
	server.GET("/login", login)
	//packet routing
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
	fmt.Printf("Get the parameters passed by the front end%s", q.Name)
}

func login(context *puffer.Context) {

}
```

###### 2、JSON data of the body of the POST request

```go
package main

import (
	"fmt"
	"github.com/druidcaesa/puffer"
	"net/http"
)

//The user gets the get parameter and must use the tag form
type Query struct {
	Name string `form:"name"`
}

type JSON struct {
	UserName string `json:"userName"`
}

func main() {
	//Quickly create a startup service
	server := puffer.New()
	//normal route
	server.POST("/login", login)
	//packet routing
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
	fmt.Printf("Get the parameters passed by the front end%s", q.Name)
}

func login(c *puffer.Context) {
	j := new(JSON)
	_, err := c.BindJsonBody(j)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Get the parameters passed by the front end%s", j.UserName)
	c.JSON(http.StatusOK, puffer.H{
		"data": j,
	})
}

```
![请求结果](http://39.105.57.46:9000/cloud-disk/WX20220607-164541@2x.png)


##### 4、middleware

##### 1、Create middleware

```go
/**
  @author: fanyanan
  @date: 2022/6/7
  @note: //log middleware
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

##### 2、Middleware registration

```go
package main

import (
	"fmt"
	"github.com/druidcaesa/puffer"
	"net/http"
	"web-demo/middlewares"
)

//The user gets the get parameter and must use the tag form
type Query struct {
	Name string `form:"name"`
}

type JSON struct {
	UserName string `json:"userName"`
}

func main() {
	//Quickly create a startup service
	server := puffer.New()
	//Register global middleware
	server.Use(middlewares.Logger())
	//packet routing
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
	fmt.Printf("Get the parameters passed by the front end%s", q.Name)
}
```
![请求结果](http://39.105.57.46:9000/cloud-disk/WX20220607-170843@2x.png)

