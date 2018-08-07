package api

import (
	"github.com/muesli/cache2go"
	"oneoaas.com/oneoaas_web/controllers"
	"oneoaas.com/oneoaas_web/util"
	"oneoaas.com/oneoaas_web/util/sms"
	"regexp"
	"time"
)

type NoticeSmsController struct {
	controllers.BaseController
}

type SMS struct {
	// 手机号
	phone string
	// 验证码
	code string
	// 发送时间
	clock int64
	// 发送次数
	count int
}

type ClientSendCount struct {
	sendCount    int
	attemptCount int
}

var CacheObj *cache2go.CacheTable

func init() {
	// 初始化 cache 对象
	CacheObj = cache2go.Cache("oneoaas_web")
}

func (this *NoticeSmsController) Prepare() {
	// 关闭页面渲染
	this.EnableRender = false
}

func (this *NoticeSmsController) SendSmsCode() {
	phone := this.GetString("phone", "")
	if phone == "" {
		this.Data["json"] = this.JsonFormat(211, nil, "手机号不能为空")
		this.ServeJSON()
		return
	} else {
		reg, _ := regexp.Compile(`^1[34578]\d{9}$`)
		if !reg.MatchString(phone) {
			this.Data["json"] = this.JsonFormat(212, nil, "手机号码验证不通过")
			this.ServeJSON()
			return
		}
	}

	var (
		smsObj *SMS
		// 缓存时间
		expire_time = time.Second * 3600
		//短信发送间隔时间
		sms_period = time.Second * 120
		// 短信发送次数
		phoneSendCount int
		// 次数限制
		phoneSendCountLimit = 3
		// client ip
		clientIpSendCount int
		// 同一IP限制发送短信数量
		clientIpSendCountLimit = 20
		// 同一IP尝试发送短信次数
		clientIpAttemptCount int
		clientSend           *ClientSendCount
		// 上次发送时间
		lastClock int64
		// 客户端IP
		clientIP string
		code     = util.RandStr(6, 0)
	)

	clientIP = this.GetClientIp()
	if CacheObj.Exists(clientIP) {
		item, err := CacheObj.Value(clientIP)
		if err != nil {
			util.FileLog.Info("服务器内部错误 msg:%s", err.Error())
			this.Data["json"] = this.JsonFormat(500, nil, "服务内部错误")
			this.ServeJSON()
			return
		}

		clientSend = item.Data().(*ClientSendCount)
		clientIpSendCount = clientSend.sendCount
		clientIpAttemptCount = clientSend.attemptCount
		if clientIpSendCount >= clientIpSendCountLimit {
			this.Data["json"] = this.JsonFormat(215, nil, "同一IP限制短信发送")
			this.ServeJSON()
			util.FileLog.Info("client ip:%s 手机号:%s 限制短信发送,已发送次数：%d 尝试次数:%d", clientIP, phone, clientIpSendCount, clientIpAttemptCount)
			clientSend.attemptCount += 1
			CacheObj.Add(clientIP, time.Second*3600*24*30, clientSend)
			return
		}
	} else {
		clientIpAttemptCount = 0
		clientSend = &ClientSendCount{
			attemptCount: clientIpAttemptCount,
		}
	}

	if CacheObj.Exists(phone) {
		item, err := CacheObj.Value(phone)
		if err != nil {
			util.FileLog.Info("服务器内部错误 msg:%s", err.Error())
			this.Data["json"] = this.JsonFormat(500, nil, "服务内部错误")
			this.ServeJSON()
			return
		}

		smsObj = item.Data().(*SMS)
		phoneSendCount = smsObj.count
		lastClock = smsObj.clock
		smsObj.count += 1
	} else {
		phoneSendCount = 0
		lastClock = 0
		smsObj = &SMS{
			phone: phone,
			count: 1,
			code:  code,
		}
	}

	if lastClock > 0 {
		if time.Now().Sub(time.Unix(lastClock, 0)) < sms_period {
			this.Data["json"] = this.JsonFormat(214, nil, "请求短信接口过于频繁")
			util.FileLog.Info("手机号:%s 请求短信接口过于频繁", smsObj.phone)
			this.ServeJSON()
			clientSend.attemptCount = clientIpAttemptCount + 1
			CacheObj.Add(clientIP, time.Second*3600*24*30, clientSend)
			return
		}
	}

	if phoneSendCount <= phoneSendCountLimit {
		// 向用户发送短信
		sent, err := sms.SendSms(phone, code)
		if err != nil {
			util.FileLog.Info("请求短信接口发生错误 msg:%s", err.Error())
			this.Data["json"] = this.JsonFormat(500, nil, "服务内部错误")
			this.ServeJSON()
			return
		}

		if sent {
			smsObj.clock = time.Now().Unix()
			CacheObj.Add(phone, expire_time, smsObj)
			// 限制恶意消耗短信
			clientSend.sendCount = clientIpSendCount + 1
			clientSend.attemptCount = clientIpAttemptCount + 1
			CacheObj.Add(clientIP, time.Second*3600*24*30, clientSend)
			this.Data["json"] = this.JsonFormat(200, nil, "短信发送成功")
			this.ServeJSON()
			util.FileLog.Info("手机号:%s 验证码:%s 发送成功", smsObj.phone, smsObj.code)
			return
		} else {
			clientSend.attemptCount = clientIpAttemptCount + 1
			CacheObj.Add(clientIP, time.Second*3600*24*30, clientSend)
			util.FileLog.Info("手机号:%s 短信发送失败", smsObj.phone)
			this.Data["json"] = this.JsonFormat(210, nil, "短信发送失败")
			this.ServeJSON()
			return
		}
	} else {
		this.Data["json"] = this.JsonFormat(212, nil, "短信发送次数超过系统限制")
		util.FileLog.Info("手机号:%s 短信发送次数超过系统限制", smsObj.phone)
		this.ServeJSON()
		clientSend.attemptCount = clientIpAttemptCount + 1
		CacheObj.Add(clientIP, time.Second*3600*24*30, clientSend)
		return
	}
}

// 验证短信
func (this *NoticeSmsController) VerifySmsCode() {
	phone := this.GetString("phone", "")
	smscode := this.GetString("smscode", "")

	if phone == "" || smscode == "" {
		this.Data["json"] = this.JsonFormat(211, nil, "手机号或者验证码不能为空")
		this.ServeJSON()
		return
	}

	reg, _ := regexp.Compile(`^1[34578]\d{9}$`)
	if !reg.MatchString(phone) {
		this.Data["json"] = this.JsonFormat(212, nil, "手机号码验证不通过")
		this.ServeJSON()
		return
	}

	if CacheObj.Exists(phone) {
		item, err := CacheObj.Value(phone)
		if err != nil {
			this.Data["json"] = this.JsonFormat(500, nil, "服务内部错误")
			util.FileLog.Info("服务器内部错误 msg:%s", err.Error())
			this.ServeJSON()
			return
		}

		smsObj := item.Data().(*SMS)

		if smscode == smsObj.code {
			CacheObj.Add(phone, time.Second*3600, smsObj)
			this.Data["json"] = this.JsonFormat(200, nil, "短信验证码验证通过")
			this.ServeJSON()
			return
		}
	} else {
		this.Data["json"] = this.JsonFormat(217, nil, "短信未发送或者过期")
		this.ServeJSON()
		return
	}

	this.Data["json"] = this.JsonFormat(200, nil, "短信验证码验证不通过")
	this.ServeJSON()
}
