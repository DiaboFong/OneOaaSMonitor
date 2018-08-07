package pc

import (
	"github.com/astaxie/beego"
)

type ProductsCmdbController struct {
	beego.Controller
}

func (c *ProductsCmdbController) Get() {
	c.TplName = "pc/products-cmdb.html"
}
