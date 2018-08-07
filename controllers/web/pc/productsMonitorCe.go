package pc

import (
	"github.com/astaxie/beego"
)

type ProductsMonitorCeController struct {
	beego.Controller
}

func (c *ProductsMonitorCeController) Get() {
	c.TplName = "pc/products-monitor-ce.html"
}
