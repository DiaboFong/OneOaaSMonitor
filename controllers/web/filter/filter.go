package filter

import (
	"github.com/astaxie/beego/context"
	"strings"
)

var FilterIndex = func(ctx *context.Context) {
	keywords := []string{"Android", "iPhone", "iPod", "iPad", "Mobile", "Windows Phone", "MQQBrowser"}
	if strings.Contains(ctx.Request.RequestURI, "/m") {
		return
	}
	for i := 0; i < len(keywords); i++ {
		if strings.Contains(ctx.Request.UserAgent(), keywords[i]) {
			ctx.Redirect(302, "/m"+ctx.Request.RequestURI)
			return
		}
	}
}
