package conf

import (
	"flag"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	ETCD *ETCD
}

type ETCD struct {
	Server string
	Ports  int
}

func init() {
	flag.StringVar(&confPath, "conf", "./conf.toml", "default config path")
}

// Init init conf
func Init() error {
	if confPath != "" {
		_, err := toml.DecodeFile(confPath, &Conf)
		return err
	}
	return nil
}
