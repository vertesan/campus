package hyper

import (
  "bufio"
  "fmt"
  "os"
  "vertesan/campus/proto/mastertag"
  "vertesan/campus/utils/rich"
)

const MASTER_RAW_DIR = "cache/masterRaw"

func DownloadOneMasterRaw(masterPack *mastertag.MasterTagPack) {
  url := masterPack.DownloadUrl
  request := prepareGetRequest(url, masterHeader)
  res, err := client.Do(request)
  if err != nil {
    panic(err)
  }
  if res.StatusCode != 200 {
    if err := res.Body.Close(); err != nil {
      panic(err)
    }
    rich.ErrorThenThrow("Status code: %d, message: %v.", res.StatusCode, res.Status)
  }
  defer res.Body.Close()

  // save response stream to a file
  if err := os.MkdirAll(MASTER_RAW_DIR, 0755); err != nil {
    panic(err)
  }
  fs, err := os.Create(fmt.Sprintf("%v/%v", MASTER_RAW_DIR, masterPack.Type))
  if err != nil {
    panic(err)
  }
  defer fs.Close()
  bufw := bufio.NewWriter(fs)
  if _, err := bufw.ReadFrom(res.Body); err != nil {
    panic(err)
  }
  if err := bufw.Flush(); err != nil {
    panic(err)
  }
}
