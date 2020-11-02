// Author: d1y<chenhonzhou@gmail.com>
// 实验性项目

package tray

import (
	"fmt"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/jsxxzy/inettray/core"
	"github.com/jsxxzy/inettray/core/config"
	"github.com/jsxxzy/inettray/icon"
)

var noLoginText = "登录"

var loginText = "已登录"

var windows = false

// 判断是否为 `windows` 系统
func isWindows() bool {
	return runtime.GOOS == "windows"
}

// New tray
func New(tray Tray) *Tray {
	return &tray
}

// Tray 托盘
type Tray struct {
	Login          *systray.MenuItem // 一层菜单的登录
	ReloadInfo     *systray.MenuItem // 登录=> 刷新信息
	WindowLogin    *systray.MenuItem // 登录=> 登录(windows)
	Ipv4           *systray.MenuItem // 登录=> ipv4
	Flow           *systray.MenuItem // 登录=> 流量
	Duration       *systray.MenuItem // 登录=> 时长
	MonthFlow      *systray.MenuItem // 登录=> 本月已用
	Logout         *systray.MenuItem // 注销
	ConfButton     *systray.MenuItem // 一层菜单的配置
	CurrUser       *systray.MenuItem // 配置=> 用户
	ReloadUser     *systray.MenuItem // 配置=> 刷新本地用户
	OpenConfigCopy *systray.MenuItem // 配置=> 打开配置文件(windows)
	HelpButton     *systray.MenuItem // 一层的帮助
	MQuit          *systray.MenuItem // 退出

	FlowCopy float64 // 流量值
	Ipv4Copy string  // ipv4
}

// SetMenuIcon 设置主菜单图标
func (tray *Tray) SetMenuIcon(iconData []byte) {
	systray.SetIcon(iconData)
}

// SetMenuToolTip 设置菜单
func (tray *Tray) SetMenuToolTip(toolTip string) {
	systray.SetTooltip(toolTip)
}

// SetMenuTitle 设置菜单标题
func (tray *Tray) SetMenuTitle(title string) {
	systray.SetTitle(title)
}

// SetMenuLoginFlag 设置菜单登录
func (tray *Tray) SetMenuLoginFlag(flag bool) {
	var s = "未登录"
	if flag {
		s = loginText
	}
	tray.SetMenuIconFlag(flag)
	systray.SetTitle(s)
}

// SetMenuIconFlag 设置主菜单图标(灰色/彩色)
func (tray *Tray) SetMenuIconFlag(flag bool) {
	// TODO 图标适配windows
	if windows {
		return
	}
	var i = icon.GrayData
	if flag {
		i = icon.RainbowData
	}
	systray.SetIcon(i)
}

// SetInfo 设置信息
func (tray *Tray) SetInfo() error {
	tmpData, tmpErr := core.GetInfo()
	if tmpErr == nil && tmpData.Time != "0" {
		tray.SetFlowMB(tmpData.FlowMB)
		tray.SetIpv4(tmpData.Ipv4)
		tray.Ipv4.SetTitle(tmpData.Ipv4)
		tray.Flow.SetTitle(tmpData.Flow)
		tray.Duration.SetTitle(tmpData.Time)
		tray.Login.SetTitle(loginText)
		tray.SetMonthTotalFlowTitle()
		tray.SetWindowsLoginButtonFlag(false)
		tray.SetMenuLoginFlag(true)
		return nil
	}
	return tray.SetInfoZero()
}

// SetInfoZero 设置初始化信息
func (tray *Tray) SetInfoZero() error {
	tray.Login.SetTitle(noLoginText)
	tray.Ipv4.SetTitle("内网ip")
	tray.Flow.SetTitle("流量")
	tray.Duration.SetTitle("使用时长")
	tray.SetFlowMB(0)
	tray.SetIpv4("")
	tray.SetMonthTotalFlowTitle()
	tray.SetMenuLoginFlag(false)
	return nil
}

// SetWindowsLoginButtonFlag 设置 `windows` 登录按钮显示和隐藏
func (tray *Tray) SetWindowsLoginButtonFlag(flag bool) {
	if flag {
		tray.WindowLogin.Show()
	} else {
		tray.WindowLogin.Hide()
	}
}

// SetWindowsOpenConfig 设置 `windows` 打开配置文件显示和隐藏
func (tray *Tray) SetWindowsOpenConfig(flag bool) {
	if flag {
		tray.OpenConfigCopy.Show()
	} else {
		tray.OpenConfigCopy.Hide()
	}
}

// SetMonthTotalFlowTitle 设置本月已用
func (tray *Tray) SetMonthTotalFlowTitle() {
	var f = core.GetHumanMonthTotalFlow()
	var s = fmt.Sprintf("本月已用: %v", f)
	tray.MonthFlow.SetTitle(s)
}

// LoginWith 登录
func (tray *Tray) LoginWith() error {
	if err := core.Login(); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tray.SetInfo()
}

// CopyIPV4 set `ipv4` to clipboard
func (tray *Tray) CopyIPV4() {
	var ipv4 = tray.Ipv4Copy
	if len(ipv4) >= 1 {
		core.SetClipboard(ipv4)
	}
}

// SetFlowMB 设置流量值
func (tray *Tray) SetFlowMB(f float64) {
	tray.FlowCopy = f
}

// SetIpv4 设置`ipv4`
func (tray *Tray) SetIpv4(ipv4 string) {
	tray.Ipv4Copy = ipv4
}

// SetUsername 设置`username`
func (tray *Tray) SetUsername() {
	var username = core.EasyGetLocalAuthUsername()
	tray.CurrUser.SetTitle(username)
}

// LogoutWith 注销
func (tray *Tray) LogoutWith() {
	var err = core.Logout()
	if err == nil {
		var f = tray.FlowCopy
		core.SetLocalMonthTotalFlow(f)
		tray.SetMonthTotalFlowTitle()
		tray.SetInfoZero()
		tray.SetMenuLoginFlag(false)
		if windows {
			tray.SetWindowsLoginButtonFlag(true)
		}
	}
}

// InitDisableMenu 初始化禁用某些`menu`
func (tray *Tray) InitDisableMenu() {
	tray.MonthFlow.Disable()
	tray.Flow.Disable()
	tray.Duration.Disable()
	tray.CurrUser.Disable()
}

// Init 初始化
func (tray *Tray) Init() {
	tray.InitDisableMenu()

	var data = icon.Data

	if !windows {
		data = icon.RainbowData
	}

	tray.SetMenuIcon(data)
	tray.SetMenuToolTip(config.MainTitleTip)

	tray.SetInfo()
	tray.SetUsername()

	if !windows {
		tray.SetWindowsOpenConfig(false)
		tray.SetWindowsLoginButtonFlag(false)
	}
}

func init() {
	windows = isWindows()
}
