package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
)

type System struct {
	Connections Connections `yaml:"connections"`
}

type Connections struct {
	Name Name
}

type Name struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func LoadFromFile(path string) (*System, error) {

	s := &System{}
	v := NewEnviper(viper.New())
	v.SetDefault("admin.enabled", true)

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/")
	v.AddConfigPath("/etc/app_name/")
	if err := v.Unmarshal(s, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		))); err != nil {
		return nil, fmt.Errorf("cannot read unmarshal: %w", err)
	}
	return s, nil
}
