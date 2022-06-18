package conf

import (
	"github.com/BurntSushi/toml"
)

type MysqlConf struct {
	Db_ip    string
	Db_port  int64
	Db       string
	Name     string
	Password string
}

type ServerConf struct {
	Mysql MysqlConf
}

func GetConfig() (ServerConf, error) {
	var config ServerConf
	if _, err := toml.DecodeFile("./conf.toml", &config); err != nil {
		return config, err
	}

	return config, nil
}
