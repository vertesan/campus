package octo

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"vertesan/campus/crypto"
	"vertesan/campus/network/hyper"
	"vertesan/campus/proto/octo"
	"vertesan/campus/utils/rich"

	"google.golang.org/protobuf/proto"
)

const OCTO_ENDPOINT = "https://api.asset.game-gakuen-idolmaster.jp/v2/pub/a/400/v/205000/list/"
const OCTO_API_KEY = "eSquJySjayO5OLLVgdTd"

func DownloadOctoList(curRevision int) *octo.Database {
  url := OCTO_ENDPOINT + fmt.Sprint(curRevision)
  headers := &http.Header{
    "User-Agent":      {"UnityPlayer/2022.3.21f1 (UnityWebRequest/1.0, libcurl/8.5.0-DEV)"},
    "Accept":          {"application/x-protobuf,x-octo-app/400"},
    "X-OCTO-KEY":      {"0jv0wsohnnsigttbfigushbtl3a8m7l5"},
    "X-Unity-Version": {"2022.3.21f1"},
  }
  rich.Info("Start to download OctoList.")
  resp, cancel, err := hyper.SendRequest(url, "GET", headers, nil, 30, 3)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  defer cancel()

  octoDb, err := DecryptOctoList(resp.Body, 0)
  if err != nil {
    panic(err)
  }
  return octoDb
}

func DecryptOctoList(reader io.Reader, offset int) (*octo.Database, error) {
  if offset > 0 {
    nothing := make([]byte, offset)
    if _, err := io.ReadFull(reader, nothing); err != nil {
      return nil, err
    }
  }
  iv := make([]byte, 16)
  if _, err := io.ReadFull(reader, iv); err != nil {
    return nil, err
  }
  key := sha256.Sum256([]byte(OCTO_API_KEY))

  buf := &bytes.Buffer{}
  crypto.Decrypt(key[:], iv, reader, buf)
  octoDb := &octo.Database{}
  if err := proto.Unmarshal(buf.Bytes(), octoDb); err != nil {
    return nil, err
  }
  return octoDb, nil
}
