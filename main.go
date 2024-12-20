package main

import (
  "flag"
  "fmt"
  "os"
  "vertesan/campus/analyser"
  "vertesan/campus/config"
  "vertesan/campus/master"
  "vertesan/campus/network"
  "vertesan/campus/octo"
  "vertesan/campus/utils/rich"

  "google.golang.org/protobuf/encoding/protojson"
)

const NEW_AB_FLAG_FILE = "cache/newab_flag"
const OCTO_CACHE_FILE = "cache/octocacheevai"
const MASTER_TAG_FILE = "cache/masterGetDec240522202758690.bin"

var (
  flagDb        = flag.Bool("db", false, "Download and decrypt master database if true.\nGenerated yaml files are saved in 'cache/masterYaml' directory.")
  flagForceDb   = flag.Bool("forcedb", false, "Download and decrypt master database without checking local version.\nTake no effect if 'db' flag is absent.")
  flagKeepRaw   = flag.Bool("keepdb", false, "Do not delete encrypted master database files after decrypting.\nTake no effect if 'db' flag is absent.")
  flagPutDb     = flag.Bool("putdb", false, "Put required DBs to CAMPUS_DB_PUT_URL.")
  flagAb        = flag.Bool("ab", false, "Download and deobfuscate assetbundles if true.\nDeobfuscated files are saved in 'cache/assets' directory.")
  flagForceAb   = flag.Bool("forceab", false, "Download assetbundles without checking version.\nTakes no effect if 'ab' is absent.\nIt's safe to set this flag to true if you only want to download a part of additional assets instead of the entire bulky thing because MD5 check will still be carried out before downloading.")
  flagKeepAbRaw = flag.Bool("keepab", false, "Do not delete obfuscated assetbundle files after deobfuscating.\nTake no effect if 'ab' flag is absent.")
  flagWebAb     = flag.Bool("webab", false, "Only download images those are needed for web use. Takes no effect if '--ab' is absent.")
  flagAnalyze   = flag.Bool("analyze", false, "Analyze dump.cs to retrieve proto schema.\nGenerated codes are saved in 'cache/GeneratedProto' directory.")
  refToken      = flag.String("token", "", "The refresh token used to retrieve login idToken from firebase.\nIf refreshToken field set in 'config.yaml' is not empty, the value in the config file will take precedence.")
)

func decryptOctoManifest() {
  // open encrypted octo cache file
  octoFs, err := os.Open(OCTO_CACHE_FILE)
  if err != nil {
    panic(err)
  }
  // keyMd5 := md5.Sum([]byte("1nuv9td1bw1udefk"))
  // ivMd5 := md5.Sum([]byte("LvAUtf+tnz"))
  octoDb, err := octo.DecryptOctoList(octoFs, 1)
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

func processMasterDbAndApiResp(flagForceDb bool) {
  // simulate login to get master database manifest
  manager := &network.NetworkManager{}
  manager.Login()
  // put response db
  if *flagPutDb {
    processApiResponse(manager)
  }
  // compare local db version with server
  cfg := config.GetConfig()
  rich.Info("Current local database version: %q.", cfg.MasterVersion)
  serverVer := manager.Client.MasterResp.MasterTag.Version
  rich.Info("Server database version: %q.", serverVer)
  if flagForceDb || cfg.MasterVersion != serverVer {
    rich.Info("New database version detected: %q.", serverVer)
    // save MasterResponse to JSON
    jsonMarshalOptions := protojson.MarshalOptions{
      Multiline:         true,
      AllowPartial:      false,
      UseProtoNames:     true,
      UseEnumNumbers:    true,
      EmitUnpopulated:   true,
      EmitDefaultValues: false,
    }
    jsonBytes, err := jsonMarshalOptions.Marshal(manager.Client.MasterResp)
    if err != nil {
      panic(err)
    }
    if err := os.WriteFile("cache/MasterResponse.json", jsonBytes, 0644); err != nil {
      panic(err)
    }
    // download master database
    master.DownloadAndDecrypt(manager.Client.MasterResp, *flagPutDb)
    if *flagPutDb {
      master.PutDb("Version", fmt.Sprintf(`{"version":"%s"}`, serverVer))
    }
  } else {
    rich.Info("Local database is already up to date, skip downloading database.")
  }
  cfg.MasterVersion = serverVer
  cfg.Save()
  os.WriteFile("cache/master_version", []byte(serverVer), 0644)
}

// APIs must be invoked every time
func processApiResponse(manager *network.NetworkManager) {
  jsonMarshalOptions := protojson.MarshalOptions{
    Multiline:     false,
    AllowPartial:  false,
    UseProtoNames: true,
    // use enum numbers in json, this option is different with yaml
    UseEnumNumbers:    true,
    EmitUnpopulated:   true,
    EmitDefaultValues: false,
  }
  // Home.Enter
  homeBytes, err := jsonMarshalOptions.Marshal(manager.Client.HomeEnterResp)
  if err != nil {
    panic(err)
  }
  master.PutDb("HomeEnter", string(homeBytes))
  // Hotice.ListAll
  noticeBytes, err := jsonMarshalOptions.Marshal(manager.Client.NoticeListAllResp)
  if err != nil {
    panic(err)
  }
  master.PutDb("NoticeListAll", string(noticeBytes))
  // PvpRate.Get
  if manager.Client.PvpRateGetResp != nil {
    pvpBytes, err := jsonMarshalOptions.Marshal(manager.Client.PvpRateGetResp)
    if err != nil {
      panic(err)
    }
    master.PutDb("PvpRateGet", string(pvpBytes))
  }
}

func main() {
  rich.Info("Start process.")

  if err := os.Remove(NEW_AB_FLAG_FILE); err != nil {
    if !os.IsNotExist(err) {
      panic(err)
    }
  }

  flag.Parse()
  cfg := config.GetConfig()

  if *refToken != "" {
    if cfg.RefreshToken == "" {
      cfg.RefreshToken = *refToken
      cfg.Save()
    }
  }

  if *flagDb {
    processMasterDbAndApiResp(*flagForceDb)
    if !*flagKeepRaw {
      os.RemoveAll(master.MASTER_RAW_PATH)
    }
  }

  if *flagAb {
    manager := &octo.OctoManager{}
    hasUpdates := manager.Work(*flagKeepAbRaw, *flagWebAb, *flagForceAb)
    cfg.OctoCacheRevision = int(manager.OctoDb.Revision)
    cfg.Save()
    if hasUpdates {
      if _, err := os.Create(NEW_AB_FLAG_FILE); err != nil {
        panic(err)
      }
    }
  }

  if *flagAnalyze {
    processAnalysis()
  }

  rich.Info("All process completed.")
}
