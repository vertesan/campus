package rpc

import (
  "context"
  "time"
  "vertesan/campus/config"
  "vertesan/campus/proto/papi"
  "vertesan/campus/proto/penum"

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/encoding"
  "google.golang.org/grpc/metadata"
)

const SERVER_URL = "api.game-gakuen-idolmaster.jp"
const TIMEOUT_SEC = 10

func init() {
  // register QuaCodec
  encoding.RegisterCodec(&QuaCodec{})
}

type CampusClient struct {
  idToken           string // issued by firebase
  authToken         string // issued by QA
  masterVersion     string
  appVersion        string
  serverUrl         string
  timeout           int
  headers           map[string]string
  conn              *grpc.ClientConn
  MasterResp        *papi.MasterGetResponse
  HomeEnterResp     *papi.HomeEnterResponse
  NoticeListAllResp *papi.NoticeListAllResponse
  PvpRateGetResp    *papi.PvpRateGetResponse
}

func NewCampusClient(idToken string, appVersion string) *CampusClient {
  if idToken == "" {
    idToken = config.GetConfig().IdToken
  }
  if appVersion == "" {
    appVersion = config.GetConfig().AppVersion
  }
  client := &CampusClient{
    idToken:    idToken,
    appVersion: appVersion,
    serverUrl:  SERVER_URL,
    timeout:    TIMEOUT_SEC,
    headers:    grpcHeader,
  }
  client.headers["x-app-version"] = appVersion
  return client
}

func (c *CampusClient) DoLogin() {
  var opts []grpc.DialOption
  creds, err := credentials.NewClientTLSFromFile("roots.pem", "")
  if err != nil {
    panic(err)
  }
  opts = append(
    opts,
    grpc.WithTransportCredentials(creds),
    grpc.WithUserAgent("grpc-csharp/2.37.0-dev grpc-c/15.0.0 (android; chttp2)"),
    grpc.WithDefaultCallOptions(
      grpc.CallContentSubtype("qua"),
    ),
  )
  c.conn, err = grpc.NewClient(c.serverUrl, opts...)
  if err != nil {
    panic(err)
  }
  defer c.conn.Close()

  c.systemFirstCheck()
  c.systemCheck()
  c.authLogin()
  c.masterGet()
  c.userGet()
  c.homeLogin()
  if c.loginBonusCheck() {
    c.loginBonusConfirm()
  }
  c.homeEnter()
  c.noticeListAll()
  if c.HomeEnterResp.PvpRateSeasonTop.StatusType != penum.PvpRateSeasonStatusType_PvpRateSeasonStatusType_OutOfTerm {
    if c.HomeEnterResp.PvpRateSeasonTop.StatusType == penum.PvpRateSeasonStatusType_PvpRateSeasonStatusType_NotAttended {
      c.pvpRateInitialize()
    }
    c.pvpRateGet()
  }
}

func (c *CampusClient) prepareContext() (*context.Context, *context.CancelFunc) {
  ctx := context.Background()
  md := metadata.New(c.headers)
  ctx = metadata.NewOutgoingContext(ctx, md)
  ctx, cancel := context.WithTimeout(ctx, time.Duration(c.timeout)*time.Second)
  return &ctx, &cancel
}
