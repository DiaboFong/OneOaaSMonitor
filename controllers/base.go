package controllers

import (
	"github.com/astaxie/beego"
	"net"
	"strings"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) JsonFormat(code int, content interface{}, msg string) (json map[string]interface{}) {
	json = map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	return json
}

// 获取用户IP地址
func (ctx *BaseController) GetClientIp() string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(ctx.Ctx.Request.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	//[::1]:61721 localhost有IPV6会被解析，正常情况为127.0.0.1:61856
	s := strings.Split(ctx.Ctx.Request.RemoteAddr, ":")
	if len(s) == 4 {
		return "[::1]"
	}
	return s[0]
}
