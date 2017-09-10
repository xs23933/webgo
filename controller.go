package webgo

import (
	"fmt"
)

type Controller struct {
	ctx       *Context
	ChildName string
}

type ControllerInterface interface {
	Init(ct *Context, cn string)
	Get()
	Post()
	Head()
	Delete()
}

func (c *Controller) Init(ctx *Context, cn string) {
	c.ctx = ctx
	c.ChildName = cn
	fmt.Println("what call Init")
}

func (c *Controller) Write(content string) {
	c.ctx.ResponseWriter.Write([]byte(content))
}
func (c *Controller) Prepare()    {}
func (c *Controller) Destructor() {}
func (c *Controller) Get() {
	c.Write("This is Default Method. Please Create this Method to You Controller.")
}
func (c *Controller) Post() {
	c.Write("This is Default Method. Please Create this Method to You Controller.")
}
func (c *Controller) Head() {
	c.Write("This is Default Method. Please Create this Method to You Controller.")
}
func (c *Controller) Delete() {
	c.Write("This is Default Method. Please Create this Method to You Controller.")
}
func (c *Controller) Put() {
	c.Write("This is Default Method. Please Create this Method to You Controller.")
}
func (c *Controller) Patch() {
	c.Write("This is Default Method. Please Create this Method to You Controller.")
}
func (c *Controller) Options() {
	c.Write("This is Default Method. Please Create this Method to You Controller.")
}
