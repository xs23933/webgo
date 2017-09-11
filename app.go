package webgo

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	WebGo       *App
	workPath    string
	appConfPath string
	appConf     AppConf
)

// App Class
type App struct {
	Router *Router
	Server *http.Server
}

type AppConf struct {
	Host       string
	StaticPath string
}

func init() {
	workPath, _ = os.Getwd()
	appConfPath = filepath.Join(workPath, "conf", "config.ini")
	if _, err := toml.DecodeFile(appConfPath, &appConf); err != nil {
		panic(err)
	}

	WebGo = NewApp()
}
func NewApp() *App {
	fmt.Println("NewApp")
	app := &App{Router: NewRouter(), Server: &http.Server{}}
	return app
}

func (app *App) Run() {
	fmt.Println("App.Run")
	end := make(chan bool, 1)

	app.Server.Addr = appConf.Host
	app.Server.Handler = app.Router
	go func() {
		if err := app.Server.ListenAndServe(); err != nil {
			end <- true
		}
	}()
	<-end
}

func Run() {
	fmt.Println("Run")
	WebGo.Run()
}

// Router 路由部分
type Handlers struct {
	pattern        string
	controllerType reflect.Type
}

type Router struct {
	handlers []*Handlers
}

func NewRouter() *Router {
	fmt.Println("NewRouter")
	router := &Router{handlers: make([]*Handlers, 0)}
	return router
}

func (p *Router) Add(pattern string, c ControllerInterface) {
	fmt.Println("Router.Add", pattern, c)

	t := reflect.Indirect(reflect.ValueOf(c)).Type()

	route := &Handlers{}
	route.pattern = pattern
	route.controllerType = t
	p.handlers = append(p.handlers, route)
}

func (p *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// static route must include staticPath keywork
	rw.Header().Set("Server", "s")
	if strings.HasPrefix(r.URL.Path[1:], appConf.StaticPath) {
		file := filepath.Join(workPath, r.URL.Path)
		http.ServeFile(rw, r, file)
		return
	}
	var runRouter *Handlers
	for _, route := range p.handlers {
		if route.pattern == r.URL.Path {
			runRouter = route
		}
		fmt.Println(route.pattern, route.controllerType)
	}

	// 找到了注册的路由
	vc := reflect.New(runRouter.controllerType)
	init := vc.MethodByName("Init")
	ct := &Context{ResponseWriter: rw, Request: r}
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(ct)
	in[1] = reflect.ValueOf(runRouter.controllerType.Name())
	init.Call(in)
	// Prepare
	in = make([]reflect.Value, 0)
	method := vc.MethodByName("Prepare")
	method.Call(in)
	// Rquest method
	rMethod := strings.ToUpper(r.Method[:1]) + strings.ToLower(r.Method[1:])
	method = vc.MethodByName(rMethod)
	method.Call(in)

	method = vc.MethodByName("Destructor")
	method.Call(in)
	fmt.Println("Router.ServeHTTP", r.URL.Path, r.Method)

}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type Controller struct {
	ctx       *Context
	ChildName string
}
type ControllerInterface interface {
	Init(ctx *Context, cn string)
	Get()
	Post()
	Head()
	Put()
	Prepare()
}

func (c *Controller) Init(ctx *Context, cn string) {
	c.ctx = ctx
	c.ChildName = cn
}

func (c *Controller) Write(content string) {
	c.ctx.ResponseWriter.Write([]byte(content))
}

func (c *Controller) Get()        {}
func (c *Controller) Post()       {}
func (c *Controller) Head()       {}
func (c *Controller) Put()        {}
func (c *Controller) Prepare()    {}
func (c *Controller) Destructor() {}

func Route(pattern string, c ControllerInterface) {
	fmt.Println("Route")
	WebGo.Router.Add(pattern, c)
}
