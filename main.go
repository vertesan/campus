package main

import (
	"flag"
	"os"
	"vertesan/campus/analyser"
	"vertesan/campus/config"
	"vertesan/campus/master"
	"vertesan/campus/network"
	"vertesan/campus/octo"
	"vertesan/campus/utils/rich"
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

func processAnalysis() {
  analyser.Analyze()
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
    network.ProcessMasterDbAndApiResp(*flagForceDb, *flagPutDb)
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
