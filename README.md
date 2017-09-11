# app
go web framework

#### Example
```go
package main

import (
	"fmt"
	"webgo"
)

type MainHandler struct {
	webgo.Controller
}

func (p *MainHandler) Get() {
	p.Write("Hello world!")
}
func main() {
	webgo.Route("/", &MainHandler{})
	webgo.Run()
}
```
