package master

import (
  "bytes"
  "database/sql"
  "fmt"
  "io"
  "net/http"
  "net/url"
  "os"
  "path"
  "regexp"
  "slices"
  "strings"
  "vertesan/campus/network/hyper"
  "vertesan/campus/network/hyper/downloader"
  "vertesan/campus/proto/mapping"
  "vertesan/campus/proto/papi"
  "vertesan/campus/utils/rich"

  _ "github.com/mutecomm/go-sqlcipher/v4"
  "google.golang.org/protobuf/encoding/protojson"
  "google.golang.org/protobuf/proto"
)

const MASTER_RAW_PATH = "cache/masterRaw"
const MASTER_JSON_PATH = "cache/masterJson"
const MASTER_YAML_PATH = "cache/masterYaml"
const ENV_CAMPUS_DB_PUT_URL = "CAMPUS_DB_PUT_URL"
const ENV_CAMPUS_DB_PUT_SECRET = "CAMPUS_DB_PUT_SECRET"

var requiredPutTypes = []string{
  "Character",
  "ProduceDescription",
  "ProduceEffectIcon",
  "Produce",
  "ProduceGroup",
  "ExamInitialDeck",
  "ProduceDescriptionProduceEffectType",
  "ProduceDescriptionProduceExamEffectType",
  "SupportCard",
  "ProduceCard",
  "ProduceItem",
  "ProduceEventSupportCard",
  "ProduceStepEventDetail",
  "ProduceEffect",
  "SupportCardProduceSkillLevelDance",
  "SupportCardProduceSkillLevelVocal",
  "SupportCardProduceSkillLevelVisual",
  "SupportCardProduceSkillLevelAssist",
  "ProduceSkill",
  "ProduceTrigger",
  "IdolCard",
  "IdolCardSkin",
  "IdolCardPotential",
  "IdolCardPotentialProduceSkill",
  "IdolCardLevelLimit",
  "IdolCardLevelLimitProduceSkill",
  "IdolCardLevelLimitStatusUp",
  "ProduceStepAuditionDifficulty",
  "ProduceExamBattleNpcGroup",
  "ProduceExamBattleConfig",
  "ProduceExamBattleScoreConfig",
  "ProduceExamGimmickEffectGroup",
  "CharacterTrueEndBonus",
  "PvpRateConfig",
  "PvpRateCommonProduceCard",
  "ExamSetting",
  "StoryEvent",
  "CharacterDetail",
  "Achievement",
  "AchievementProgress",
  "ProduceExamEffect",
  "EventLabel",
  "MemoryAbility",
  "ResultGradePattern",
  "GuildReaction",
}

func DownloadAndDecrypt(masterTagResp *papi.MasterGetResponse, putDb bool) {
  DownloadAllMaster(masterTagResp)
  DecryptAll(masterTagResp, putDb)
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

func PutDb(name string, jsonDb string) {
  remoteUrl := os.Getenv(ENV_CAMPUS_DB_PUT_URL)
  secret := os.Getenv(ENV_CAMPUS_DB_PUT_SECRET)
  if remoteUrl == "" || secret == "" {
    rich.ErrorThenThrow("Environment variable CAMPUS_DB_PUT_URL or ENV_CAMPUS_DB_PUT_SECRET does not exist.")
  }
  url, err := url.JoinPath(remoteUrl, name)
  if err != nil {
    panic(err)
  }
  headers := &http.Header{
    "Content-Type":  {"multipart/form-data; boundary=---011000010111000001101001"},
    "Authorization": {fmt.Sprintf("Bearer %s", secret)},
  }
  payload := strings.NewReader(fmt.Sprintf("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"metadata\"\r\n\r\n{}\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"value\"\r\n\r\n%s\r\n-----011000010111000001101001--\r\n\r\n", jsonDb))
  hyper.SendRequest(url, "PUT", headers, payload, 10, 3)
  rich.Info("Database %q is successfully put to remote db.", name)
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

func DecryptAll(masterTagResp *papi.MasterGetResponse, putDb bool) {
  rich.Info("Start to decrypt database.")
  jsonMarshalOptions := protojson.MarshalOptions{
    Multiline:     false,
    AllowPartial:  false,
    UseProtoNames: true,
    // use enum numbers in json, this option is different with yaml
    UseEnumNumbers:    true,
    EmitUnpopulated:   true,
    EmitDefaultValues: false,
  }
  for _, masterTagPack := range masterTagResp.MasterTag.MasterTagPacks {
    if _, ok := mapping.ProtoMap["Master."+masterTagPack.Type]; !ok {
      rich.Warning("'Master.%s' not found in existing map, perhaps a new database entry is introduced.", masterTagPack.Type)
      rich.Warning("Skip processing 'Master.%s'", masterTagPack.Type)
      continue
    }
    // read db
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
    jsonList := []string{}
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

      jsonBytes, err := jsonMarshalOptions.Marshal(instance)
      if err != nil {
        panic(err)
      }
      jsonList = append(jsonList, string(jsonBytes))

      yamlBytes = bytes.TrimSuffix(yamlBytes, []byte("\n"))
      reg := regexp.MustCompile(`\n(?P<ctt>.+)`)
      yamlBytes = reg.ReplaceAll(yamlBytes, []byte("\n  $ctt"))
      yamlList = append(yamlList, yamlBytes)
    }
    var yamlDb []byte
    if len(yamlList) != 0 {
      yamlDb = append([]byte("- "), bytes.Join(yamlList, []byte("\n- "))...)
    }
    jsonDb := "[" + strings.Join(jsonList, ",") + "]"
    writeYaml(masterTagPack.Type, yamlDb)

    // determine whether this Type needs to be put
    if putDb && slices.Contains(requiredPutTypes, masterTagPack.Type) {
      PutDb(masterTagPack.Type, jsonDb)
    }

    writeJson(masterTagPack.Type, &jsonDb)
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
  filePath := path.Join(MASTER_JSON_PATH, name+".json")
  if err := os.WriteFile(filePath, []byte(*data), 0644); err != nil {
    panic(err)
  }
}
