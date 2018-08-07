package pc

import (
	"github.com/astaxie/beego"
)

type ProductsController struct {
	beego.Controller
}

func (c *ProductsController) Get() {
	c.TplName = "pc/products.html"
}
