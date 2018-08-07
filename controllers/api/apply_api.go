package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	model "oneoaas.com/oneoaas_web/models"
	oneUtil "oneoaas.com/oneoaas_web/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/h3llowor1d/checkmail"
	"net/http"
	"oneoaas.com/oneoaas_web/controllers"
	"oneoaas.com/oneoaas_web/models"
	"oneoaas.com/oneoaas_web/util"
	"regexp"
	"strings"
	"unicode/utf8"
	"time"
	"encoding/base64"
	"net"
	"net/mail"
	"crypto/tls"
	"net/smtp"
	"net/url"
)

type ApplyController struct {
	controllers.BaseController
}

func (this *ApplyController) Prepare() {
	this.EnableRender = false
}

func (this *ApplyController) GetLisenceKey() {
	machinecode := this.GetString("machinecode", "")
	beego.Info("申请创业板，钉钉自动发送的机器码是:"+machinecode)
	// 校验机器码
	reg, _ := regexp.Compile(`^[0-9A-Za-z]{32}$`)
	if !reg.MatchString(machinecode) {
		this.Data["json"] = this.JsonFormat(210, nil, "机器码校验不通过")
		this.ServeJSON()
		return
	}

	user := model.GetUserByMachinecode(machinecode)
	if user == nil{
		this.Data["json"] = this.JsonFormat(210, nil, "机器码不存在，无法申请")
		this.ServeJSON()
		return
	}
	beego.Info("申请创业板，钉钉自动发送 user Machinecode:"+user.Machinecode)
	username_email := beego.AppConfig.String("username_email")
	password_email := beego.AppConfig.String("password_email")
	host_email := beego.AppConfig.String("host_email")
	port_email,_ := beego.AppConfig.Int("port_email")

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
	lisence,expireDate := oneUtil.CreateLisence(user.Machinecode,30,0,0,"ee")

	from := mail.Address{"", username_email}
	to   := mail.Address{"", user.Email}
	subj := "apply oneoaas monitor-ee license success"
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
		"<img src='http://www.oneoaas.com/static/img/news/logo-blue.png'>" +
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
		"<p>shell# vim /usr/share/oneoaas-monitor-se/conf/app.conf</p>" +
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
		"<p>shell# curl -X POST -H 'Content-Type: application/json' -d '{\"jsonrpc\": \"2.0\",\"method\":\"user.login\",\"params\":{\"user\":\"Admin\",\"password\":\"zabbix\"},\"id\":0}' http://10.10.10.38/api_jsonrpc.php</p>" +
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

	//如果Licenses为0,则发送
	if len(user.Licenses) == 0{
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
		this.Data["json"] = this.JsonFormat(200, nil, "邮件发送成功")
		this.ServeJSON()
	}else{
		//如果存在则不发送
		beego.Info("邮箱:"+user.Email+",电话:"+user.Phone+",机器码:"+machinecode+",对应的授权码已经存在")
		//msg := fmt.Sprintf(`
		//	对应的授权码已经存在
		//	手机号:%s
		//	邮箱:%s
		//	公司:%s
		//	岗位:%s
		//	机器码:%s
		//	授权码:%s
		//`, user.Phone, user.Email, user.Company, user.Work, machinecode,user.Licenses[0].LicenseKey)
		//sendDingdingMessage(msg)
		this.Data["json"] = this.JsonFormat(200, nil, "授权码已经存在")
		this.ServeJSON()
	}

}


