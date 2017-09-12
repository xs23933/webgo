package main

import (
	"github.com/xs23933/webgo"
)

type MainHandler struct {
	webgo.Controller
}

func (p *MainHandler) Get() {
	webgo.D(p.Ctx.Params["id"])
	p.SetCookie("a", p.Ctx.Params["id"], "domain.com", true)
	p.Write(`<html>
<head>
<link href="/css/app.css" rel="stylesheet" />
<script src='/js/app.js'></script>
</head>
<body>
<h1>Hello world</h1>
<img src='/images/logo.jpg' />
</body>
</html>
	`)
}
func main() {
	webgo.Route("/", &MainHandler{})
	webgo.Route("/user", &MainHandler{})
	webgo.Route("/user/:id", &MainHandler{})
	webgo.Route("/main/:id:int", &MainHandler{})
	webgo.Route("/info/:id(\\w+)", &MainHandler{})
	webgo.Route("/app/:id:string", &MainHandler{})
	webgo.Run()
}
