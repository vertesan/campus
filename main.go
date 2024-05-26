package main

import (
  "crypto/md5"
  "flag"
  "os"
  "vertesan/campus/analyser"
  "vertesan/campus/config"
  "vertesan/campus/master"
  "vertesan/campus/network"
  "vertesan/campus/octo"
  "vertesan/campus/utils/rich"

  "google.golang.org/protobuf/encoding/protojson"
)

const OCTO_CACHE_FILE = "cache/octocacheevai"
const MASTER_TAG_FILE = "cache/masterGetDec240522202758690.bin"

var (
  flagDb      = flag.Bool("db", false, "Download and decrypt master database if true. Generated yaml files are saved in 'cache/masterYaml' directory.")
  flagKeepRaw = flag.Bool("keep", false, "Do not delete encrypted master database files after decrypting. Take no effect if 'db' flag is not set.")
  flagAnalyze = flag.Bool("analyze", false, "Analyze dump.cs to retrieve proto schemas. Generated codes are saved in 'cache/GeneratedProto' directory.")
  refToken    = flag.String("token", "", "The refresh token used to retrieve login idToken from firebase.\nIf refreshToken field set in 'config/config.yaml' is not empty, the value in the config file will take precedence.")
)

func decryptOctoManifest() {
  // open encrypted octo cache file
  octoFs, err := os.Open(OCTO_CACHE_FILE)
  if err != nil {
    panic(err)
  }
  keyMd5 := md5.Sum([]byte("1nuv9td1bw1udefk"))
  ivMd5 := md5.Sum([]byte("LvAUtf+tnz"))
  octoDb, err := octo.DecryptManifest(octoFs, keyMd5[:], ivMd5[:], 1)
  if err != nil {
    panic(err)
  }
  // convert protobuf object to json string
  marshalOptions := &protojson.MarshalOptions{
    Multiline:      true,
    Indent:         "  ",
    AllowPartial:   true,
    UseProtoNames:  true,
    UseEnumNumbers: true,
  }
  octoJson, err := marshalOptions.Marshal(octoDb)
  if err != nil {
    panic(err)
  }
  // write manifest json file
  if err := os.WriteFile("cache/octo.json", octoJson, 0644); err != nil {
    panic(err)
  }
}

func processAnalysis() {
  analyser.Analyze()
}

func processMasterDb() {
  // simulate login to get master database manifest
  manager := &network.NetworkManager{}
  manager.Login()

  // compare local db version with server
  cfg := config.GetConfig()
  rich.Info("Current local database version: %q.", cfg.MasterVersion)
  serverVer := manager.Client.MasterResp.MasterTag.Version
  rich.Info("Server database version: %q.", serverVer)
  if cfg.MasterVersion != serverVer {
    rich.Info("New database version detected: %q.", serverVer)
    // download master database
    master.DownloadAndDecrypt(manager.Client.MasterResp)
  } else {
    rich.Info("Local database is already up to date, skip downloading database.")
  }
  cfg.MasterVersion = serverVer
  cfg.Save()
  os.WriteFile("cache/master_version", []byte(serverVer), 0644)
}

func main() {
  rich.Info("Start process.")
  flag.Parse()
  cfg := config.GetConfig()

  if *refToken != "" {
    if cfg.RefreshToken == "" {
      cfg.RefreshToken = *refToken
      cfg.Save()
    }
  }

  if *flagDb {
    processMasterDb()
    if !*flagKeepRaw {
      os.RemoveAll(master.MASTER_RAW_PATH)
    }
  }

  if *flagAnalyze {
    processAnalysis()
  }

  rich.Info("All process completed.")
}
