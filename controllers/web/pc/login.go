package pc

import (
	"github.com/astaxie/beego"
	"strings"
	"oneoaas.com/oneoaas_web/util"
	"oneoaas.com/oneoaas_web/models"
	"oneoaas.com/oneoaas_web/controllers"
)

type LoginController struct {
	controllers.BaseController
}

func (c *LoginController) GetView() {
	c.TplName = "pc/login.html"
}

func (c *LoginController) Login() {

	sessionStore,err := beego.GlobalSessions.SessionStart(c.Ctx.ResponseWriter,c.Ctx.Request)
	defer sessionStore.SessionRelease(c.Ctx.ResponseWriter)
	if err != nil{
		beego.Error("登录session开启错误")
		c.Data["json"] =  map[string]interface{}{
			"code": 400,
			"msg": "登录失败",
		}
		c.ServeJSON()
		return
	}

	sess := c.StartSession()
	sessionId := ""
	if sess == nil {
		c.CruSession = sessionStore
		sessionId = sessionStore.SessionID()
		beego.Info("设置session id "+sessionId)
		c.SetSession("oneoaasessionId",sessionId)
	}else{
		//session 已经存在
		beego.Info("oneoaasessionId 已经存在")
		c.CruSession.SessionRelease(c.Ctx.ResponseWriter)
		sessionId = c.CruSession.Get("oneoaasessionId").(string)
	}

	username := c.GetString("username")
	if len(strings.TrimSpace(username)) ==0 {
		c.Data["json"] =  map[string]interface{}{
			"code": 400,
			"msg": "用户名或密码不能为空",
		}
		c.ServeJSON()
		return
	}
	password := c.GetString("password")
	if len(strings.TrimSpace(password)) ==0{
		c.Data["json"] =  map[string]interface{}{
			"code": 400,
			"msg": "用户名或密码不能为空",
		}
		c.ServeJSON()
		return
	}
	encodePassword := util.StrtoMd5(password)
	user := models.GetUserByUsernameAndPassword(username,encodePassword)
	if user == nil {
		c.Data["json"] =  map[string]interface{}{
			"code": 400,
			"msg": "登录失败",
		}
		c.ServeJSON()
		//c.Redirect("/login",200)
		return
	}else {
		c.Ctx.SetCookie("who",username)
		c.Data["json"] =  map[string]interface{}{
			"code": 200,
			"msg": "登录成功",
		}
		//c.Redirect("/users",200)
	}
	c.ServeJSON()
}

//退出函数
func (c *LoginController) LogOut() {
	c.Ctx.SetCookie("oneoaasessionId", "")
	c.Ctx.SetCookie("who", "")
	c.Redirect("/login",302)
}
