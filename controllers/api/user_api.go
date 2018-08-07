package api

import (
	"oneoaas.com/oneoaas_web/controllers"
	"oneoaas.com/oneoaas_web/models"
	"github.com/astaxie/beego"
	"strings"
	//"os/user"
	"strconv"
	"oneoaas.com/oneoaas_web/util"
	"fmt"
	"time"
	"net/mail"
	"encoding/base64"
	"crypto/tls"
	"net"
	"net/smtp"
)

type UserController struct {
	controllers.BaseController
}

//添加用户
func (this *UserController) AddUser() {
	var user models.User
	username := this.GetString("username", "")
	email := this.GetString("email", "")
	phone := this.GetString("phone", "")
	company := this.GetString("company", "")
	work := this.GetString("work", "")
	machinecode := this.GetString("machinecode", "")
	if username == "" || email == "" || phone == "" || company == "" || work == "" || machinecode == "" {
		this.Data["json"] = this.JsonFormat(204, nil, "信息不完整")
		this.ServeJSON()
		return
	}

	//CacheObj.Exists()
	v := models.GetSmsCodeByPhone(phone)
	//判断手机号是否通过短信验证码
	if v.Phone != phone || v.IsTrue == 0 {
		this.Data["json"] = this.JsonFormat(204, nil, "手机号码未通过短信验证")
	} else {
		user.Username = username
		user.Email = email
		user.Phone = phone
		user.Company = company
		user.Work = work
		user.Machinecode = machinecode
		_, err := models.AddUser(user)
		if err != nil {
			this.Data["json"] = this.JsonFormat(204, nil, "用户申请发送失败,请联系管理员")
		} else {
			this.Data["json"] = this.JsonFormat(200, nil, "用户申请发送成功,请等待审核")
		}
	}
	this.ServeJSON()
}

//func (this *UserController) ListUsers(){
//	this.Data["json"] = models.GetUsers()
//	this.ServeJSON()
//}

type Result struct {
	RecordsTotal int `json:"recordsTotal"`
	Data []models.User `json:"data"`
}

func (this *UserController) ListUsers() {
	if beego.GlobalSessions.GetActiveSession() <=0{
		beego.Info("激活的session 为空")
		this.Redirect("/login",302)
		return
	}


	if this.CruSession == nil{
		//如果是
		sessionId := this.Ctx.GetCookie("oneoaasessionId")
		beego.Info("sessionId is "+sessionId)
		if len(sessionId) == 0{
			beego.Error("cookie 中没有 sessionid")
			this.Redirect("/login",302)
			return
		}
		_,err:= beego.GlobalSessions.GetSessionStore(sessionId)
		if err != nil{
			//数据库没有找到当前的sessionid
			beego.Error("数据库没有找到当前的sessionid")
			this.Redirect("/login",302)
			return
		}
	}else {
		oneoaasessionId := this.CruSession.Get("oneoaasessionId")
		if oneoaasessionId == nil{
			beego.Error("当前没有oneoaasessionId 无法登录")
			this.Redirect("/login",302)
			return
		}
	}

	list,err := models.GetUsers()
	if err != nil{
		beego.Error("获取用户列表错误")
	}
	res := Result{
		RecordsTotal:len(list),
		Data:list,
	}
	this.Data["json"] = res
	this.ServeJSON()
	return
}

