// Author: d1y<chenhonzhou@gmail.com>
// 实验性项目

package main

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

func main() {
	systray.Run(onReady, onExit)
}

var windows = false

// 判断是否为 `windows` 系统
func isWindows() bool {
	return runtime.GOOS == "windows"
}

func onReady() {

	tmpData, tmpErr := core.GetInfo()
	localAuthUsername := core.EasyGetLocalAuthUsername()

	systray.SetIcon(icon.Data)

	systray.SetTooltip(config.MainTitleTip)
	// systray.SetTitle(config.AppName)
	loginButton := systray.AddMenuItem(noLoginText, config.LoginTip)

	reloadInfoButton := loginButton.AddSubMenuItem("刷新信息", "获取使用信息")

	wnidowsLoginButton := loginButton.AddSubMenuItem("登录", "")

	if !windows || tmpErr == nil {
		wnidowsLoginButton.Hide()
	}

	// ip
	ipv4 := loginButton.AddSubMenuItem("内网ip", "单击可复制ipv4地址")
	// 流量
	flow := loginButton.AddSubMenuItem("流量", "使用流量")
	// 使用时长
	duration := loginButton.AddSubMenuItem("使用时长", "使用时长")

	if tmpErr == nil {
		ipv4.SetTitle(tmpData.Ipv4)
		flow.SetTitle(tmpData.Flow)
		duration.SetTitle(tmpData.Time)
		loginButton.SetTitle(loginText)
	}

	flow.Disable()
	duration.Disable()

	logoutButton := systray.AddMenuItem("注销", config.LogoutTip)

	confButton := systray.AddMenuItem("配置", config.ConfigTip)
	currUser := confButton.AddSubMenuItem(localAuthUsername, "当前的账号")
	currUser.Disable()

	reloadUser := confButton.AddSubMenuItem("刷新本地用户", "刷新本地用户")

	openConfigCopy := confButton.AddSubMenuItem("打开配置文件", "打开配置文件")

	if !windows {
		openConfigCopy.Hide()
	}

	helpButton := systray.AddMenuItem("帮助", config.HelpTip)
	mQuit := systray.AddMenuItem("退出", config.ExitTip)

	var ipv4copy = ""

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	go func() {
		for {
			select {
			case <-ipv4.ClickedCh: // ipv4
				if len(ipv4copy) >= 1 {
					core.SetClipboard(ipv4copy)
				}
			case <-reloadInfoButton.ClickedCh: // 获取使用信息
				data, err := core.GetInfo()
				if err != nil {
					loginButton.SetTitle(noLoginText)
					ipv4copy = ""
				} else {
					loginButton.SetTitle(loginText)
					ipv4.SetTitle(data.Ipv4)
					ipv4copy = data.Ipv4
					flow.SetTitle(data.Flow)
					duration.SetTitle(data.Time)
				}
			case <-reloadUser.ClickedCh: // 获取本地用户
				var username = core.EasyGetLocalAuthUsername()
				currUser.SetTitle(username)
			case <-helpButton.ClickedCh: // 帮助
				core.OpenHelp()
			case <-openConfigCopy.ClickedCh:
				core.OpenConfig()
			case <-confButton.ClickedCh: // 配置
				core.OpenConfig()

			case <-wnidowsLoginButton.ClickedCh: // 登录
				// TODO: 重复代码
				if err := core.Login(); err != nil {
					fmt.Println(err.Error())
				} else {
					tmpData, tmpErr := core.GetInfo()
					if tmpErr == nil && tmpData.Time != "0" {
						ipv4.SetTitle(tmpData.Ipv4)
						flow.SetTitle(tmpData.Flow)
						duration.SetTitle(tmpData.Time)
						loginButton.SetTitle(loginText)
					}
					wnidowsLoginButton.Hide()
				}
			case <-loginButton.ClickedCh: // 登录
				if err := core.Login(); err != nil {
					fmt.Println(err.Error())
				} else {
					tmpData, tmpErr := core.GetInfo()
					if tmpErr == nil && tmpData.Time != "0" {
						ipv4.SetTitle(tmpData.Ipv4)
						flow.SetTitle(tmpData.Flow)
						duration.SetTitle(tmpData.Time)
						loginButton.SetTitle(loginText)
					}
				}
			case <-logoutButton.ClickedCh: // 注销
				if core.Logout() == nil {
					ipv4.SetTitle("内网ip")
					ipv4copy = ""
					flow.SetTitle("流量")
					duration.SetTitle("使用时长")
					loginButton.SetTitle(noLoginText)
					if windows {
						wnidowsLoginButton.Show()
					}
				}
			}
		}
	}()

	// Sets the icon of a menu item. Only available on Mac and Windows.
	// mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}

func init() {
	windows = isWindows()
}
