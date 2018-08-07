package pc

import (
	"strings"
	"github.com/h3llowor1d/checkmail"
	"regexp"
	"oneoaas.com/oneoaas_web/controllers"
	model "oneoaas.com/oneoaas_web/models"
	oneUtil "oneoaas.com/oneoaas_web/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"encoding/base64"
	"crypto/tls"
	"net/mail"
	"fmt"
	"net/smtp"
	"net"
	"time"
	"io"
)

type LicenseController struct {
	controllers.BaseController
}

func (this *LicenseController) GetLisenceView() {
	this.TplName = "pc/license.html"
}

type EmailConfig struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Host string  `json:"host"`
	Port int  `json:"port"`
}

func (this *LicenseController) GetLisenceKey() {
	email := this.GetString("email", "")
	vendor := this.GetString("vendor", "")
	machinecode := this.GetString("machinecode", "")

	logs.Info("代理商编号,"+vendor)

	// 校验email
	unsupportDomains := []string{"qq.com", "sina.com", "sina.cn", "163.com", "google.com", "foxmail.com", "gmail.com", "hotmail.com", "126.com"}
	for _, domain := range unsupportDomains {
		if strings.Contains(email, "@"+domain) {
			this.Data["json"] = this.JsonFormat(210, nil, "请勿使用个人邮箱")
			this.ServeJSON()
			return
		}
	}

	if err := checkmail.ValidateFormat(email); err != nil {
		this.Data["json"] = this.JsonFormat(210, nil, "邮箱格式校验不通过")
		this.ServeJSON()
		return
	}


	// 校验机器码
	reg, _ := regexp.Compile(`^[0-9A-Za-z]{32}$`)
	if !reg.MatchString(machinecode) {
		this.Data["json"] = this.JsonFormat(210, nil, "机器码校验不通过")
		this.ServeJSON()
		return
	}

	//验证代理商是否存在
	user := model.GetUserByVendorNum(vendor)
	if user == nil{
		this.Data["json"] = this.JsonFormat(210, nil, "代理商编号错误")
		this.ServeJSON()
		return
	}

	username_email := beego.AppConfig.String("username_email")
	password_email := beego.AppConfig.String("password_email")
	host_email := beego.AppConfig.String("host_email")
	port_email,_ := beego.AppConfig.Int("port_email")
	//config := `{"username":"","password":"2ZEgsJ6pQwoEqsmu","host":"smtp.exmail.qq.com","port":465}`

	if len(username_email)==0{
		this.Data["json"] = this.JsonFormat(210, nil, "发送邮箱没有配置")
		this.ServeJSON()
		return
	}

	if len(password_email)==0{
		this.Data["json"] = this.JsonFormat(210, nil, "发送邮箱没有密码")
		this.ServeJSON()
		return
	}

	if len(host_email)==0{
		this.Data["json"] = this.JsonFormat(210, nil, "邮件发送服务器没有配置")
		this.ServeJSON()
		return
	}

	//默认30天
	lisence,expireDate := oneUtil.CreateLisence(machinecode,30,0,0,"ee")

	from := mail.Address{"", username_email}
	to   := mail.Address{"", email}
	subj := "apply oneoaas monitor license success"
	//subj := "颁发OneOaaS-Monitor产品授权证书成功"
	// todo 格式化邮件模板
	body := "<div class='oneoaas'>" +
		"<style type='text/css'>" +
		".oneoaas{width:95%;margin: 20px auto}" +
		".oneoaas a{text-decoration: none;}" +
		".oneoaas .main{border-top: 1px solid #999;border-bottom: 1px solid #999;}" +
		".oneoaas p{font-size:13px;color:#333;}" +
		".oneoaas p.p_title{font-size:16px;color:#333;font-weight:bold;}" +
		".oneoaas h2.h2_text{color:#3580ff}" +
		".oneoaas h4{color:#999;width:90%;word-wrap:break-word;}" +
		".oneoaas span{color:#3580ff;}" +
		"</style>" +
		"<h2>" +
		"<a href='http://www.oneoaas.com/' target='_blank'>" +
		"<img src='http://www.oneoaas.com/static/img/pc/footer-logo.png'>" +
		"</a>" +
		"</h2>" +
		"<div class='main'>" +
		"<h2>尊敬的用户:</h2>" +
		"<h2>您好! 欢迎使用OneOaaS Monitor企业版。您所申请的授权码为:</h2>" +
		"<h4>"+lisence+"</h4>" +
		"<h2 class='h2_text'>授权码有效期为30天</h5>" +
		"</div>" +
		"<div>" +
		"<p class='p_title'>一、配置说明：</p>" +
		"<p>shell# vim /usr/share/oneoaas-monitor-ee/conf/app.conf</p>" +
		"<p>httpport = 4005</p>" +
		"<p>#database</p>" +
		"<p>#zabbix 数据库的类型，默认为MySQL</p>" +
		"<p>dbtype = 'mysql'</p>" +
		"<p>#zabbix 数据库的用户</p>" +
		"<p>dbuser = 'zabbix'</p>" +
		"<p>#zabbix 数据库的密码</p>" +
		"<p>dbpass = 'zabbix'</p>" +
		"<p>#zabbix 数据库的IP地址</p>" +
		"<p>dbhost = '127.0.0.1'</p>" +
		"<p>#zabbix 数据库的端口</p>" +
		"<p>dbport = 3306</p>" +
		"<p>#zabbix 数据库的库名称</p>" +
		"<p>dbname = 'zabbix'</p>" +
		"<p>#zabbix Web的地址，即Web API，Zabbix Web地址为http://10.10.10.38，,则配置为</p>" +
		"<p>ZbxUrl = 'http://10.10.10.38'</p>" +
		"<p># 带路径的配置，如果您的Zabbix WEB访问地址为http://10.10.10.38/zabbix,则配置为  </p>" +
		"<p>#ZbxUrl = 'http://10.10.10.38/zabbix'</p>" +
		"<p>#授权码</p>" +
		"<p>LicenseKey = '您申请到的授权码'</p>" +
		"<p>配置填写正确后，启动oneoaas-monitor-ee</p>" +
		"<p>shell# /etc/init.d/oneoaas-monitor-ee start</p>" +
		"<p class='p_title'>二、访问地址</p>" +
		"<p>http://您的IP地址:4005</p>" +
		"<p>用户: 您的Zabbix用户</p>" +
		"<p>密码: 您的Zabbix密码</p>" +
		"<p>需要注意的是ZbxUrl必须配置正确，如果您不确定，则可以用curl测试</p>" +
		"<p class='p_title'>三、手动测试</p>" +
		"<p>shell# curl  -X POST -H 'Content-Type: application/json' -d '{'jsonrpc': '2.0','method':'user.login','params':{'user':'Admin','password':'zabbix'},'id':0}'  http://10.10.10.38/api_jsonrpc.php</p>" +
		"<p>输出如下结果，则说明http://10.10.10.38为正确的Zabbix API地址</p>" +
		"<p>{'jsonrpc':'2.0','result':'b93e7c8387a5c55b64c29f6e7830622e','id':0}</p>" +
		"<p>如使用过程出现其他问题，可加客服QQ <span>3408827719</span></p>" +
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
		this.Data["json"] = this.JsonFormat(210, nil, "建立邮箱服务器TCP连接失败")
		this.ServeJSON()
		return
	}

	conn = tls.Client(conn,tlsconfig)
	if err != nil {
		beego.Error("建立邮箱服务器SSL连接失败"+err.Error())
		this.Data["json"] = this.JsonFormat(210, nil, "建立邮箱服务器SSL连接失败")
		this.ServeJSON()
		return
	}

	c,err := smtp.NewClient(conn, host)
	if err != nil {
		beego.Error("创建邮箱服务器客户端失败"+err.Error())
		this.Data["json"] = this.JsonFormat(210, nil, "创建邮箱服务器客户端失败")
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
				this.Data["json"] = this.JsonFormat(210, nil, "邮件发送认证失败")
				this.ServeJSON()
				return
			}
		}
	}


	// To && From
	if err = c.Mail(from.Address); err != nil {
		beego.Error("邮件发送地址错误"+err.Error())
		this.Data["json"] = this.JsonFormat(210, nil, "邮件发送地址错误")
		this.ServeJSON()
		return
	}

	if err = c.Rcpt(to.Address); err != nil {
		beego.Error("邮件接收地址错误"+err.Error())
		this.Data["json"] = this.JsonFormat(210, nil, "邮件接收地址错误")
		this.ServeJSON()
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		beego.Error("邮件发送内容解析错误"+err.Error())
		this.Data["json"] = this.JsonFormat(210, nil, "邮件发送内容解析错误")
		this.ServeJSON()
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		beego.Error("邮件发送失败"+err.Error())
		this.Data["json"] = this.JsonFormat(210, nil, "邮件发送失败")
		this.ServeJSON()
		return
	}

	err = w.Close()
	if err != nil {
		beego.Error("关闭邮件发送失败"+err.Error())
		this.Data["json"] = this.JsonFormat(210, nil, "关闭邮件发送失败")
		this.ServeJSON()
		return
	}

	//退出客户端
	c.Quit()

	licenseModel := model.License{
		ApplyDate:time.Now(),
		Duration:"30",
		ExpireDate:expireDate,
		User:user,
		LicenseKey:lisence,
	}

	_,err = licenseModel.AddLicense()
	if err != nil{
		beego.Error("保存License信息错误"+err.Error())
	}

	user.Licenses = append(user.Licenses, &licenseModel)
	if err = user.UpdateUser() ; err != nil{
		beego.Error("更新用户License信息错误"+err.Error())
	}


	this.Data["json"] = this.JsonFormat(200, nil, "邮件发送成功")
	this.ServeJSON()
}

