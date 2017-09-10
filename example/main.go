package main

import "github.com/xs23933/webgo"
import "fmt"

type MainHandler struct {
	webgo.Controller
}

func (p *MainHandler) Get() {
	fmt.Println("call MainHandler")
	p.Write("Fuck Men")
}

type UserHandler struct {
	webgo.Controller
}

func (p *UserHandler) Get() {
	fmt.Println("call UserHandler")
	p.Write("This is User Page")
}

func main() {
	webgo.Route("/", &MainHandler{})
	webgo.Route("/user", &UserHandler{})
	// webgo.Route("/user/(\d+)", &UserHandler{})
	webgo.Run()
}
