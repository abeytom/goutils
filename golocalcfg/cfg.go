package golocalcfg

import (
	"fmt"
	"github.com/abeytom/goutils/gofile"
	"os"
	"path/filepath"
	yaml "sigs.k8s.io/yaml/goyaml.v2"
)

type LocalConfig map[string]interface{}

func LoadDefaultCfg() (LocalConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	return LoadCfg(homeDir)
}

func LoadCfg(basedir string) (LocalConfig, error) {
	configPath := filepath.Join(basedir, ".config/.local-config.yml")
	if !gofile.IsFile(configPath) {
		return map[string]interface{}{}, nil
	}

	// Open the YAML file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Decode the YAML into a map
	var config map[string]interface{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode YAML file: %w", err)
	}
	return config, nil
}

func (lc LocalConfig) GetStringOrEmpty(key string) string {
	val, exists := lc[key]
	if !exists {
		return ""
	}
	s, ok := val.(string)
	if ok {
		return s
	} else {
		return ""
	}
}
