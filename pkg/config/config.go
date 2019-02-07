package config

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
	"log"
)


func Init(c *pflag.FlagSet) {
	fs := afero.NewOsFs()
	viper.SetFs(fs)
	viper.AutomaticEnv()
	viper.SetConfigName("pb")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	ifErrr(viper.BindPFlags(c))
	fmt.Println(viper.SafeWriteConfig())
	ifErrr(viper.ReadInConfig())
}

func AddConfigPaths(path ...string) func() {
	return func() {
		for _, p := range path {
			viper.AddConfigPath(p)
		}
	}
}

func ifErrr (err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}