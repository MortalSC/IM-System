package config

import (
	libLog "github.com/MortalSC/IM-System/lib/log"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Cfg = InitConfig()

type Config struct {
	viper  *viper.Viper
	SrvCfg *ServerConfig
	GC     *GrpcConfig
}

func InitConfig() *Config {
	v := viper.New()
	conf := &Config{viper: v}

	workDir, _ := os.Getwd()

	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath("/etc/IM-System/auth-service/user")
	conf.viper.AddConfigPath(workDir + "/auth-service/config")

	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	// 读取服务配置
	conf.InitServerConfig()
	// 读取日志配置
	conf.InitZapLog()
	// 读取redis配置
	conf.InitRedisOptions()
	// 读取grpc配置
	conf.InitGrpcConfig()

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

func (c *Config) InitRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"), // no password set
		DB:       c.viper.GetInt("redis.db"),          // use default DB
	}
}

type GrpcConfig struct {
	Name string
	Addr string
}

func (c *Config) InitGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	c.GC = gc
}
