package web

import (
	"oneoaas.com/oneoaas_web/controllers"
	"oneoaas.com/oneoaas_web/util"
)

type DownLoadController struct {
	controllers.BaseController
}

func (this *DownLoadController) Get() {
	filename := this.Ctx.Input.Param(":filename")
	if filename == "" {
		this.Abort("404")
	}
	filepath := "download/" + filename
	exists, err := util.PathExists(filepath)
	if err != nil {
		this.Abort("500")
	}

	if exists {
		this.Ctx.Output.Download(filepath)
	} else {
		this.Abort("404")
	}
}
