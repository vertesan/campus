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

// unused but preserved for potential future use
const OCTO_ENDPOINT_ANDROID = "https://api.asset.game-gakuen-idolmaster.jp/v2/pub/a/400/v/205100/list/"
const OCTO_API_KEY_ANDROID = "eSquJySjayO5OLLVgdTd"

const OCTO_ENDPOINT = "https://api.asset.game-gakuen-idolmaster.jp/v2/pub/a/400/v/705100/list/"
const OCTO_API_KEY = "x5HFaJCJywDyuButLM0f"

func DownloadOctoList(curRevision int) *octo.Database {
  url := OCTO_ENDPOINT + fmt.Sprint(curRevision)
  headers := &http.Header{
    "User-Agent":      {"UnityPlayer/6000.0.77f1 (UnityWebRequest/1.0, libcurl/8.10.1-DEV)"},
    "Accept":          {"application/x-protobuf,x-octo-app/400"},
    "X-OCTO-KEY":      {"0jv0wsohnnsigttbfigushbtl3a8m7l5"},
    "X-Unity-Version": {"6000.0.77f1"},
  }
  rich.Info("Start to download OctoList.")
  resp, cancel, err := hyper.SendRequest(url, "GET", headers, nil, 30, 3)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  defer cancel()

  contentLen := resp.ContentLength
  octoDb, err := DecryptOctoList(resp.Body, 0, contentLen)
  if err != nil {
    panic(err)
  }
  return octoDb
}

func DecryptOctoList(reader io.Reader, offset int, contentLen int64) (*octo.Database, error) {
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

  // discard gap if any
  gap := contentLen - (contentLen-int64(offset)-16)/16*16 - int64(offset) - 16
  redundant := make([]byte, gap)
  if _, err := io.ReadFull(reader, redundant); err != nil {
    return nil, err
  }

  buf := &bytes.Buffer{}
  crypto.Decrypt(key[:], iv, reader, buf)
  octoDb := &octo.Database{}
  if err := proto.Unmarshal(buf.Bytes(), octoDb); err != nil {
    return nil, err
  }
  return octoDb, nil
}
