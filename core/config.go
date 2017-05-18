package core

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

var (
	AppConfig map[string]interface{}
)

func init() {
	AppConfig = make(map[string]interface{})
}

// GetConfigFile 读取配置文件路径
func GetConfigFile(filename, confDir string) string {
	return fmt.Sprintf("%s/%s", confDir, filename)
}

// LoadAppConfig 读取配置文件内容
func LoadAppConfig(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if _, err := toml.Decode(string(content), &AppConfig); err != nil {
		panic(err)
	}
}

// GetAppConfig 读取某一项配置
func GetAppConfig(key string) interface{} {
	if value, ok := AppConfig[key]; ok {
		return value
	}
	Logger.Warn("%s not exists in app config")
	return nil
}
