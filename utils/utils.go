package utils

import (
  "bufio"
  "os"
  "vertesan/campus/utils/rich"

  "github.com/goccy/go-json"
  "github.com/goccy/go-yaml"
)

func UNUSED(x ...interface{}) {}

func WriteToJsonFile(instance any, dst string) {
  jsonBytes, err := json.MarshalIndent(instance, "", "  ")
  if err != nil {
    panic(err)
  }
  jsonDbFile, err := os.Create(dst)
  if err != nil {
    panic(err)
  }
  jFileBuf := bufio.NewWriter(jsonDbFile)
  jFileBuf.Write(jsonBytes)
  jFileBuf.Flush()
  rich.Info("Completed writing json file '%s'.", jsonDbFile.Name())
  jsonDbFile.Close()
}

func ReadFromJsonFile(src string, v *any) error {
  jBytes, err := os.ReadFile(src)
  if err != nil {
    return err
  }
  if err := json.Unmarshal(jBytes, v); err != nil {
    return err
  }
  return nil
}

func WriteToYamlFile(instance any, dst string) {
  yamlBytes, err := yaml.Marshal(instance)
  if err != nil {
    panic(err)
  }
  yamlDbFile, err := os.Create(dst)
  if err != nil {
    panic(err)
  }
  yFileBuf := bufio.NewWriter(yamlDbFile)
  yFileBuf.Write(yamlBytes)
  yFileBuf.Flush()
  rich.Info("Completed writing yaml file '%s'.", yamlDbFile.Name())
  yamlDbFile.Close()
}
