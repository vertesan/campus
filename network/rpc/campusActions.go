package rpc

import (
  "fmt"
  papi "vertesan/campus/proto/papi"
  "vertesan/campus/utils"
  "vertesan/campus/utils/rich"
)

func (c *CampusClient) systemFirstCheck() {
  client := papi.NewSystemClient(c.conn)
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.SystemCheckRequest{}
  rich.Info("Calling SystemFirstCheck...")
  resp, err := client.Check(*ctx, req)
  if err != nil {
    panic(err)
  }
  utils.UNUSED(resp)
}

func (c *CampusClient) systemCheck() {
  client := papi.NewSystemClient(c.conn)
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.SystemCheckRequest{
    IdToken: c.idToken,
  }
  rich.Info("Calling SystemCheck...")
  resp, err := client.Check(*ctx, req)
  if err != nil {
    panic(err)
  }
  utils.UNUSED(resp)
}

func (c *CampusClient) authLogin() {
  client := papi.NewAuthClient(c.conn)
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.AuthLoginRequest{
    IdToken: c.idToken,
  }
  rich.Info("Calling AuthLogin...")
  resp, err := client.Login(*ctx, req)
  if err != nil {
    panic(err)
  }
  if resp.GameAuthToken != "" {
    c.authToken = resp.GameAuthToken
    c.headers["x-auth-token"] = c.authToken
  } else {
    rich.ErrorThenThrow("Get a nil authToken from server while login.")
  }
}

func (c *CampusClient) masterGet() {
  client := papi.NewMasterClient(c.conn)
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling MasterGet...")
  resp, err := client.Get(*ctx, req)
  if err != nil {
    panic(err)
  }
  if resp.MasterTag.Version != "" {
    c.masterVersion = resp.MasterTag.Version
    c.headers["x-master-version"] = c.masterVersion
    c.MasterResp = resp
  } else {
    rich.ErrorThenThrow("Get a nil MasterTag.Version from server while get masterdb.")
  }
}

func (c *CampusClient) userGet() {
  client := papi.NewUserClient(c.conn)
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling UserGet...")
  resp, err := client.Get(*ctx, req)
  if err != nil {
    panic(err)
  }
  utils.UNUSED(resp)
}

func (c *CampusClient) homeLogin() {
  client := papi.NewHomeClient(c.conn)
  c.headers["x-app-request-id"] = fmt.Sprint(getTicks())
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling HomeLogin...")
  resp, err := client.Login(*ctx, req)
  if err != nil {
    panic(err)
  }
  utils.UNUSED(resp)
}

func (c *CampusClient) loginBonusCheck() bool {
  client := papi.NewLoginBonusClient(c.conn)
  c.headers["x-app-request-id"] = fmt.Sprint(getTicks())
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling LoginBonusCheck...")
  resp, err := client.Check(*ctx, req)
  if err != nil {
    panic(err)
  }
  if len(resp.List) > 0 {
    return true
  }
  return false
}

func (c *CampusClient) loginBonusConfirm() {
  client := papi.NewLoginBonusClient(c.conn)
  c.headers["x-app-request-id"] = fmt.Sprint(getTicks())
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling LoginBonusConfirm...")
  resp, err := client.Confirm(*ctx, req)
  if err != nil {
    panic(err)
  }
  utils.UNUSED(resp)
}

func (c *CampusClient) homeEnter() {
  client := papi.NewHomeClient(c.conn)
  c.headers["x-app-request-id"] = fmt.Sprint(getTicks())
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling HomeEnter...")
  resp, err := client.Enter(*ctx, req)
  if err != nil {
    panic(err)
  }
  c.HomeEnterResp = resp
  utils.UNUSED(resp)
}

func (c *CampusClient) noticeListAll() {
  client := papi.NewNoticeClient(c.conn)
  c.headers["x-app-request-id"] = fmt.Sprint(getTicks())
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling NoticeListAll...")
  resp, err := client.ListAll(*ctx, req)
  if err != nil {
    panic(err)
  }
  c.NoticeListAllResp = resp
  utils.UNUSED(resp)
}

func (c *CampusClient) pvpRateGet() {
  client := papi.NewPvpRateClient(c.conn)
  c.headers["x-app-request-id"] = fmt.Sprint(getTicks())
  ctx, cancel := c.prepareContext()
  defer (*cancel)()
  req := &papi.Empty{}
  rich.Info("Calling PvpRateGet...")
  resp, err := client.Get(*ctx, req)
  if err != nil {
    panic(err)
  }
  c.PvpRateGetResp = resp
  utils.UNUSED(resp)
}
