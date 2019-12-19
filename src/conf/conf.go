package conf

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type config struct {
	Db struct {
		Dbms    string
		Name    string
		User    string
		Pass    string
		Host    string
		Port    string
		ConnOpt string `mapstructure:"conn_opt"`
	}
	Sv struct {
		Timeout int64
		Port    string
		Debug   bool
	}
	Auth struct {
		Jwtkey string
	}
}

var C config

func Readconf() {

	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath("conf")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spew.Dump(C)
}
