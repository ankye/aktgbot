package main

import (
	log "github.com/gonethopper/libs/logs"
)

//AppConfig app基础配置
type AppConfig struct {
	Botkey string `yaml:"botkey"`
}

//Config 配置信息表
type Config struct {
	App *AppConfig `yaml:"app"`
	Log *log.LogConfig
}

//NewConfig 创建配置文件
func NewConfig() *Config {
	c := new(Config)
	c.App = new(AppConfig)

	return c
}
