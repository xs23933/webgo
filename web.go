package webgo

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var (
	WebGo         *App
	appPath       string
	appConfigPath string
	appConfig     AppConfig
)

func init() {
	appPath, _ = os.Getwd()
	WebGo = InitApp()
	appConfigPath = filepath.Join(appPath, "conf", "app.conf")
}

type App struct {
	Handler *Handlers
	Server  *http.Server
}

// Config Define
type AppConfig struct {
	Name     string
	Server   server
	Database database
}

type database struct {
	Server   string
	Port     int
	User     string
	Password string
	Type     string
}
type server struct {
	IP string
}

func InitApp() *App {
	// get Project Run Path
	appConfigPath := filepath.Join(appPath, "conf", "app.conf")
	fmt.Println(appConfigPath)
	if _, err := toml.DecodeFile(appConfigPath, &appConfig); err != nil {
		panic(err)
	}
	app := &App{Handler: NewHandlers(), Server: &http.Server{}}
	return app
}

func (app *App) Run() {
	var (
		endRunning = make(chan bool, 1)
	)
	fmt.Println(app.Handler)
	app.Server.Handler = app.Handler
	app.Server.Addr = appConfig.Server.IP

	go func() {
		if err := app.Server.ListenAndServe(); err != nil {
			endRunning <- true
		}
	}()
	<-endRunning
}

func Run() {
	WebGo.Run()
}

type Handler struct {
	pattern string
	h       http.Handler
}

type Handlers struct {
	handlers map[string]*Handler
}

func NewHandlers() *Handlers {
	handles := &Handlers{handlers: make(map[string]*Handler)}
	CtrlPath := filepath.Join(appPath, "controllers")
	fmt.Println(CtrlPath)
	return handles
}

func (p *Handlers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("what happend! u call ServeHTTP")
}
