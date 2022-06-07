package puffer

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) handlerIsNil() HandlerFunc {
	return func(context *Context) {
		context.JSON(http.StatusOK, H{
			"message": "操作成功",
		})
	}
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler ...HandlerFunc) {
	if len(handler) == 0 {
		group.addRoute("GET", pattern, group.handlerIsNil())
		return
	}
	group.addRoute("GET", pattern, handler[0])
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler ...HandlerFunc) {
	if len(handler) == 0 {
		group.addRoute("POST", pattern, group.handlerIsNil())
		return
	}
	group.addRoute("POST", pattern, handler[0])
}

// PUT defines the method to add PUT request
func (group *RouterGroup) PUT(pattern string, handler ...HandlerFunc) {
	if len(handler) == 0 {
		group.addRoute("PUT", pattern, group.handlerIsNil())
		return
	}
	group.addRoute("PUT", pattern, handler[0])
}

// DELETE defines the method to add DELETE request
func (group *RouterGroup) DELETE(pattern string, handler ...HandlerFunc) {
	if len(handler) == 0 {
		group.addRoute("DELETE", pattern, group.handlerIsNil())
		return
	}
	group.addRoute("DELETE", pattern, handler[0])
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//Default Default instance
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
