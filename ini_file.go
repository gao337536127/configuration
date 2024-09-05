package configuration

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"gopkg.in/ini.v1"
)

var defaultConfigurationBytes = make([][]byte, 0)

func AppendDefaultConfigurationBytes(configBytes []byte) {
	defaultConfigurationBytes = append(defaultConfigurationBytes, configBytes)
}

var configFilePaths = make([]string, 0)

func AppendConfigFile(path string) {
	configFilePaths = append(configFilePaths, path)
}

type IniConfig struct {
	once sync.Once
	cfg  *ini.File
}

func (i *IniConfig) GetConfig(section, key, defaultValue string) (string, error) {
	i.once.Do(func() {
		i.cfg = ini.Empty()
		for _, bs := range defaultConfigurationBytes {
			err := i.cfg.Append(bs)
			if err != nil {
				return
			}
		}
		for _, path := range configFilePaths {
			err := i.cfg.Append(path)
			if err != nil {
				return
			}
		}
	})

	if i.cfg == nil {
		return defaultValue, errors.New("read config error")
	}

	sec := i.cfg.Section(section)
	if sec == nil {
		return defaultValue, fmt.Errorf("the section named '%s' could not be located", section)
	}

	k := sec.Key(key)
	if k == nil {
		return defaultValue, fmt.Errorf("the key named '%s' could not be located", key)
	}

	return k.Validate(func(s string) string {
		if len(s) == 0 {
			return defaultValue
		}
		return s
	}), nil
}

func (i *IniConfig) GetConfigWithEnvironment(environment, section, key, defaultValue string) (string, error) {
	env := os.Getenv(environment)
	if env != "" {
		return env, nil
	} else {
		return i.GetConfig(section, key, defaultValue)
	}
}
