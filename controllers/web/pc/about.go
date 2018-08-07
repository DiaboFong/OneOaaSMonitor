package pc

import (
	"github.com/astaxie/beego"
)

type AboutController struct {
	beego.Controller
}

func (c *AboutController) Get() {
	c.TplName = "pc/about.html"
}
