package routers

import (
	"oneoaas.com/oneoaas_web/controllers/api"
	"oneoaas.com/oneoaas_web/controllers/web"
	//"oneoaas.com/oneoaas_web/controllers/web/mobile"
	"oneoaas.com/oneoaas_web/controllers/web/pc"

	"github.com/astaxie/beego"
)

func init() {
	//PC端路由
	//beego.Router("/", &pc.IndexController{})
	beego.Router("/products", &pc.ProductsController{})
	beego.Router("/download", &pc.DownloadController{})
	beego.Router("/products/cmdb", &pc.ProductsCmdbController{})
	beego.Router("/products/monitor-ee", &pc.ProductsMonitorEeController{})
	beego.Router("/products/monitor-ce", &pc.ProductsMonitorCeController{})
	beego.Router("/community", &pc.CommunityController{})
	beego.Router("/about", &pc.AboutController{})
	//受到到保护的代理商相关API
	beego.Router("/vendor", &pc.LicenseController{},"get:GetLisenceView;post:GetLisenceKey")

	//受到到保护的代理商相关API
	// 申请
	beego.Router("/apply/:machinecode(^[0-9A-Za-z]{32}$)", &pc.ApplyController{}, "get:Get")

	//移动端路由
/*	beego.Router("/m", &mobile.IndexController{})
	beego.Router("/m/products", &mobile.ProductsController{})
	beego.Router("/m/products/cmdb", &mobile.ProductsCmdbController{})
	beego.Router("/m/products/monitor-ee", &mobile.ProductsMonitorEeController{})
	beego.Router("/m/products/monitor-ce", &mobile.ProductsMonitorCeController{})
	beego.Router("/m/community", &mobile.CommunityController{})
	beego.Router("/m/about", &mobile.AboutController{})*/

	//用户信息
	/*
		username String
		email String
		phone String
		company String
		work String
		machinecode String
	*/
	beego.Router("/api/monitor/user", &api.UserController{}, "post:AddUser")

	//-------受到保护的API
	//beego.Router("/users", &api.UserController{}, "get:ListUsers")
	beego.Router("/users/grant", &api.UserController{}, "get:GrantUser")

	beego.Router("/users/list/license", &pc.LicenseController{}, "get:ListLicenses")
	// 用户管理界面
	beego.Router("/users/user", &api.UserController{}, "get:GetListUserView")
	// 用户数据
	beego.Router("/users/OneS/data", &api.UserController{}, "get:ListUsers")

	//-------受到保护的API

	//图片验证码
	//beego.Router("/captcha", &notice.CaptchaController{}, "get:GetCaptcha")
	//beego.Router("/api/verify", &notice.NoticeSmsController{}, "post:VerifyCaptcha")
	//beego.Handler("/api/captcha/*.png", captcha.Server(240, 80))
	// 发送短信验证码
	beego.Router("/api/send_sms_code", &api.NoticeSmsController{}, "post:SendSmsCode")

	// 申请试用
	beego.Router("/api/apply", &api.ApplyController{}, "post:Apply")
	beego.Router("/api/apply/license", &api.ApplyController{}, "get:GetLisenceKey")

	// 下载
	beego.Router("/download/v2.0/:filename", &web.DownLoadController{}, "get:Get")

	// login 登录API验证
	beego.Router("/api/login", &pc.LoginController{}, "post:Login")

	// login 进入登录页面
	beego.Router("/", &pc.LoginController{},"get:GetView")

	// logout 退出系统
	beego.Router("/logout", &pc.LoginController{},"get:LogOut")


	// license 管理界面
	beego.Router("/license", &pc.LicenseController{},"get:GetLicenseView")




}
