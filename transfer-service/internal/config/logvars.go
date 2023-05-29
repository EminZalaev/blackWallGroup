package config

import (
	"github.com/spf13/viper"
	"log"
)

// LogVars log viper vars
func LogVars(vars ...string) {
	for _, v := range vars {
		log.Printf("* %s=%s", v, viper.GetString(v))
	}
}
