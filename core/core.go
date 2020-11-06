// Author: d1y<chenhonzhou@gmail.com>
// 实验性项目

package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jinzhu/now"
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
//
// error: https://stackoverflow.com/a/26129063
func getHumanFlow(f float64) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
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
	Time   string
	Flow   string
	FlowMB float64
	Ipv4   string
}

// GetInfo 获取数据
func GetInfo() (Info, error) {
	info, err := inet.QueryInfo()
	if err != nil || info.Error() != nil {
		return Info{}, errors.New("查询信息失败")
	}
	xTime, _ := strconv.Atoi(info.Time)
	var time = getHumanTime(xTime)
	var flow = getHumanFlow(info.Flow)
	if len(flow) == 0 {
		flow = "0kb"
	}
	return Info{
		Time:   time,
		Flow:   flow,
		FlowMB: info.Flow,
		Ipv4:   info.V4ip,
	}, nil
}

// OpenHelp 打开帮助
func OpenHelp() error {
	return browser.OpenURL(config.GithubHelp)
}

// OpenConfig 打开配置文件
func OpenConfig() error {
	if runtime.GOOS == "windows" {
		tmpRun := exec.Command("notepad", config.ConfigFile)
		return tmpRun.Run()
	}
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

// SetClipboard 设置剪贴板
func SetClipboard(text string) error {
	return clipboard.WriteAll(text)
}

// FlowConfigFile 流量文件结构体
type FlowConfigFile struct {
	file *os.File
}

// Get 获取流量
func (ff *FlowConfigFile) Get() (float64, error) {
	byteData, err := ioutil.ReadFile(ff.file.Name())
	if err != nil {
		return 0, err
	}
	var s = string(byteData)
	s = strings.TrimSuffix(s, "\n")
	byteFloat, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return byteFloat, nil
}

// Clean 清零
//
// 如果到了月底就要清零一次流量
//
// TODO: 可能以后会做几个月加在一起吧
//
// (估计不会做的这么简单吧)
//
func (ff *FlowConfigFile) Clean() error {
	var zeroValue float64 = 0
	return ff.Write(zeroValue)
}

// 写入流量操作
func (ff *FlowConfigFile) Write(flow float64) error {
	var value = fmt.Sprintf("%f", flow)
	var err = ioutil.WriteFile(ff.file.Name(), []byte(value), 0777)
	return err
}

// Add 添加流量
func (ff *FlowConfigFile) Add(val float64) (float64, error) {
	var flow, err = ff.Get()
	if err != nil {
		return val, err
	}
	var outputValue = flow + val
	return val, ff.Write(outputValue)
}

// GetHumanMonthTotalFlow 获取本月共用了多少流量
func GetHumanMonthTotalFlow() string {
	vendor := FlowConfigFile{
		file: config.FlowFile,
	}
	f, e := vendor.Get()
	if e != nil || f == 0 {
		fmt.Println(e)
		return "0kb"
	}
	var s = getHumanFlow(f)
	return s
}

// SetLocalMonthTotalFlow 设置本月使用总流量
//
// 获取文件修改时间: https://blog.csdn.net/liangguangchuan/article/details/78952979
//
// 获取每个月的第一天: https://www.coder.work/article/25353
//
func SetLocalMonthTotalFlow(timeLessFlow float64) (float64, error) {

	// fmt.Println("1. 创建文件句柄")
	var vendor = FlowConfigFile{
		file: config.FlowFile,
	}

	// fmt.Println("2. 获取文件句柄")
	flowFileInfo, err := os.Stat(config.ConfigFile)
	if err != nil {
		return 0, err
	}
	// fmt.Println("3. 比对时间")
	var changeTime = flowFileInfo.ModTime().Unix()
	var startTime = now.BeginningOfMonth().Unix()
	var curr = now.BeginningOfDay()
	var currday = curr.Unix()
	var nextDay = curr.AddDate(0, 0, 1).Unix()

	var pushFlow = timeLessFlow

	if currday == startTime {
		// fmt.Println("当前为本月第一天")
		if changeTime <= nextDay {
			// fmt.Println("文件修改在当天, 不需要清零")
		} else {
			vendor.Clean()
			pushFlow = 0
			// fmt.Println("文件写入的时间为上月, 需要清零")
		}
	}

	// fmt.Println("4. 增加流量")

	return vendor.Add(pushFlow)
}
