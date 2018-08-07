package models

import (
	"errors"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type SmsCode struct {
	//Emails []map[string]interface{} `json:"Emails"`
	Id     int64  `json:"smscodeid" orm:"column(smscodeid);pk;auto"`
	Phone  string `json:"phone" orm:"column(phone);size(255);null"`
	Code   string `json:"code" orm:"column(code);size(10);null"`
	Count  int64  `json:"count" orm:"column(count);size(10);0"`
	IsTrue int64  `json:"istrue" orm:"column(istrue);size(10);0"`
}

func checkSmsCode(ctx SmsCode) error {
	valid := validation.Validation{}
	b, _ := valid.Valid(ctx)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

func GetSmsCodeByPhone(phone string) *SmsCode {
	smsCode := new(SmsCode)
	err := Orm.QueryTable("sms_code").Filter("phone", phone).One(smsCode)

	if err != nil {
		beego.Debug("查询sms_code错误", err.Error())
		return smsCode
	}

	return smsCode
}

func AddSmsCode(ctx SmsCode) (int64, error) {
	id, err := Orm.Insert(&ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func PutSmsCode(ctx *SmsCode) (err error) {
	_, err = Orm.Update(ctx)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DelSmsCode(phone string) error {
	_, err := Orm.QueryTable("sms_code").Filter("phone", phone).Delete()
	if err != nil {
		beego.Error("删除sms_code错误", err.Error())
		return err
	}
	return nil
}