// 申请试用
func (this *ApplyController) Apply() {
	username := this.GetString("username", "")
	email := this.GetString("email", "")
	phone := this.GetString("phone", "")
	company := this.GetString("company", "")
	work := this.GetString("work", "")
	smscode := this.GetString("smscode", "")
	machinecode := this.GetString("machinecode", "")
	if username == "" || email == "" || phone == "" || company == "" || work == "" || machinecode == "" || smscode == "" {
		this.Data["json"] = this.JsonFormat(210, nil, "参数不完整")
		this.ServeJSON()
		return
	}

	if l := utf8.RuneCountInString(username); l==0 {
		this.Data["json"] = this.JsonFormat(210, nil, "姓名不能为空")
		this.ServeJSON()
		return
	}


	if l := utf8.RuneCountInString(work); l==0 {
		this.Data["json"] = this.JsonFormat(210, nil, "工作岗位不能为空")
		this.ServeJSON()
		return
	}

	if l := utf8.RuneCountInString(company); l==0 {
		this.Data["json"] = this.JsonFormat(210, nil, "公司名称不能为空")
		this.ServeJSON()
		return
	}

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

	// 校验手机号码
	reg, _ := regexp.Compile(`^1[34578]\d{9}$`)
	if !reg.MatchString(phone) {
		this.Data["json"] = this.JsonFormat(210, nil, "手机号码校验不通过")
		this.ServeJSON()
		return
	}

	// 校验机器码
	reg, _ = regexp.Compile(`^[0-9A-Za-z]{32}$`)
	if !reg.MatchString(machinecode) {
		this.Data["json"] = this.JsonFormat(210, nil, "机器码校验不通过")
		this.ServeJSON()
		return
	}

	// 校验短信
	if CacheObj.Exists(phone) {
		item, err := CacheObj.Value(phone)
		if err != nil {
			this.Data["json"] = this.JsonFormat(500, nil, "缓存查询错误")
			util.FileLog.Error("从缓存查询手机号错误 msg:%s", err.Error())
			this.ServeJSON()
			return
		}

		smsObj := item.Data().(*SMS)
		if smscode != smsObj.code {
			this.Data["json"] = this.JsonFormat(210, nil, "短信验证码校验不通过")
			this.ServeJSON()
			return
		}
	} else {
		this.Data["json"] = this.JsonFormat(210, nil, "短信验证码校验不通过")
		this.ServeJSON()
		return
	}

	var user models.User
	user.Username = username
	user.Email = email
	user.Phone = phone
	user.Company = company
	user.Work = work
	user.Machinecode = machinecode
	user.VendorNum = ""
	if err := user.GetOne(); err != nil {
		if err != orm.ErrNoRows {
			beego.Error("创业板申请，通过手机查询用户错误"+err.Error())
			this.Data["json"] = this.JsonFormat(500, nil, "用户不存在")
			util.FileLog.Error("通过手机查询用户错误 msg:%s", err.Error())
			this.ServeJSON()
			return
		}
	}else {
		this.Data["json"] = this.JsonFormat(220, nil, "该手机号已提交申请")
		this.ServeJSON()
		return
	}

	if _, err := user.AddUser(); err != nil {
		beego.Error("创业板申请，添加用户失败"+err.Error())
		this.Data["json"] = this.JsonFormat(500, nil, "用户已存在")
		util.FileLog.Error("添加用户错误 msg:%s", err.Error())
		this.ServeJSON()
		return
	}

	// todo 未加密
	v := url.Values{}
	v.Add("machinecode", machinecode)
	oneoaas_domain := beego.AppConfig.String("oneoaas_domain")
	requestUrl := ""
	if len(oneoaas_domain) != 0{
		requestUrl = oneoaas_domain+"/api/apply/license?"+v.Encode()
	}else{
		requestUrl = "http://www.oneoaas.com/api/apply/license?"+v.Encode()
	}

	msg := fmt.Sprintf(`
		申请人:%s
		手机号:%s
		邮箱:%s
		公司:%s
		岗位:%s
		机器码:%s
	`, username, phone, email, company, work, machinecode)

	sendDingdingLinkMessage(msg,"申请Monitor创业板授权码",requestUrl)
	// 给钉钉发消息

	this.Data["json"] = this.JsonFormat(200, nil, "提交申请成功,请耐心等待人工审核,审核结果将发送至邮箱")
	this.ServeJSON()
}

func sendDingdingMessage(msg string) {
	if msg == "" {
		return
	}
	url := beego.AppConfig.String("ding_webhook")
	if url == "" {
		return
	}

	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": msg,
		},
	}

	encoded, _ := json.Marshal(data)

	// Setup our HTTP request
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(encoded))

	if err != nil {
		util.FileLog.Error("钉钉消息发送错误 msg:%s", err.Error())
		return
	}
	request.Header.Add("Content-Type", "application/json")

	// Execute the request
	_, err = client.Do(request)

	if err != nil {
		util.FileLog.Error("钉钉消息发送错误 msg:%s", err.Error())
		return
	}
}

func sendDingdingLinkMessage(msg string,title string ,messageUrl string) {
	if msg == "" {
		return
	}
	url := beego.AppConfig.String("ding_webhook")
	if url == "" {
		return
	}

	data := map[string]interface{}{
		"msgtype": "link",
		"link":map[string]interface{}{
			"text": msg,
			"title": title,
			"picUrl": "",
			"messageUrl": messageUrl,
		},
	}

	encoded, _ := json.Marshal(data)

	// Setup our HTTP request
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(encoded))

	if err != nil {
		util.FileLog.Error("钉钉消息发送错误 msg:%s", err.Error())
		return
	}
	request.Header.Add("Content-Type", "application/json")

	// Execute the request
	_, err = client.Do(request)

	if err != nil {
		util.FileLog.Error("钉钉消息发送错误 msg:%s", err.Error())
		return
	}
}

