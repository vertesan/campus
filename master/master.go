package master

import (
  "bytes"
  "database/sql"
  "fmt"
  "io"
  "net/http"
  "os"
  "path"
  "vertesan/campus/network/hyper/downloader"
  "vertesan/campus/proto/mapping"
  "vertesan/campus/proto/papi"
  "vertesan/campus/utils/rich"

  _ "github.com/mutecomm/go-sqlcipher/v4"
  // "google.golang.org/protobuf/encoding/protojson"
  "google.golang.org/protobuf/proto"
)

const MASTER_RAW_PATH = "cache/masterRaw"
const MASTER_JSON_PATH = "cache/masterJson"
const MASTER_YAML_PATH = "cache/masterYaml"

func DownloadAndDecrypt(masterTagResp *papi.MasterGetResponse) {
  DownloadAllMaster(masterTagResp)
  DecryptAll(masterTagResp)
}

func UnmarshalPlain(reader io.Reader) (*papi.MasterGetResponse, error) {
  masterGetResp := &papi.MasterGetResponse{}
  buf := &bytes.Buffer{}
  if _, err := io.Copy(buf, reader); err != nil {
    panic(err)
  }
  proto.Unmarshal(buf.Bytes(), masterGetResp)
  return masterGetResp, nil
}

func DownloadAllMaster(masterTagResp *papi.MasterGetResponse) {
  masterHeader := &http.Header{
    "User-Agent":      {"UnityPlayer/2022.3.21f1 (UnityWebRequest/1.0, libcurl/8.5.0-DEV)"},
    "Accept":          {"*/*"},
    "X-Unity-Version": {"2022.3.21f1"},
  }
  dler := downloader.NewDownloader(30, masterHeader, MASTER_RAW_PATH, 5)
  entries := []*downloader.Entry{}
  for _, masterPack := range masterTagResp.MasterTag.MasterTagPacks {
    entries = append(entries, &downloader.Entry{
      Url:          masterPack.DownloadUrl,
      SaveFileName: masterPack.Type,
    })
  }
  dler.SetEntries(entries)
  if err := dler.DownloadAll(); err != nil {
    panic(err)
  }
}

func DecryptAll(masterTagResp *papi.MasterGetResponse) {
  rich.Info("Start to decrypt database.")
  for _, masterTagPack := range masterTagResp.MasterTag.MasterTagPacks {
    dbPath := path.Join(MASTER_RAW_PATH, masterTagPack.Type)
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
    // jsonList := []string{}
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
      yamlBytes, err := YamlMarshal(instance)
      if err != nil {
        panic(err)
      }

      // jsonBytes, err := marshalOptions.Marshal(instance)
      // if err != nil {
      //   panic(err)
      // }
      // jsonList = append(jsonList, string(jsonBytes))

      yamlBytes = bytes.TrimSuffix(yamlBytes, []byte("\n"))
      yamlBytes = bytes.Replace(yamlBytes, []byte("\n"), []byte("\n  "), -1)
      yamlList = append(yamlList, yamlBytes)
    }
    var yamlDb []byte
    if len(yamlList) != 0 {
      yamlDb = append([]byte("- "), bytes.Join(yamlList, []byte("\n- "))...)
    }
    // jsonDb := "[\n" + strings.Join(jsonList, ",") + "]"
    writeYaml(masterTagPack.Type, yamlDb)
    // writeJson(masterTagPack.Type, &jsonDb)
    rich.Info("Database %q is successfully decrypted.", masterTagPack.Type)
  }
  rich.Info("Database decrypting completed.")
}

func writeYaml(name string, data []byte) {
  if err := os.MkdirAll(MASTER_YAML_PATH, 0755); err != nil {
    panic(err)
  }
  filePath := path.Join(MASTER_YAML_PATH, name+".yaml")
  if err := os.WriteFile(filePath, data, 0644); err != nil {
    panic(err)
  }
}

func writeJson(name string, data *string) {
  if err := os.MkdirAll(MASTER_JSON_PATH, 0755); err != nil {
    panic(err)
  }
  filePath := path.Join(MASTER_JSON_PATH, "name"+".json")
  if err := os.WriteFile(filePath, []byte(*data), 0644); err != nil {
    panic(err)
  }
}
