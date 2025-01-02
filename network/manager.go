package network

import (
	"fmt"
	"os"
	"vertesan/campus/config"
	"vertesan/campus/master"
	"vertesan/campus/network/hyper"
	"vertesan/campus/network/jwto"
	"vertesan/campus/network/rpc"
	"vertesan/campus/utils/rich"

	"google.golang.org/protobuf/encoding/protojson"
)

type NetworkManager struct {
  Client *rpc.CampusClient
}

func (m *NetworkManager) Login() {
  config := config.GetConfig()

  // retrieve app version
  gameVersion, err := hyper.GetPlayVersion()
  if err != nil {
    rich.Error("%v", err)
    rich.Error("Failed to retrieve game app version from GooglePlay.")
    // fallback
    gameVersion = config.AppVersion
  } else {
    config.AppVersion = gameVersion
    config.Save()
  }

  // retrieve firebase token
  if config.IdToken == "" || jwto.IsJwtExpired(config.IdToken) {
    idToken, refToken := hyper.GetFirebaseToken(config.RefreshToken)
    config.IdToken = idToken
    config.RefreshToken = refToken
    config.Save()
  }

  // simulate login
  m.Client = rpc.NewCampusClient(config.IdToken, gameVersion)
  m.Client.DoLogin()
}

func ProcessMasterDbAndApiResp(flagForceDb bool, flagPutDb bool) {
  // simulate login to get master database manifest
  manager := &NetworkManager{}
  manager.Login()
  // put response db
  if flagPutDb {
    ProcessApiResponse(manager)
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
    master.DownloadAndDecrypt(manager.Client.MasterResp, flagPutDb)
    if flagPutDb {
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
func ProcessApiResponse(manager *NetworkManager) {
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
