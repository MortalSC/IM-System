package config

import (
	libLog "github.com/MortalSC/IM-System/lib/log"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Cfg = InitConfig()

type Config struct {
	viper  *viper.Viper
	SrvCfg *ServerConfig
}

func InitConfig() *Config {
	v := viper.New()
	conf := &Config{viper: v}

	workDir, _ := os.Getwd()

	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath("/etc/IM-System/api-center/user")
	conf.viper.AddConfigPath(workDir + "/api-center/config")

	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	// 读取服务配置
	conf.InitServerConfig()
	// 读取日志配置
	conf.InitZapLog()

	return conf
}

// ServerConfig 服务配置
type ServerConfig struct {
	Name string
	Addr string
}

func (c *Config) InitServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SrvCfg = sc
}

func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	lc := &libLog.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := libLog.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}
