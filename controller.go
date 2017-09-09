package webgo

import (
	"net/http"
)

var (
	Request  http.Request
	Response *http.ResponseWriter
)

type Controller struct {
	res  *http.ResponseWriter
	req  *http.Request
	path string
}

type ControllerInterface interface {
	Get()
	Post()
	Delete()
	Put()
	Head()
	Render() error
}

func (c *Controller) Init(req *http.Request, res *http.ResponseWriter) {
	c.req = req
	c.res = res
	c.path = "/"
}
