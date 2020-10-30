// Author: d1y<chenhonzhou@gmail.com>
// 实验性项目

package main

import (
	"fmt"

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

func onReady() {
	systray.SetIcon(icon.Data)

	systray.SetTooltip(config.MainTitleTip)
	systray.SetTitle(config.AppName)
	loginButton := systray.AddMenuItem(noLoginText, config.LoginTip)

	reloadInfoButton := loginButton.AddSubMenuItem("刷新信息", "获取使用信息")

	// ip
	ipv4 := loginButton.AddSubMenuItem("内网ip", "内网ipv4")
	// 流量
	flow := loginButton.AddSubMenuItem("流量", "使用流量")
	// 使用时长
	duration := loginButton.AddSubMenuItem("使用时长", "使用时长")

	flow.Disable()
	duration.Disable()

	logoutButton := systray.AddMenuItem("注销", config.LogoutTip)

	confButton := systray.AddMenuItem("配置", config.ConfigTip)
	currUser := confButton.AddSubMenuItem("未知", "当前的账号")
	currUser.Disable()
	reloadUser := confButton.AddSubMenuItem("刷新本地用户", "刷新本地用户")

	helpButton := systray.AddMenuItem("帮助", config.HelpTip)
	mQuit := systray.AddMenuItem("退出", config.ExitTip)

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	go func() {
		for {
			select {
			case <-reloadInfoButton.ClickedCh: // 获取使用信息
				data, err := core.GetInfo()
				if err != nil {
					loginButton.SetTitle(noLoginText)
				} else {
					loginButton.SetTitle(loginText)
					ipv4.SetTitle(data.Ipv4)
					flow.SetTitle(data.Flow)
					duration.SetTitle(data.Time)
				}
			case <-reloadUser.ClickedCh: // 获取本地用户
				var username = core.EasyGetLocalAuthUsername()
				currUser.SetTitle(username)
			case <-helpButton.ClickedCh: // 帮助
				core.OpenHelp()
			case <-confButton.ClickedCh: // 配置
				core.OpenConfig()
			case <-loginButton.ClickedCh: // 登录
				if err := core.Login(); err != nil {
					fmt.Println(err.Error())
				}
			case <-logoutButton.ClickedCh: // 注销
				core.Logout()
			}
		}
	}()

	// Sets the icon of a menu item. Only available on Mac and Windows.
	// mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}
