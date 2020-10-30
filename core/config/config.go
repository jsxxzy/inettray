// Author: d1y<chenhonzhou@gmail.com>
// 实验性项目
//
// 参考项目: https://github.com/tj/go-config

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jsxxzy/inettray/core/ini"
)

// MainTitleTip 大标题介绍
const MainTitleTip = "职院校园网"

// AppName app name
const AppName = "inet"

// LoginTip 登录
const LoginTip = "登录校园网"

// LogoutTip 注销
const LogoutTip = "注销校园网"

// ConfigTip 配置
const ConfigTip = "配置你的账号密码"

// HelpTip 帮助
const HelpTip = ""

// ExitTip 退出
const ExitTip = "退出程序并不会注销"

// GithubHelp 帮助
const GithubHelp = "https://github.com/jsxxzy/inet"

// ConfigFile 配置文件
var ConfigFile = ""

var configFileName = ".inet.conf"

// Auth 鉴权
type Auth struct {
	// Username 账号
	Username string
	// password 密码
	Password string
}

// 获取用户`home`目录
func getHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

// 初始化配置文件
func initConfigFile() error {
	if exists(ConfigFile) {
		return nil
	}
	var initStr = `
# 请填入账号和密码即可
# https://github.com/jsxxzy/inettray

username = 
password = `
	return ioutil.WriteFile(ConfigFile, []byte(initStr), 0777)
}

// check file/dir exists
//
// https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetConfigFile 获取 `config` 配置文件
func GetConfigFile() (Auth, error) {
	conf, err := ini.NewConfig(ConfigFile) // string(byteData))
	if err != nil {
		return Auth{}, err
	}
	u := conf.String("username")
	p := conf.String("password")
	return Auth{
		Username: u,
		Password: p,
	}, nil
}

func init() {
	homeDir, err := getHomeDir()
	if err != nil {
		panic(err)
	}
	ConfigFile = filepath.Join(homeDir, configFileName)
	initConfigFile()
}
