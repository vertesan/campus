package octo

import (
  "bytes"
  "io"
  "vertesan/campus/crypto"
  "vertesan/campus/proto/octo"

  "google.golang.org/protobuf/proto"
)

func DecryptManifest(reader io.Reader, key []byte, iv []byte, offset int) (*octo.Database, error) {
  buf := &bytes.Buffer{}
  nothing := make([]byte, offset)
  if _, err := reader.Read(nothing); err != nil {
    panic(err)
  }
  crypto.Decrypt(key, iv, reader, buf)
  octoDb := &octo.Database{}
  err := proto.Unmarshal(buf.Bytes()[16:], octoDb)
  return octoDb, err
}
