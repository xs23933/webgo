package webgo

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var (
	WebGo         *Application
	appConfigPath string
	appPath       string
	ctrlPath      string
	appConfig     AppConfig
	StaticDir     map[string]string
)

// AppConfig config define
// host e.g :8080 or 0.0.0.0::8080
type AppConfig struct {
	Host     string
	Database database
}

type database struct {
	Server   string
	Port     int
	User     string
	Password string
	Type     string
}

// init package
func init() {
	appPath, _ = os.Getwd()
	StaticDir = make(map[string]string)
	appConfigPath = filepath.Join(appPath, "conf", "app.conf")
	ctrlPath = filepath.Join(appPath, "controller")
	StaticDir["favicon"] = "public/images"
	StaticDir["images"] = "public"
	StaticDir["js"] = "public"
	StaticDir["css"] = "public"
	WebGo = NewApplication()
}

// Application is main point
type Application struct {
	Router *Router
	Server *http.Server
}

// NewApplication is Application constructor
func NewApplication() *Application {
	if _, err := toml.DecodeFile(appConfigPath, &appConfig); err != nil {
		panic(err)
	}
	app := &Application{Router: NewRouter(), Server: &http.Server{}}
	return app
}

// Run call start application
// WegGo.Run() start application
func (app *Application) Run() {
	endRunning := make(chan bool, 1)
	app.Server.Addr = appConfig.Host
	app.Server.Handler = app.Router

	go func() {
		if err := app.Server.ListenAndServe(); err != nil {
			endRunning <- true
		}
	}()
	<-endRunning
}

func (app *Application) Route(path string, c ControllerInterface) *Application {
	app.Router.Add(path, c)
	return app
}

func Route(path string, c ControllerInterface) *Application {
	WebGo.Route(path, c)
	return WebGo
}

// Run Start Application
func Run() {
	// fmt.Println(ctrlPath, appPath, appConfigPath, appConfig)
	WebGo.Run()
}
