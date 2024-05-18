package master

import (
  "bytes"
  "database/sql"
  "fmt"
  "io"
  "os"
  "vertesan/campus/proto/mapping"
  "vertesan/campus/proto/mastertag"

  "github.com/bufbuild/protoyaml-go"
  _ "github.com/mutecomm/go-sqlcipher/v4"
  "google.golang.org/protobuf/proto"
)

var (
  masterRawPath = "cache/masterRaw/"
  // masterJsonPath = "cache/masterJson/"
  masterYamlPath = "cache/masterYaml/"
  // marshalOptions = &protojson.MarshalOptions{
  //   Multiline:         false,
  //   Indent:            "",
  //   AllowPartial:      false,
  //   UseProtoNames:     true,
  //   UseEnumNumbers:    true,
  //   EmitUnpopulated:   false,
  //   EmitDefaultValues: false,
  // }
  yamlMarshalOptions = &protoyaml.MarshalOptions{
    Indent:          2,
    AllowPartial:    false,
    UseProtoNames:   true,
    UseEnumNumbers:  false,
    EmitUnpopulated: false,
  }
)

func UnmarshalPlain(reader io.Reader) (*mastertag.MasterGetResponse, error) {
  masterGetResp := &mastertag.MasterGetResponse{}
  buf := &bytes.Buffer{}
  if _, err := io.Copy(buf, reader); err != nil {
    panic(err)
  }
  proto.Unmarshal(buf.Bytes(), masterGetResp)
  return masterGetResp, nil
}

func DownloadAllMaster(masterTagResp *mastertag.MasterGetResponse) {

}

func DecryptAll(masterTagResp *mastertag.MasterGetResponse) {
  for _, masterTagPack := range masterTagResp.MasterTag.MasterTagPacks {
    dbPath := masterRawPath + masterTagPack.Type
    key := masterTagPack.CryptoKey
    dbname := fmt.Sprintf("%s?_pragma_key=x'%s'", dbPath, key)
    db, err := sql.Open("sqlite3", dbname)
    if err != nil {
      panic(err)
    }
    defer db.Close()
    rows, err := db.Query(fmt.Sprintf("select data from %s;", masterTagPack.Type))
    if err != nil {
      panic(err)
    }
    defer rows.Close()
    yamlList := [][]byte{}
    for rows.Next() {
      var data []byte
      if err = rows.Scan(&data); err != nil {
        panic(err)
      }
      instance := mapping.ProtoMap["Master."+masterTagPack.Type]
      if err := proto.Unmarshal(data, instance); err != nil {
        panic(err)
      }
      yamlBytes, err := yamlMarshalOptions.Marshal(instance)
      if err != nil {
        panic(err)
      }
      yamlBytes = bytes.TrimSuffix(yamlBytes, []byte("\n"))
      yamlBytes = bytes.Replace(yamlBytes, []byte("\n"), []byte("\n  "), -1)
      yamlList = append(yamlList, yamlBytes)
    }
    var yamlDb []byte
    if len(yamlList) != 0 {
      yamlDb = append([]byte("- "), bytes.Join(yamlList, []byte("\n- "))...)
    }
    writeYaml(masterTagPack.Type, yamlDb)
  }
}

func writeYaml(name string, data []byte) {
  if err := os.MkdirAll(masterYamlPath, 0755); err != nil {
    panic(err)
  }
  filePath := masterYamlPath + name + ".yaml"
  if err := os.WriteFile(filePath, data, 0644); err != nil {
    panic(err)
  }
}

// func writeJson(name string, data *string) {
//   filePath := masterJsonPath + name + ".json"
//   if err := os.WriteFile(filePath, []byte(*data), 0644); err != nil {
//     panic(err)
//   }
// }
