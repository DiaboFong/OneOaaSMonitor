package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error401() {
	c.Data["content"] = "page not found"
	c.TplName = "pc/401.html"
}

func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.TplName = "pc/404.html"
}

func (c *ErrorController) Error501() {
	c.Data["content"] = "server error"
	c.TplName = "pc/500.html"
}
