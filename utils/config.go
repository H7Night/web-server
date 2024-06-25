package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string
	JwtKey   string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
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
