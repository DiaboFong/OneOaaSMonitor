package pc

import (
	"github.com/astaxie/beego"
)

type DownloadController struct {
	beego.Controller
}

func (c *DownloadController) Get() {
	c.TplName = "pc/download.html"
}
