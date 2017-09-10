package webgo

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type Controllers struct {
	pattern        string
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
	methods        map[string]string
	hasMethod      bool
}

// Router class
type Router struct {
	routers    []*Controllers
	fixrouters []*Controllers
}

// NewRouter is Router constructor
func NewRouter() *Router {
	router := &Router{}
	return router
}

func (p *Router) Add(pattern string, c ControllerInterface) {
	parts := strings.Split(pattern, "/")
	params := make(map[int]string)
	j := 0

	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "(.+)"
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			} else if lindex := strings.LastIndex(part, ":"); lindex != 0 {
				switch part[lindex:] {
				case ":int":
					expr = "([0-9]+)"
					part = part[:lindex]
				case ":string":
					expr = `([\w]+)`
					part = part[:lindex]
				}
			}
			params[j] = part
			parts[i] = expr
			j++
		}
		if strings.HasPrefix(part, "*") {
			expr := "(.+)"
			if part == "*.*" {
				params[j] = ":path"
				parts[i] = "([^.]+).([^.]+)"
				j++
				params[j] = ":ext"
				j++
			} else {
				params[j] = ":splat"
				parts[i] = expr
				j++
			}
		}
	}
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	if j == 0 {
		route := &Controllers{}
		route.pattern = pattern
		route.controllerType = t
		p.fixrouters = append(p.fixrouters, route)
	} else {
		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {
			panic(regexErr)
			return
		}

		route := &Controllers{}
		route.regex = regex
		route.params = params
		route.pattern = pattern

		route.controllerType = t
		p.routers = append(p.routers, route)
	}
}

// ServeHTTP define from http.Server.Handler
func (p *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w := &Response{writer: rw}

	w.Header().Set("Server", "s")
	var runRouter *Controllers
	// var findRouter bool
	params := make(map[string]string)

	//static file server
	for prefix, staticDir := range StaticDir {
		if strings.HasPrefix(r.URL.Path[1:], prefix) {
			file := staticDir + r.URL.Path //[len(prefix)+1:]
			fmt.Println(file, prefix)
			http.ServeFile(w, r, file)
			w.started = true
			return
		}
	}

	requestPath := r.URL.Path

	// var requestBody []byte

	for _, route := range p.fixrouters {
		n := len(requestPath)
		if n == 1 {
			if requestPath == route.pattern {
				runRouter = route
				// findRouter = true
				break
			} else {
				continue
			}
		}
		if (requestPath[n-1] != '/' && route.pattern == requestPath) ||
			(requestPath[n-1] == '/' && len(route.pattern) >= n-1 && requestPath[0:n-1] == route.pattern) {
			runRouter = route
			// findRouter = true
			break
		}
	}

	if runRouter != nil {
		vc := reflect.New(runRouter.controllerType)

		init := vc.MethodByName("Init")
		in := make([]reflect.Value, 2)
		ct := &Context{ResponseWriter: w, Request: r, Params: params}

		in[0] = reflect.ValueOf(ct)
		in[1] = reflect.ValueOf(runRouter.controllerType.Name())
		init.Call(in)

		in = make([]reflect.Value, 0)
		method := vc.MethodByName("Prepare")
		method.Call(in)
		// call method Change to First Upper case
		requestMethod := strings.ToUpper(r.Method[:1]) + strings.ToLower(r.Method[1:])
		method = vc.MethodByName(requestMethod)
		fmt.Println(requestMethod, r.URL.Path)
		method.Call(in)
		method = vc.MethodByName("Destructor")
		method.Call(in)
	}

	if w.started == false {
		http.NotFound(w, r)
	}
}

type Response struct {
	writer  http.ResponseWriter
	started bool
	status  int
}

func (w *Response) Header() http.Header {
	return w.writer.Header()
}

func (w *Response) Write(p []byte) (int, error) {
	w.started = true
	return w.writer.Write(p)
}

func (w *Response) WriteHeader(code int) {
	w.status = code
	w.started = true
	w.writer.WriteHeader(code)
}
