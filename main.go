package main

import (
	"crypto/md5"
	"os"
	"vertesan/campus/analyser"
	"vertesan/campus/master"
	"vertesan/campus/proto/mastertag"
	"vertesan/campus/octo"
	"vertesan/campus/rich"

	"google.golang.org/protobuf/encoding/protojson"
)

const OCTO_CACHE_FILE = "cache/octocacheevai"
const MASTER_TAG_FILE = "cache/masterGetDec240516175118527.bin"

var (
  marshalOptions = &protojson.MarshalOptions{
    Multiline:      true,
    Indent:         "  ",
    AllowPartial:   true,
    UseProtoNames:  true,
    UseEnumNumbers: true,
  }
)

func decryptOctoManifest() {
  // open encrypted octo cache file
  octoFs, err := os.Open(OCTO_CACHE_FILE)
  if err != nil {
    panic(err)
  }
  keyMd5 := md5.Sum([]byte("1nuv9td1bw1udefk"))
  ivMd5 := md5.Sum([]byte("LvAUtf+tnz"))
  octoDb, err := octo.DecryptManifest(octoFs, keyMd5[:], ivMd5[:], 1)
  if err != nil {
    panic(err)
  }
  // convert protobuf object to json string

  octoJson, err := marshalOptions.Marshal(octoDb)
  if err != nil {
    panic(err)
  }
  // write manifest json file
  if err := os.WriteFile("cache/octo.json", octoJson, 0644); err != nil {
    panic(err)
  }
}

func getMasterDb() *mastertag.MasterGetResponse {
  fs, err := os.Open(MASTER_TAG_FILE)
  if err != nil {
    panic(err)
  }
  masterGetResp, err := master.UnmarshalPlain(fs)
  if err != nil {
    panic(err)
  }
  masterGetRespJson, err := marshalOptions.Marshal(masterGetResp)
  if err != nil {
    panic(err)
  }
  if err := os.WriteFile("cache/masterGetResp.json", masterGetRespJson, 0644); err != nil {
    panic(err)
  }
  return masterGetResp
}

func main() {
  rich.Info("Start operation")
  // decryptOctoManifest()
  // masterGetResp := getMasterDb()
  // master.DownloadAllMaster(masterGetResp)
  // master.DecryptAll(masterGetResp)
  analyser.Analyze()
}
