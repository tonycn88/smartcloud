package config

import (
	"fmt"
	"os"

	v2 "gopkg.in/yaml.v2"
)

const configfile = "config.yaml"

type Permission struct {
	Model  string `yaml:"model"`
	Policy string `yaml:"policy"`
}

type TLS struct {
	Cert string `yaml:"servercert"`
	Key  string `yaml:"serverkey"`
}

type Server struct {
	Protocol string `yaml:"protocol"`
	Port     string `yaml:"port"`
	Tls      TLS    `yaml:"tls"`
}

type Authentication struct {
	Type  string `yaml:"type"`
	Basic struct {
		Enable bool `yaml:"enabled"`
	}
	Token struct {
		Enable bool   `yaml:"enabled"`
		Secret string `yaml:"secret"`
	}
	Ca struct {
		Enable bool   `yaml:"enabled"`
		CaCrt  string `yaml:"ca"`
	}
}

// type Certs struct {
// 	Key  string `yaml:"key"`
// 	Crt  string `yaml:"crt"`
// 	Type string `yaml:"type"`
// 	Ca   string `yaml:"ca"`
// }

type Database struct {
	Datatype string `yaml:"dbtype"`
	Url      string `yaml:"url"`
}

type Config struct {
	DbConfig Database `yaml:"database"`
	// CertsConfig Certs    `yaml:"certs"`
	Server     Server         `yaml:"server"`
	Rootdir    string         `yaml:"rootdir"`
	Rootpath   string         `yaml:"rootpath"`
	Auth       Authentication `yaml:"authentication"`
	Permission Permission     `yaml:"permission"`
}

func Read() *Config {
	data, err := os.ReadFile(configfile)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}

	conf := new(Config)
	if err := v2.Unmarshal(data, conf); err != nil { //使用yaml.Unmarshal将yaml文件中的信息反序列化给Config结构体
		fmt.Printf("err: %v\n", err)
		return nil
	}
	fmt.Printf("conf: %v\n", conf)
	// fmt.Printf("conf.SecretKey: %v\n", conf.SecretKey) //通过结构体语法取值

	out, err := v2.Marshal(conf) //序列化为yaml格式文件
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	fmt.Printf("out: %v\n", string(out))
	return conf
}