func (this *UserController) GrantUser() {
	// 获取 id ,并去掉尾部 ，
	str := strings.TrimSuffix(this.GetString("id"),",")
	if str =="" {
		this.Data["json"] = this.JsonFormat(204, nil, "数据为空")
		this.ServeJSON()
		return
	}
	var ids []string
	// 将字符串分割成数组
	ids = strings.Split(str,",")
	for _,value := range ids{
		// 将数组内 id 转换为 int64
		id,error := strconv.ParseInt(value,10,64)
		if error != nil{
			this.Data["json"] = this.JsonFormat(204, nil, "字符串转换成数字失败")
			this.ServeJSON()
			break
		}
		user := models.GetUserById(id)
		if user == nil {
			this.Data["json"] = this.JsonFormat(204, nil, "用户不存在")
			this.ServeJSON()
			break
		}

		salt := util.GenerateSalt()
		password2 := salt
		//加密密码
		user.Password = util.StrtoMd5(password2)
		//明文密码
		user.Password2 = password2
		result := user.UpdateUserVendorNum()
		if result {
			//todo 自动生成密码，需要发送邮件
			this.Data["json"] = this.JsonFormat(200, nil, "授权该用户为代理商成功")

			// 发送邮件
			username_email := beego.AppConfig.String("username_email")
			password_email := beego.AppConfig.String("password_email")
			host_email := beego.AppConfig.String("host_email")
			port_email,_ := beego.AppConfig.Int("port_email")
			//config := `{"username":"","password":"2ZEgsJ6pQwoEqsmu","host":"smtp.exmail.qq.com","port":465}`

			if len(username_email)==0{
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,发送邮箱没有配置")
				this.ServeJSON()
				return
			}

			if len(password_email)==0{
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,发送邮箱没有密码")
				this.ServeJSON()
				return
			}

			if len(host_email)==0{
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,邮件发送服务器没有配置")
				this.ServeJSON()
				return
			}

			from := mail.Address{"", username_email}
			to   := mail.Address{"", user.Email}
			subj := "apply oneoaas monitor vendor success"
			//subj := "代理商用户名和密码"
			// todo 格式化邮件模板
			body := "<div class='oneoaas'>" +
				"<style type='text/css'>" +
				".oneoaas{width:95%;margin: 20px auto}" +
				".oneoaas a{text-decoration: none;}" +
				".oneoaas .main{border-top: 1px solid #999;border-bottom: 1px solid #999;}" +
				".oneoaas h2.h2_text{color:#3580ff}" +
				"</style>" +
				"<h2>" +
				"<a href='http://www.oneoaas.com/' target='_blank'>" +
				"<img src='http://www.oneoaas.com/static/img/pc/footer-logo.png'>" +
				"</a>" +
				"</h2>" +
				"<div class='main'>" +
				"<h2>尊敬的用户:</h2>" +
				"<h2>您好! 恭喜成为OneOaaS Monitor代理商。</h2>" +
				"<h2>用户名为：" + user.Username + "</h2>" +
				"<h2>密&nbsp;码&nbsp;&nbsp;为：" + user.Password2 + "</h2>" +
				"<h2>代理商编号为：" + user. VendorNum+ "</h2>" +
				"<h2>请前往：<a href='http://www.oneoaas.com/login'>http://www.oneoaas.com/login</a></h2>" +
				"</div>" +
				"</div>"
			// Setup headers
			headers := make(map[string]string)
			headers["From"] = from.String()
			headers["To"] = to.String()
			headers["Subject"] = subj
			headers["Content-Type"] = "text/html; charset=\"utf-8\""
			headers["Content-Transfer-Encoding"] = "base64"

			// Setup message
			message := ""
			for k, v := range headers {
				message += fmt.Sprintf("%s: %s\r\n", k, v)
			}
			message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))


			// Connect to the SMTP Server
			servername := fmt.Sprintf("%s:%d",host_email,port_email)

			host, _, _ := net.SplitHostPort(servername)

			// TLS config 使用ssl认证发送
			tlsconfig := &tls.Config {
				InsecureSkipVerify: true,
				ServerName: host,
			}

			conn, err := net.DialTimeout("tcp", servername, 10*time.Second)
			if err != nil {
				beego.Error("建立邮箱服务器TCP连接失败"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,建立邮箱服务器TCP连接失败")
				this.ServeJSON()
				return
			}

			conn = tls.Client(conn,tlsconfig)
			if err != nil {
				beego.Error("建立邮箱服务器SSL连接失败"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,建立邮箱服务器SSL连接失败")
				this.ServeJSON()
				return
			}

			c,err := smtp.NewClient(conn, host)
			if err != nil {
				beego.Error("创建邮箱服务器客户端失败"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,创建邮箱服务器客户端失败")
				this.ServeJSON()
				return
			}



			if ok, auths := c.Extension("AUTH"); ok {
				if strings.Contains(auths, "CRAM-MD5") {
					//
				} else if strings.Contains(auths, "LOGIN") && !strings.Contains(auths, "PLAIN") {
					//
				} else {
					auth := smtp.PlainAuth("",username_email, password_email, host_email)
					if err = c.Auth(auth); err != nil {
						c.Close()
						beego.Error("邮件发送认证失败"+err.Error())
						this.Data["json"] = this.JsonFormat(210, nil, "授权成功,邮件发送认证失败")
						this.ServeJSON()
						return
					}
				}
			}


			// To && From
			if err = c.Mail(from.Address); err != nil {
				beego.Error("邮件发送地址错误"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,邮件发送地址错误")
				this.ServeJSON()
				return
			}

			if err = c.Rcpt(to.Address); err != nil {
				beego.Error("邮件接收地址错误"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,邮件接收地址错误")
				this.ServeJSON()
				return
			}

			// Data
			w, err := c.Data()
			if err != nil {
				beego.Error("邮件发送内容解析错误"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,邮件发送内容解析错误")
				this.ServeJSON()
				return
			}

			_, err = w.Write([]byte(message))
			if err != nil {
				beego.Error("邮件发送失败"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,邮件发送失败")
				this.ServeJSON()
				return
			}

			err = w.Close()
			if err != nil {
				beego.Error("关闭邮件发送失败"+err.Error())
				this.Data["json"] = this.JsonFormat(210, nil, "授权成功,关闭邮件发送失败")
				this.ServeJSON()
				return
			}

			//退出客户端
			c.Quit()
			this.Data["json"] = this.JsonFormat(200, nil, "授权成功,邮件发送成功")
			this.ServeJSON()
		} else {
			this.Data["json"] = this.JsonFormat(204, nil, "授权该用户为代理商失败")
			this.ServeJSON()
			break
		}
	}

	return
}

func (c *UserController)Prepare(){

}

func (c *UserController) GetListUserView() {
	if beego.GlobalSessions.GetActiveSession() <=0{
		beego.Info("激活的session 为空")
		c.Redirect("/login",302)
		return
	}


	if c.CruSession == nil{
		//如果是
		sessionId := c.Ctx.GetCookie("oneoaasessionId")
		beego.Info("sessionId is "+sessionId)
		if len(sessionId) == 0{
			beego.Error("cookie 中没有 sessionid")
			c.Redirect("/login",302)
			return
		}
		_,err:= beego.GlobalSessions.GetSessionStore(sessionId)
		if err != nil{
			//数据库没有找到当前的sessionid
			beego.Error("数据库没有找到当前的sessionid")
			c.Redirect("/login",302)
			return
		}
	}else {
		oneoaasessionId := c.CruSession.Get("oneoaasessionId")
		if oneoaasessionId == nil{
			beego.Error("当前没有oneoaasessionId 无法登录")
			c.Redirect("/login",302)
			return
		}
	}

	who := c.Ctx.GetCookie("who")
	if who == "oneoaas"{
		c.Data["IsAdmin"]=true
	}else{
		c.Data["IsAdmin"]=false
	}
	c.TplName = "pc/user.html"
}
