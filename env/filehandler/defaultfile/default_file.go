package defaultfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zouyx/agollo/v3/env"
	"github.com/zouyx/agollo/v3/env/filehandler"
	"os"

	"github.com/zouyx/agollo/v3/component/log"
	jsonConfig "github.com/zouyx/agollo/v3/env/config/json"
)

const suffix = ".json"

func init() {
	filehandler.SetFileHandler(&DefaultFile{})
}

//DefaultFile 默认备份文件读写
type DefaultFile struct {
}

var (
	configFileMap  = make(map[string]string, 1)
	jsonFileConfig = &jsonConfig.ConfigFile{}
)

//WriteWithRaw decorator for WriteConfigFile
func (fileHandler *DefaultFile) WriteWithRaw(f func(config *env.ApolloConfig, configPath string) error) func(config *env.ApolloConfig, configPath string) error {
	return func(config *env.ApolloConfig, configPath string) error {
		filePath := fmt.Sprintf("%s/%s", configPath, config.NamespaceName)
		file, e := os.Create(filePath)
		if e != nil {
			return e
		}
		defer file.Close()
		_, e = file.WriteString(config.Configurations["content"])
		if e != nil {
			return e
		}

		return f(config, configPath)
	}
}

//WriteConfigFile write config to file
func (fileHandler *DefaultFile) WriteConfigFile(config *env.ApolloConfig, configPath string) error {
	return jsonFileConfig.Write(config, fileHandler.GetConfigFile(configPath, config.NamespaceName))
}

//GetConfigFile get real config file
func (fileHandler *DefaultFile) GetConfigFile(configDir string, namespace string) string {
	fullPath := configFileMap[namespace]
	if fullPath == "" {
		filePath := fmt.Sprintf("%s%s", namespace, suffix)
		if configDir != "" {
			configFileMap[namespace] = fmt.Sprintf("%s/%s", configDir, filePath)
		} else {
			configFileMap[namespace] = filePath
		}
	}
	return configFileMap[namespace]
}

//LoadConfigFile load config from file
func (fileHandler *DefaultFile) LoadConfigFile(configDir string, namespace string) (*env.ApolloConfig, error) {
	configFilePath := fileHandler.GetConfigFile(configDir, namespace)
	log.Info("load config file from :", configFilePath)
	c, e := jsonFileConfig.Load(configFilePath, func(b []byte) (interface{}, error) {
		config := &env.ApolloConfig{}
		e := json.NewDecoder(bytes.NewBuffer(b)).Decode(config)
		return config, e
	})

	if c == nil || e != nil {
		log.Errorf("loadConfigFile fail,error:", e)
		return nil, e
	}

	return c.(*env.ApolloConfig), e
}
