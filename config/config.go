package config

import (
  "os"

  "github.com/goccy/go-yaml"
)

const CONFIG_FILE_PATH = "config.yaml"

type Config struct {
  AppVersion           string `yaml:"appVersion"`
  MasterVersion        string `yaml:"masterVersion"`
  OctoCacheRevision    int    `yaml:"octoCacheRevision"`
  OctoManifestRevision int    `yaml:"octoManifestRevision"`
  RefreshToken         string `yaml:"refreshToken"`
  IdToken              string `yaml:"idToken"`
}

func (c *Config) Save() {
  yamlBytes, err := yaml.Marshal(c)
  if err != nil {
    panic(err)
  }
  if err := os.WriteFile(CONFIG_FILE_PATH, yamlBytes, 0644); err != nil {
    panic(err)
  }
}

var _config *Config

func createConfigFile() {
  newFileBytes, err := yaml.Marshal(_config)
  if err != nil {
    panic(err)
  }
  if err := os.WriteFile(CONFIG_FILE_PATH, newFileBytes, 0644); err != nil {
    panic(err)
  }
}

func init() {
  _config = &Config{}
  yamlBytes, err := os.ReadFile(CONFIG_FILE_PATH)
  if err != nil {
    if os.IsNotExist(err) {
      createConfigFile()
    } else {
      panic(err)
    }
  }
  if err := yaml.Unmarshal(yamlBytes, _config); err != nil {
    panic(err)
  }
}

func GetConfig() *Config {
  return _config
}
