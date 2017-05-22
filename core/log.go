package core

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/logs"
)

// LogConfig 日志配置
type LogConfig struct {
	Log_adapter_console        bool
	Log_console_level          int
	Log_file                   string
	Log_file_level             int
	Log_enable_func_call_depth bool
	Log_async                  bool
	Log_chan_length            int
	Log_maxlines               int
	Log_maxsize                int
	Log_daily                  bool
	Log_maxdays                int
	Log_rotate                 bool
	Log_multifile              bool
	Log_separate               []string
}

// Logger 日志对象
var Logger *logs.BeeLogger

// LoadLoggerConfig 加载日志配置
func LoadLoggerConfig(file string) {
	var logConfig LogConfig

	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if _, err := toml.Decode(string(content), &logConfig); err != nil {
		panic(err)
	}

	// 加载log
	Logger = logs.NewLogger()
	// 设置异步输出
	if logConfig.Log_async {
		Logger.Async(int64(logConfig.Log_chan_length))
	}
	// 设置输出文件名、文件行数
	if logConfig.Log_enable_func_call_depth {
		Logger.EnableFuncCallDepth(true)
	}
	// 设置控制台输出
	if logConfig.Log_adapter_console {
		consoleConfig := make(map[string]int)
		consoleConfig["level"] = logConfig.Log_console_level
		byt, _ := json.Marshal(consoleConfig)
		Logger.SetLogger(logs.AdapterConsole, string(byt))
	}

	fileConfig := make(map[string]interface{})
	fileConfig["filename"] = logConfig.Log_file
	fileConfig["maxlines"] = logConfig.Log_maxlines
	fileConfig["maxsize"] = logConfig.Log_maxsize
	fileConfig["daily"] = logConfig.Log_daily
	fileConfig["maxdays"] = logConfig.Log_maxdays
	fileConfig["rotate"] = logConfig.Log_rotate
	if logConfig.Log_multifile {
		fileConfig["separate"] = logConfig.Log_separate
		byt, _ := json.Marshal(fileConfig)
		Logger.SetLogger(logs.AdapterMultiFile, string(byt))
	} else {
		byt, _ := json.Marshal(fileConfig)
		Logger.SetLogger(logs.AdapterFile, string(byt))
	}

	// 据说不这样做，会有一些性能问题
	Logger.SetLevel(logConfig.Log_file_level)
}
