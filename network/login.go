package network

import (
  "vertesan/campus/config"
  "vertesan/campus/network/hyper"
  "vertesan/campus/network/jwto"
  "vertesan/campus/network/rpc"
  "vertesan/campus/utils/rich"
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
