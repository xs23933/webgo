package webgo

import "testing"

type XHandler struct {
	Controller
}

var r = NewRouter()

func TestD(t *testing.T) {
	x := "/user/123/"
	D(x[len(x)-1:])
}

func TestRouteAdd(t *testing.T) {
	D("Test begin\n\n\n\n\n")
	r.Add("/:id:int", &XHandler{})
	r.Add("/user/:uid", &XHandler{})
	r.Add("/:ooxx([0-9]+)", &XHandler{})
	r.Add("/user/:id:int", &XHandler{})
	r.Add("/app/:cid:int", &XHandler{})
	r.Add("/goods/:ooxx([0-9]+)", &XHandler{})
	// r.Add("/", &XHandler{})
	route := r.match("/app/123")
	if route != nil {
		Success(route.controllerType)
	}
	route = r.match("/user/123/")
	if route != nil {
		Success(route.controllerType)
	}
	route = r.match("/")
	if route != nil {
		Success(route.controllerType)
	}
	route = r.match("/user/")
	if route != nil {
		Success(route.controllerType)
	}

}
