package utils

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string //模式
	HttpPort string //端口
	JwtKey   string //Key

	// 数据库
	DbHost     string //IP
	DbPort     string //端口
	DbUser     string //用户
	DbPassWord string //密码
	DbName     string //库名
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("config.ini not found", err)
	}
	loadServer(file)
	loadDatabase(file)
}

func loadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("89js82js72")
}

func loadDatabase(file *ini.File) {
	DbHost = file.Section("database").Key("DbHost").MustString("127.0.0.1")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("root")
	DbName = file.Section("database").Key("DbName").MustString("test")
}
