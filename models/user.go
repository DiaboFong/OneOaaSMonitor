package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"oneoaas.com/oneoaas_web/util"
)

type User struct {
	Id          int64  `json:"userid" orm:"column(userid);pk;auto"`
	Username    string `json:"username" orm:"column(username);size(20)"`
	//加密密码
	Password    string `json:"-" orm:"column(password);size(32)"`
	//明文密码
	Password2   string `json:"password2" orm:"column(password2);size(32)"`
	Email       string `json:"email" orm:"column(email);size(255)"`
	Company     string `json:"company" orm:"column(company);size(20)"`
	Work        string `json:"work" orm:"column(work);size(10)"`
	Phone       string `json:"phone" orm:"column(phone);size(11);unique"`
	VendorNum   string `json:"vendor_num" orm:"column(vendor_num);size(60);null"`
	Machinecode string `json:"machinecode" orm:"column(machinecode);size(32);unique"`
	//用户可能有多个license
	Licenses    []*License `orm:"reverse(many)";null`
}

func checkUser(ctx User) error {
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

func GetUser(phone string) User {
	var user User
	err := Orm.QueryTable("user").Filter("phone", phone).One(&user)
	if err != nil {
		beego.Debug("通过手机查询用户错误", err.Error())
		return user
	}
	return user
}

//通过用户名称 和 密码获取用户
func GetUserByUsernameAndPassword(username string,password string) *User {
	var user User
	err := Orm.QueryTable("user").Filter("username", username).Filter("password", password).One(&user)
	if err != nil {
		beego.Debug("用户登录，查询用户或密码错误", err.Error())
		return nil
	}
	return &user
}

func GetUserById(userId int64) *User {
	var user User
	err := Orm.QueryTable("user").Filter("userid", userId).One(&user)
	if err != nil {
		beego.Error("查询用户错误"+err.Error())
		return nil
	}
	return &user
}

func GetUserByVendorNum(vendorNum string) *User {
	var user User
	err := Orm.QueryTable("user").Filter("vendor_num", vendorNum).One(&user)
	if err != nil {
		beego.Error("通过供应商编号查询用户错误"+err.Error())
		return nil
	}
	return &user
}

//通过手机号查询用户
func GetUserByPhone(phoneNum string) *User {
	var user User
	err := Orm.QueryTable("user").Filter("phone", phoneNum).RelatedSel().One(&user)
	if err != nil {
		beego.Error("通过手机号查询用户错误"+err.Error())
		return nil
	}
	return &user
}

//通过手机号查询用户
func GetUserByMachinecode(machinecode string) *User {
	var user User
	err := Orm.QueryTable("user").Filter("machinecode",machinecode).One(&user)
	Orm.LoadRelated(&user,"Licenses")
	if err != nil {
		beego.Error("通过手机号查询用户错误"+err.Error())
		return nil
	}
	return &user
}

func AddUser(ctx User) (int64, error) {
	id, err := Orm.Insert(&ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *User) GetOne() error {
	return Orm.QueryTable("user").Filter("phone", u.Phone).One(u)
}

func (u *User) AddUser() (int64, error) {
	return Orm.Insert(u)
}

func DelUser(phone string) error {
	_, err := Orm.QueryTable("user").Filter("phone", phone).Delete()
	if err != nil {
		beego.Error("删除user错误", err.Error())
		return err
	}
	return nil
}

//更新用户代理商编号
func (u *User) UpdateUserVendorNum() bool {
	userid := fmt.Sprintf("%d",u.Id)
	beego.Info("授权用户ID是"+userid)
	u.VendorNum = "ONEOAAS-VENDOR-"+userid+util.StrtoMd5(u.Username + u.Email+ u.Company)
	_, err := Orm.Update(u)
	if err != nil{
		return false
	}

	return true
}

func (u *User) UpdateUser() error {
	_, err := Orm.Update(u)
	if err != nil{
		return err
	}

	return nil
}

//获取用户列表
func GetUsers() ([]User,error) {
	var user []User
	_,err := Orm.QueryTable("user").All(&user)
	if err != nil{
		beego.Error("查询用户错误", err.Error())
		return nil,err
	}
	return user,nil
}

