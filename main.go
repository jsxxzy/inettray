// Author: d1y<chenhonzhou@gmail.com>
// 实验性项目

package main

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/jsxxzy/inettray/core"
	"github.com/jsxxzy/inettray/core/config"
	"github.com/jsxxzy/inettray/core/tray"
	"github.com/jsxxzy/lockfile"
)

func main() {
	singleApp, code, _ := lockfile.NewSingleApp("inet")
	if code == lockfile.AppRunOtherProcess {
		fmt.Println("[inet] 只允许启动一个进程")
		return
	}
	singleApp.Free(func() {
		fmt.Println("[inet] 已退出程序")
		os.Exit(0)
	})
	systray.Run(onReady, onExit)
}

func onReady() {

	loginButton := systray.AddMenuItem("登录", config.LoginTip)
	reloadInfoButton := loginButton.AddSubMenuItem("刷新信息", "获取使用信息")
	windowsLoginButton := loginButton.AddSubMenuItem("登录", "")
	ipv4 := loginButton.AddSubMenuItem("内网ip", "单击可复制ipv4地址")
	flow := loginButton.AddSubMenuItem("流量", "使用流量")
	duration := loginButton.AddSubMenuItem("使用时长", "使用时长")
	monthFlow := loginButton.AddSubMenuItem("本月已用", "本月已用")
	logoutButton := systray.AddMenuItem("注销", config.LogoutTip)
	confButton := systray.AddMenuItem("配置", config.ConfigTip)
	currUser := confButton.AddSubMenuItem("", "当前的账号")
	reloadUser := confButton.AddSubMenuItem("刷新本地用户", "刷新本地用户")
	openConfigCopy := confButton.AddSubMenuItem("打开配置文件", "打开配置文件")
	helpButton := systray.AddMenuItem("帮助", config.HelpTip)
	mQuit := systray.AddMenuItem("退出", config.ExitTip)

	var App = tray.New(tray.Tray{
		Login:          loginButton,
		ReloadInfo:     reloadInfoButton,
		WindowLogin:    windowsLoginButton,
		Ipv4:           ipv4,
		Flow:           flow,
		Duration:       duration,
		MonthFlow:      monthFlow,
		Logout:         logoutButton,
		ConfButton:     confButton,
		CurrUser:       currUser,
		ReloadUser:     reloadUser,
		OpenConfigCopy: openConfigCopy,
		HelpButton:     helpButton,
		MQuit:          mQuit,
	})

	App.Init()

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	go func() {
		for {
			select {
			case <-monthFlow.ClickedCh: // 本月已用
				App.SetMonthTotalFlowTitle()
			case <-ipv4.ClickedCh: // ipv4
				App.CopyIPV4()
			case <-reloadInfoButton.ClickedCh: // 获取使用信息
				App.SetInfo()
			case <-reloadUser.ClickedCh: // 获取本地用户
				App.SetUsername()
			case <-helpButton.ClickedCh: // 帮助
				core.OpenHelp()
			case <-openConfigCopy.ClickedCh:
				core.OpenConfig()
			case <-confButton.ClickedCh: // 配置
				core.OpenConfig()
			case <-windowsLoginButton.ClickedCh: // 登录
				App.LoginWith()
			case <-loginButton.ClickedCh: // 登录
				App.LoginWith()
			case <-logoutButton.ClickedCh: // 注销
				App.LogoutWith()
			}
		}
	}()

}

func onExit() {
	fmt.Println("已退出")
}
