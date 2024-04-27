package conf

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func conversionLogLevel(logLevel string) int {
	switch logLevel {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func init() {
	InitConf()
	config := make(map[string]interface{})
	config["filename"] = GameConf.LogPath
	config["level"] = conversionLogLevel(GameConf.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		logs.Error("marsha1 failed,err", err)
		fmt.Println("marsha1 failed,err", err)
	}
	fmt.Println("marsha1 failed", string(configStr))

	err = logs.SetLogger(logs.AdapterFile, string(configStr))
	if err != nil {
		logs.Error("init logger failed:%v", err)
	}
}
