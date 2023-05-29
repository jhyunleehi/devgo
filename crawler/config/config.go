package crawler

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfigFile(envName string) {

	viper.SetConfigName(envName) // name of config file (without extension)
	viper.SetConfigType("yaml")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("conf/") // path to look for the config file in
	err := viper.ReadInConfig()  // Find and read the config file

	if err != nil { // Handle errors reading the config file
		log.Printf(err.Error())
		panic("A fatal error occurred while reading the configuration file.")
	}
}

// SetDefault Func is...
func SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

// Get Func is get the value corresponding to the key.
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetInt Func is get the int value corresponding to the key.
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool Func is get the boolean value corresponding to the key.
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetStringSlice Func is get the string-slice value corresponding to the key.
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetStringMap Func is get the string-map value corresponding to the key.
func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// GetStringMapString Func is get the string-map-string value corresponding to the key.
func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// GetStringMapStringSlice Func is get the string-map-string-slice value corresponding to the key.
func GetStringMapStringSlice(key string) map[string][]string {
	return viper.GetStringMapStringSlice(key)
}
