// Author: d1y<chenhonzhou@gmail.com>
// 实验性项目

package core

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/jsxxzy/inet"
	"github.com/jsxxzy/inettray/core/config"
	"github.com/pkg/browser"
)

var suffixes [5]string

// =======

// Round round offset
//
// https://gist.github.com/anikitenko/b41206a49727b83a530142c76b1cb82d
func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// 获取格式化好的时间
func getHumanTime(h int) string {
	if h < 60 {
		return fmt.Sprintf("%v分钟", h)
	}
	if h == 60 {
		return "1小时"
	}
	m := h % 60
	var p float64 = 60
	b := float64(h) / p
	c := math.Floor(b)
	return fmt.Sprintf("%v小时%v分钟", c, m)
}

// getHumanFlow 转换流量格式转为阳间格式
//
// https://gist.github.com/anikitenko/b41206a49727b83a530142c76b1cb82d
func getHumanFlow(f float64) string {
	size := f * 1024 * 1024 // This is in bytes
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	var result = strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
	return result
}

// Info 数据
type Info struct {
	Time string
	Flow string
	Ipv4 string
}

// GetInfo 获取数据
func GetInfo() (Info, error) {
	info, err := inet.QueryInfo()
	if err != nil {
		return Info{}, errors.New("查询信息失败")
	}
	xTime, _ := strconv.Atoi(info.Time)
	var time = getHumanTime(xTime)
	var flow = getHumanFlow(info.Flow)
	return Info{
		Time: time,
		Flow: flow,
		Ipv4: info.V4ip,
	}, nil
}

// OpenHelp 打开帮助
func OpenHelp() error {
	return browser.OpenURL(config.GithubHelp)
}

// OpenConfig 打开配置文件
func OpenConfig() error {
	return browser.OpenFile(config.ConfigFile)
}

// EasyGetLocalAuthUsername 用最简单的方法获取用户名
//
// 不安全的方法, 慎用!!
func EasyGetLocalAuthUsername() string {
	auth, err := config.GetConfigFile()
	if err != nil || len(auth.Username) <= 1 {
		return "未知"
	}
	return auth.Username
}

// Login 登录
func Login() error {
	auth, err := config.GetConfigFile()
	if len(auth.Username) <= 1 && len(auth.Password) <= 1 {
		return errors.New("账号密码配置错误")
	}
	if err != nil {
		return err
	}
	loginInfo, err := inet.Login(auth.Username, auth.Password)
	return loginInfo.Error()
}

// Logout 注销
func Logout() error {
	return inet.Logout()
}
