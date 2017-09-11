package main

import (
	"github.com/xs23933/webgo"
)

type MainHandler struct {
	webgo.Controller
}

func (p *MainHandler) Get() {
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
	webgo.Route("/user/", &MainHandler{})
	webgo.Run()
}
