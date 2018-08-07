package pc

import (
	"github.com/astaxie/beego"
)

type ApplyController struct {
	beego.Controller
}

func (this *ApplyController) Get() {
	machinecode := this.Ctx.Input.Param(":machinecode")

	this.Data["Machinecode"] = machinecode
	this.TplName = "pc/apply.html"
}