type Result struct {
	RecordsTotal int `json:"recordsTotal"`
	Data []model.License `json:"data"`
}

func (this *LicenseController) ListLicenses() {

	who := this.Ctx.GetCookie("who")

	list,err := model.GetLicenses(who)
	if err != nil{
		beego.Error("获取License列表错误")
	}
	res := Result{
		RecordsTotal:len(list),
		Data:list,
	}
	this.Data["json"] = res
	this.ServeJSON()
	return
}

func TimeBuild(strTime string) time.Time {
	tm, _ := time.Parse("2006-01-02 15:04:05", strTime)
	return tm
}

func AddOneToEachElement(slice []*model.License) {
	for i := range slice {
		slice[i]=nil
	}
}

type smtpClient interface {
	Hello(string) error
	Extension(string) (bool, string)
	StartTLS(*tls.Config) error
	Auth(smtp.Auth) error
	Mail(string) error
	Rcpt(string) error
	Data() (io.WriteCloser, error)
	Quit() error
	Close() error
}

func smtpNewClient(conn net.Conn, host string) (smtpClient, error) {
	return smtp.NewClient(conn, host)
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}



func (c *LicenseController)Prepare(){

}

func (c *LicenseController) GetLicenseView() {
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

	c.TplName = "pc/license-manage.html"
}
