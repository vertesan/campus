package rpc

import (
	"context"
	"fmt"
	"time"
	papi "vertesan/campus/proto/papi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
)

const SERVER_URL = "api.game-gakuen-idolmaster.jp"
const TIMEOUT_SEC = 10

func init() {
  // register QuaCodec
  encoding.RegisterCodec(&QuaCodec{})
}

func systemFirstCheck(chn *grpc.ClientConn) {
  client := papi.NewSystemClient(chn)
  ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_SEC*time.Second)
  defer cancel()
  req := &papi.SystemCheckRequest{}
  systemCheckResp, err := client.Check(ctx, req)
  if err != nil {
    panic(err)
  }
  fmt.Printf("%s", systemCheckResp.OctoHost)
}

func RunScenarios() {
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

  conn, err := grpc.NewClient(SERVER_URL, opts...)
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  systemFirstCheck(conn)
}
