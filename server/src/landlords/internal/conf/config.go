package conf

import (
	"os"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

type GameConfInfo struct {
	Version  string
	HttpPort int
	LogPath  string
	LogLevel string
}

var (
	GameConf = &GameConfInfo{}
)

func InitConf() (err error) {
	environment := os.Getenv("ENV")
	if environment != "dev" && environment != "testing" && environment != "product" {
		environment = "product"
	}
	logs.Info("the running environment is : %s", environment)
	conf, err := config.NewConfig("ini", "internal/conf/app.conf")
	if err != nil {
		logs.Error("new conf failed ,err : %v", err)
		return
	}
	environment += "::"

	GameConf.Version = conf.String(environment + "version")
	if len(GameConf.Version) == 0 {
		GameConf.Version = "1.0.0"
	}
	logs.Debug("read conf success , Version : %v", GameConf.Version)

	GameConf.HttpPort, err = conf.Int(environment + "http_port")
	if err != nil {
		logs.Error("config http_port failed,err: %v", err)
		return
	}

	logs.Debug("read conf success , httpPort : %v", GameConf.HttpPort)

	//todo log config
	GameConf.LogPath = conf.String(environment + "log_path")
	if len(GameConf.LogPath) == 0 {
		GameConf.LogPath = "internal/logs/game.log"
	}

	logs.Debug("read conf success , LogPath :  %v", GameConf.LogPath)
	GameConf.LogLevel = conf.String(environment + "log_level")
	if len(GameConf.LogLevel) == 0 {
		GameConf.LogLevel = "debug"
	}
	logs.Debug("read conf success , LogLevel :  %v", GameConf.LogLevel)

	return
}
