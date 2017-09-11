# app
go web framework


#### Install

```go
go get github.com/xs23933/webgo
```
#### Example
```go
package main

import (
	"fmt"
	"github.com/xs23933/webgo"
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
