package pc

import (
	"github.com/astaxie/beego"
)

type ProductsMonitorEeController struct {
	beego.Controller
}

func (c *ProductsMonitorEeController) Get() {
	c.TplName = "pc/products-monitor-ee.html"
}
