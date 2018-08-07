package pc

import (
	"github.com/astaxie/beego"
)

type CommunityController struct {
	beego.Controller
}

func (c *CommunityController) Get() {
	c.TplName = "pc/community.html"
}
