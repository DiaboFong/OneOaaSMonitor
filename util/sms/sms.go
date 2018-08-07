package sms

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"time"

	"github.com/astaxie/beego"
)

type Code struct {
	Code string `json:"%code%"`
}
type Sms struct {
	MsgType    int64  `json:"msgType"`
	Phone      string `json:"phone"`
	SmsUser    string `json:"smsUser"`
	TemplateId int    `json:"templateId"`
	Vars       Code   `json:"vars"`
	Signature  string `json:"signature"`
}

var (
	sms_url  = "http://www.sendcloud.net/smsapi/send"
	sms_user = beego.AppConfig.String("sms_user")
	sms_key  = beego.AppConfig.String("sms_key")
	sms_tid  = beego.AppConfig.String("sms_tid")
)

func SendSms(phone string, code string) (bool, error) {
	vars := map[string]string{
		`%code%`: code,
	}
	jsonVars, _ := json.Marshal(&vars)
	PostParams := url.Values{
		`smsUser`:    {sms_user},
		`templateId`: {sms_tid},
		`msgType`:    {`0`},
		`phone`:      {phone},
		`vars`:       {string(jsonVars)},
	}
	paramsKeyS := make([]string, 0, len(PostParams))
	for k, _ := range PostParams {
		paramsKeyS = append(paramsKeyS, k)
	}
	sort.Strings(paramsKeyS)
	sb := sms_key + "&"
	for _, v := range paramsKeyS {
		sb += fmt.Sprintf("%s=%s&", v, PostParams.Get(v))
	}
	sb += sms_key
	hashMd5 := md5.New()
	io.WriteString(hashMd5, sb)
	sign := fmt.Sprintf("%x", hashMd5.Sum(nil))
	PostParams.Add("signature", sign)

	PostBody := bytes.NewBufferString(PostParams.Encode())
	http.DefaultClient.Timeout = 5 * time.Second
	ResponseHandler, err := http.Post(sms_url, "application/x-www-form-urlencoded", PostBody)
	if err != nil {
		return false, err
	}
	defer ResponseHandler.Body.Close()
	BodyByte, _ := ioutil.ReadAll(ResponseHandler.Body)
	v := string(BodyByte)

	reg, _ := regexp.Compile("请求成功")
	result := reg.MatchString(v)
	return result, nil
}
