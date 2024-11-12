package config

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func InitConfig(path string) GlobalConfig {
	var gc GlobalConfig
	file, err := os.ReadFile(path)
	if err != nil {
		panic("missing config file: " + path)
	}

	err = yaml.NewDecoder(bytes.NewReader(file)).Decode(&gc)
	if err != nil {
		panic("error parse config file: " + err.Error())
	}

	gc.print()
	return gc
}

type GlobalConfig struct {
	HostName string `yaml:"hostName"`
	HttpConf struct {
		ListenAddr   string `yaml:"listenAddr"`
		CertHostname string `yaml:"certHostname"`
		TimeOut      string `yaml:"timeOut"`
		IdleTimeout  string `yaml:"idleTimeout"`
		CertFile     string `yaml:"certFile"`
		KeyFile      string `yaml:"keyFile"`
		JWTSecret    string `yaml:"jwtSecret"`
	} `yaml:"httpConf"`
	DbConf struct {
		DbType string `yaml:"dbType"`
		DbURL  string `yaml:"dbURL"`
	} `yaml:"dbConf"`
}

func (c *GlobalConfig) print() {
	fmt.Println("Global Config:")
	fmt.Println("HostName:", c.HostName)

	fmt.Println("HTTP Configuration:")
	fmt.Println("  ListenAddr:", c.HttpConf.ListenAddr)
	fmt.Println("  CertHostname:", c.HttpConf.CertHostname)
	fmt.Println("  TimeOut:", c.HttpConf.TimeOut)
	fmt.Println("  IdleTimeout:", c.HttpConf.IdleTimeout)
	fmt.Println("  CertFile:", c.HttpConf.CertFile)
	fmt.Println("  KeyFile:", c.HttpConf.KeyFile)

	fmt.Println("Database Configuration:")
	fmt.Println("  DbType:", c.DbConf.DbType)
	fmt.Println("  DbURL:", c.DbConf.DbURL)
}
