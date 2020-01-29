package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load the configuration
func Load() {
	// Setting default value
	viper.SetDefault("APP_URL", "0.0.0.0")
	viper.SetDefault("APP_PORT", "8080")

	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("error reading config : %s ", err))
	}
	viper.AutomaticEnv()

}

// ReadInt reads a int from configuration based on its key
func ReadInt(key string) int {
	isSet(key)
	return viper.GetInt(key)
}

// ReadString reads a string from configuration based on its key
func ReadString(key string) string {
	isSet(key)
	return viper.GetString(key)
}

// ReadBool reads a bool from configuration based on its key
func ReadBool(key string) bool {
	isSet(key)
	return viper.GetBool(key)
}

func isSet(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Errorf("key %s is not set", key))
	}
}
