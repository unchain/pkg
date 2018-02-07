package xconfig

import (
	"github.com/spf13/viper"
	"github.com/pkg/errors"
)

func LoadConfig(cfg interface{}, path string) error {
	var err error 
	
	v := viper.New()
	v.SetConfigFile(path)
	err = v.ReadInConfig()
	
	if err != nil {
		return errors.Wrapf(err, "failed to read in config file %s", path)
	}
	
	err = v.Unmarshal(cfg)
	
	if err != nil {
		return errors.Wrapf(err, "failed to parse config file %s", path)
	}
	
	return nil
}

