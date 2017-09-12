package webgo

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
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
		Error("配置文件读取失败", err)
	}
	D("D")
	Info("Info")
	Error("Error")
	Warning("Warning")
	Success("Success")
	WebGo = NewApp()
}
func NewApp() *App {
	Success("NewApp")
	app := &App{Router: NewRouter(), Server: &http.Server{}}
	return app
}

func (app *App) Run() {
	Success("App.Run")
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
	Success("Run")
	WebGo.Run()
}

// Router 路由部分
type Handlers struct {
	pattern        string
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
}

type Router struct {
	handlers []*Handlers
}

func NewRouter() *Router {
	Success("NewRouter")
	router := &Router{handlers: make([]*Handlers, 0)}
	return router
}

func (p *Router) Add(pattern string, c ControllerInterface) {
	Info("Router.Add", pattern, c)

	parts := strings.Split(pattern, "/")
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		// 搜索是否包含:id 类似声明
		if strings.HasPrefix(part, ":") {
			// 有表达式
			expr := "([^/]+)"
			// similar to expressjs: /user/:id([0-9]+)
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			} else if lidx := strings.LastIndex(part, ":"); lidx != 0 { // /user/:id:int or :id:string
				switch part[lidx:] {
				case ":int":
					expr = "([0-9]+)"
					part = part[:lidx]
				case ":string":
					expr = `([\w]+)`
					part = part[:lidx]
				}
			}

			params[j] = part[1:]
			parts[i] = expr
			j++
		}
	}
	//now create the Route
	t := reflect.Indirect(reflect.ValueOf(c)).Type()
	route := &Handlers{}
	route.controllerType = t
	if j != 0 {
		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {
			panic(regexErr)
			return
		}
		route.params = params
		route.regex = regex
	}
	route.pattern = pattern
	p.handlers = append(p.handlers, route)
}
func (p *Router) match(pattern string) *Handlers {

	return nil
}
func (p *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Server", "s")
	Info(r.Method, r.URL.Path)
	// static route must include staticPath keywork
	var runRouter *Handlers
	params := make(map[string]string)
	switch {
	case strings.HasPrefix(r.URL.Path, "/favicon"), strings.HasPrefix(r.URL.Path, "/robots"):
		file := filepath.Join(workPath, appConf.StaticPath, "images", r.URL.Path)
		http.ServeFile(rw, r, file)
		return
	case strings.HasPrefix(r.URL.Path, "/js"), strings.HasPrefix(r.URL.Path, "/css"), strings.HasPrefix(r.URL.Path, "/images"):
		file := filepath.Join(workPath, appConf.StaticPath, r.URL.Path)
		http.ServeFile(rw, r, file)
		return
	case strings.HasPrefix(r.URL.Path[1:], appConf.StaticPath):
		file := filepath.Join(workPath, r.URL.Path)
		http.ServeFile(rw, r, file)
		return
	default:
		pattern := r.URL.Path
		n := len(pattern)
		// 去除网址结尾斜杠
		if n > 1 && "/" == pattern[n-1:] {
			pattern = pattern[:n-1]
		}
		for _, route := range p.handlers {
			Info(pattern, ":", route.pattern)
			if n == 1 {
				if route.pattern == pattern {
					runRouter = route
					break
				} else {
					continue
				}
			}
			// 匹配 /user /user/ 同等效果
			if route.pattern == pattern {
				runRouter = route
				break
			}

			// Route Match
			if nil != route.regex {
				// 正则匹配路径
				if !route.regex.MatchString(pattern) {
					continue
				}
				matches := route.regex.FindStringSubmatch(pattern)
				if len(matches[0]) != len(pattern) {
					continue
				}
				if len(route.params) > 0 {
					values := r.URL.Query()
					for i, match := range matches[1:] {
						// TODO 添加传值
						values.Add(route.params[i], match)
						// r.Form.Add(route.params[i], match)
						params[route.params[i]] = match
					}
					//reassemble query params and add to RawQuery
					r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
				}
				runRouter = route
				break
			}
		}
	}

	if runRouter != nil {
		// 找到了注册的路由
		Success("Haha", runRouter.controllerType)
		vc := reflect.New(runRouter.controllerType)
		init := vc.MethodByName("Init")
		ct := &Context{ResponseWriter: rw, Request: r, Params: params}
		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(ct)
		in[1] = reflect.ValueOf(runRouter.controllerType.Name())
		init.Call(in)
		// Prepare
		in = make([]reflect.Value, 0)
		method := vc.MethodByName("Prepare")
		method.Call(in)
		// Rquest method
		rMethod := fmt.Sprintf("%s%s", strings.ToUpper(r.Method[:1]), strings.ToLower(r.Method[1:]))
		method = vc.MethodByName(rMethod)
		method.Call(in)

		method = vc.MethodByName("Destructor")
		method.Call(in)
	} else {
		Error("Other Method", r.URL.Path, r.Method)
		http.NotFound(rw, r)
	}

}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Params         map[string]string
}

type Controller struct {
	Ctx       *Context
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
	c.Ctx = ctx
	c.ChildName = cn
}

func (c *Controller) Write(content string) {
	c.Ctx.ResponseWriter.Write([]byte(content))
}

func (c *Controller) SetCookie(name string, value string, params ...interface{}) {
	var b bytes.Buffer
	D(value)
	fmt.Fprintf(&b, "%s=%s; Max-Age=2147483647; path=/", name, value)
	c.Ctx.ResponseWriter.Header().Set("Set-Cookie", b.String())
}

func (c *Controller) Get() {
	Warning(c.ChildName, "Get Method Not Create")
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Post() {
	Warning(c.ChildName, "Get Method Not Create")
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Delete() {
	Warning(c.ChildName, "Get Method Not Create")
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Head() {
	Warning(c.ChildName, "Get Method Not Create")
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Put() {
	Warning(c.ChildName, "Get Method Not Create")
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Patch() {
	Warning(c.ChildName, "Get Method Not Create")
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Options() {
	Warning(c.ChildName, "Get Method Not Create")
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Prepare()    {}
func (c *Controller) Destructor() {}

func Route(pattern string, c ControllerInterface) {
	Success("Route")
	WebGo.Router.Add(pattern, c)
}

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红

	info = "[INFO]"
	trac = "[TRAC]"
	erro = "[ERRO]"
	warn = "[WARN]"
	succ = "[SUCC]"
)

var goLogger = log.New(os.Stdout, "", log.Ltime)

func D(v ...interface{}) {
	funcName, _, line, _ := runtime.Caller(1)
	fName := runtime.FuncForPC(funcName).Name()
	goLogger.Printf("\x1b[%dm%s %v\x1b[0m %s:%d\n", color_yellow, trac, v, fName, line)
}

func Info(v ...interface{}) {
	goLogger.Printf("\x1b[%dm%s %v\x1b[0m\n", color_blue, info, v)
}

func Error(v ...interface{}) {
	funcName, _, line, _ := runtime.Caller(1)
	fName := runtime.FuncForPC(funcName).Name()
	goLogger.Printf("\x1b[%dm%s %v\x1b[0m %s:%d\n", color_red, erro, v, fName, line)
}

func Warning(v ...interface{}) {
	funcName, _, line, _ := runtime.Caller(1)
	fName := runtime.FuncForPC(funcName).Name()
	goLogger.Printf("\x1b[%dm%s\x1b[0m %v %s:%d\n", color_magenta, warn, v, fName, line)
}

func Success(v ...interface{}) {
	goLogger.Printf("\x1b[%dm%s\x1b[0m %v\n", color_green, succ, v)
}
