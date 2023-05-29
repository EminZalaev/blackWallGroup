package config

import (
	"log"

	"github.com/spf13/viper"
)

// LogVars log viper vars
func LogVars(vars ...string) {
	for _, v := range vars {
		log.Printf("* %s=%s", v, viper.GetString(v))
	}
}
