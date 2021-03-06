package puffer

import (
	"encoding/json"
	"fmt"
	"github.com/druidcaesa/puffer/cookie"
	"github.com/druidcaesa/puffer/file"
	"github.com/druidcaesa/puffer/utils"
	"log"
	"net/http"
)

type H map[string]interface{}

type ContextFunc interface {
	// BindQuery Get Get request parameters binding struct
	BindQuery(v interface{}) (bool, error)
	// BindJsonBody Get POST request parameters
	BindJsonBody(v interface{}) (bool, error)
	// GetCookie Get cookie
	GetCookie(key string) (*http.Cookie, error)
	// SetCookie set cookies
	SetCookie(key, value, path, domain string, maxAge int, secure, httpOnly bool)
	// GetQuery Get Get request parameters
	GetQuery(key string) string
	// GetParameter Get dynamic route URI parameters
	GetParameter(key string) string
	// GetPostForm Get form submission data
	GetPostForm(key string) string
	// FormFile Get File Upload
	FormFile(fileKey string) (*file.File, error)
}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	engine   *Engine
	Cookie   cookie.Cookie
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

// BindQuery Get request parameter binding
func (c *Context) BindQuery(v interface{}) (bool, error) {
	u := new(utils.Tag)
	u.R = c.Req
	return u.BindForm(v)
}

// BindJsonBody JSON parameter binding for POST
func (c *Context) BindJsonBody(v interface{}) (bool, error) {
	u := new(utils.Tag)
	u.R = c.Req
	return u.BindJson(v)
}

// GetCookie get cookie
func (c *Context) GetCookie(key string) (*http.Cookie, error) {
	c.Cookie.SetReq(c.Req)
	return c.Cookie.Cookie(key)
}

/**
 * @author fanyanan
 * @description //GetQuery Get Get request parameters
 * @date 16:35 2022/6/7
 * @param key parameter key
 * @return string
 **/
func (c *Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

/**
 * @author fanyanan
 * @description //Get dynamic request parameters
 * @date 16:37 2022/6/7
 * @param Dynamic request parameter key
 * @return string
 **/
func (c *Context) GetParameter(key string) string {
	return c.Params[key]
}

/**
 * @author fanyanan
 * @description set cookie function
 * @date 14:11 2022/6/7
 * @param key cookie key
 * @param value cookie value
 * @param domain domain name
 * @param maxAge Maximum aging unit second
 * @param secure Can it be accessed via https
 * @param httpOnly Whether to allow js to get
 **/
func (c *Context) SetCookie(key, value, path, domain string, maxAge int, secure, httpOnly bool) {
	c.Cookie.SetResp(c.Writer)
	c.Cookie.SetCookie(key, value, path, domain, maxAge, secure, httpOnly)
}

// GetPostForm get post form request parameters
// params key request parameters key
func (c *Context) GetPostForm(key string) string {
	r := c.Req
	err := r.ParseForm()
	if err != nil {
		log.Fatalf("An exception occurs----->%s", err.Error())
		return ""
	}
	for k, v := range r.Form {
		if k == key {
			return v[0]
		}
	}
	return ""
}

// FormFile Get Upload File
func (c *Context) FormFile(fileKey string) (*file.File, error) {
	f := new(file.File)
	f.SetReq(c.Req)
	formFile, err := f.GetFormFile(fileKey)
	if err != nil {
		return nil, err
	}
	return formFile, nil
}
