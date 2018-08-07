package models

import (
	"time"
	"github.com/astaxie/beego"
)

type License struct {
	Id            int64      `json:"id"           orm:"column(id);pk;auto"`
	User          *User	     `json:"user"         orm:"rel(fk)"`
	ExpireDate    time.Time  `json:"expire_date"  orm:"column(expire_date);type(datetime)"`
	ApplyDate     time.Time  `json:"apply_date"   orm:"column(apply_date);type(datetime)"`
	Duration      string 	 `json:"duration"     orm:"column(duration);size(10)"`
	LicenseKey    string 	 `json:"license_key"  orm:"column(license_key);size(200)"`
}

func (license *License) AddLicense() (int64, error) {
	return Orm.Insert(license)
}

func GetLicense(license_key string) *License {
	var license License
	err := Orm.QueryTable("license").Filter("license_key", license_key).One(&license)
	if err != nil {
		beego.Error("通过license_key查询license错误"+err.Error())
		return nil
	}
	return &license

}

func GetLicenses(who string) ([]License,error) {
	var licenses []License
	if who=="oneoaas"{
		_,err := Orm.QueryTable("license").RelatedSel("user").All(&licenses)
		if err != nil{
			beego.Error("查询 oneoaas License错误", err.Error())
			return nil,err
		}
	}else{
		_,err := Orm.QueryTable("license").Filter("User__Username",who).RelatedSel("user").All(&licenses)
		if err != nil{
			beego.Error("查询"+who +"License错误", err.Error())
			return nil,err
		}
	}

	return licenses,nil
}

func (license *License) DeleteLicense() (int64, error) {
	return Orm.Insert(license)
}





