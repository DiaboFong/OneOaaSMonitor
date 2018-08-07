package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"time"
	"github.com/shirou/gopsutil/cpu"
	"fmt"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
)

var GlobalSessions *session.Manager

var fontKinds = [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}

func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

func RandStr(size int, kind int) string {
	ikind, result := kind, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			ikind = rand.Intn(3)
		}
		scope, base := fontKinds[ikind][0], fontKinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var SIGN_KEY_SE = "HJKALSDHJDHUQWGFHUBONEOAASXZNCBMNZBHVCNQAJHSGDFKQHG"
var SIGN_KEY_EE = "ONEOAASHJKALSDHJDHUQWGFHUBONEOAASXZNCBMNZBHVCNQAJHSGDFKQHGONEOAAS"

//创建Lisence
func CreateLisence(hostinfo string,days int,hours int,years int,version string )(string,time.Time){
	beego.Info("Lisence类型是:"+version)
	if hostinfo == "" {
		fmt.Println("请输入字符串")
		os.Exit(1)
	}
	var expire_time time.Duration = time.Hour
	if days > 0 {
		beego.Info(fmt.Sprintf("申请时长是 %d 天",days))
		expire_time *= time.Duration(24 * days)
	} else if hours > 0 {
		beego.Info(fmt.Sprintf("申请时长是 %d 小时",hours))
		expire_time *= time.Duration(hours)
	} else if years > 0 {
		beego.Info(fmt.Sprintf("申请时长是 %d 年",hours))
		expire_time *= time.Duration(365 * 24 * years)
	} else {
		beego.Info(fmt.Sprintf("默认申请时长是 %d 半年",24 * 30 * 6))
		expire_time *= time.Duration(24 * 30 * 6)
	}
	tokenString := CreateToken(hostinfo, version, expire_time)

	beego.Info("以下为授权码")
	beego.Info(tokenString)
	beego.Info("有效期:%s~%s\n", time.Now().String(), time.Now().Add(expire_time).String())
	return tokenString,time.Now().Add(expire_time)
}

func CreateToken(hostinfo, version string, expire time.Duration) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	hostinfo = StrtoMd5(hostinfo)
	beego.Info("创建Token的")
	claims["hostinfo"] = hostinfo
	claims["exp"] = time.Now().Add(expire).Unix()
	var signKey []byte
	if version == "se" {
		signKey = []byte(SIGN_KEY_SE + hostinfo)
	} else {
		signKey = []byte(SIGN_KEY_EE + hostinfo)
	}

	tokenString, _ := token.SignedString(signKey)
	return tokenString
}

func GetHostInfoMd5() (string, error) {
	host, err := host.Info()
	if err != nil {
		fmt.Printf("获取设备信息失败 %v\n", err)
		return "", err
	}

	hostid := host.HostID
	cpus, err := cpu.Info()

	if err != nil {
		fmt.Printf("获取设备信息失败 %v\n", err)
		return "", err
	}

	if len(cpus) == 0 {
		return "", errors.New("cpu不存在")
	}

	modelName := cpus[0].ModelName
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("获取设备信息失败 %v\n", err)
		return "", err
	}

	var macAddr string
	for _, vv := range interfaces {
		if vv.HardwareAddr != "" {
			macAddr = vv.HardwareAddr
			break
		}
	}

	return StrtoMd5(hostid + modelName + macAddr), nil
}

func StrtoMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}


func GenerateSalt() string {
	const randomLength = 16
	var salt []byte
	var asciiPad int64
	asciiPad = 32

	for i := 0; i < randomLength; i++ {
		salt = append(salt, byte(rand.Int63n(94)+asciiPad))
	}
	return string(salt)
}


