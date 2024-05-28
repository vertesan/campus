package octo

import (
  "io"
  "net/http"
  "os"
  "path"
  "slices"
  "strings"
  "vertesan/campus/config"
  "vertesan/campus/network/hyper/downloader"
  "vertesan/campus/proto/octo"
  "vertesan/campus/utils/rich"

  "github.com/goccy/go-json"
  "google.golang.org/protobuf/encoding/protojson"
)

const OCTO_CACHE_FILE_PATH = "cache/octo_cache.json"
const OCTO_LOCAL_RECORD_PATH = "cache/octo_record.json"
const OCTO_DL_RECORD_PATH = "cache/octo_downloaded.json"
const OCTO_RAW_DIR = "cache/raw"
const OCTO_ASSET_DIR = "cache/assets"

type OctoManager struct {
  OctoDb *octo.Database
}

func (m *OctoManager) saveOctoJsonCache() {
  marshalOpt := protojson.MarshalOptions{
    Multiline:      true,
    Indent:         "  ",
    UseEnumNumbers: false,
  }
  jsonBytes, err := marshalOpt.Marshal(m.OctoDb)
  if err != nil {
    panic(err)
  }
  if err := os.WriteFile(OCTO_CACHE_FILE_PATH, jsonBytes, 0644); err != nil {
    panic(err)
  }
}

func (m *OctoManager) saveLocalJson(record map[string]string, path string) {
  jsonBytes, err := json.MarshalIndent(record, "", "  ")
  if err != nil {
    panic(err)
  }
  if err := os.WriteFile(path, jsonBytes, 0644); err != nil {
    panic(err)
  }
}

func (m *OctoManager) Work(keepRaw bool) {
  cfg := config.GetConfig()
  curRevision := cfg.OctoCacheRevision
  rich.Info("Current octo revision: %d.", curRevision)
  m.OctoDb = DownloadOctoList(curRevision)
  rich.Info("Server octo revision: %d.", m.OctoDb.Revision)
  if int(m.OctoDb.Revision) == curRevision {
    rich.Info("Octo revision is already up to date, skip downloading assets.")
    return
  }
  if err := os.MkdirAll(OCTO_RAW_DIR, 0755); err != nil {
    panic(err)
  }
  if err := os.MkdirAll(OCTO_ASSET_DIR, 0755); err != nil {
    panic(err)
  }

  m.saveOctoJsonCache()
  urlFormat := m.OctoDb.UrlFormat

  // load local record md5
  localRecordBytes, err := os.ReadFile(OCTO_LOCAL_RECORD_PATH)
  if err != nil {
    if !os.IsNotExist(err) {
      panic(err)
    }
  }

  localRecord := make(map[string]string)
  downloaded := make(map[string]string)
  if localRecordBytes != nil {
    json.Unmarshal(localRecordBytes, &localRecord)
  }

  // prepare downloader
  assetDownloadHeader := &http.Header{
    "User-Agent": {"Dalvik/2.1.0 (Linux; U; Android 11; GM1910 Build/RKQ1.201022.002)"},
  }
  dler := downloader.NewDownloader(600, assetDownloadHeader, OCTO_RAW_DIR, 20)
  entries := []*downloader.Entry{}

  // add assetbundles to entry list
  for _, asset := range m.OctoDb.AssetBundleList {
    // check file MD5 first. if already exists, skip downloading
    md5, ok := localRecord[asset.Name]
    if ok && asset.Md5 == md5 {
      rich.Warning("The MD5 of AssetBundle %q matches one of the local files, skip downloading.", asset.Name)
      continue
    } else {
      // else, write it into map
      localRecord[asset.Name] = asset.Md5
      downloaded[asset.Name] = asset.ObjectName
    }
    url := strings.Replace(urlFormat, "{o}", asset.ObjectName, 1)
    entries = append(entries, &downloader.Entry{
      Url:          url,
      SaveFileName: asset.Name,
      Type:         "ab",
    })
  }
  // add resources to entry list
  for _, resource := range m.OctoDb.ResourceList {
    // check file MD5 first. if already exists, skip downloading
    md5, ok := localRecord[resource.Name]
    if ok && resource.Md5 == md5 {
      rich.Warning("The MD5 of Resource %q matches one of the local files, skip downloading.", resource.Name)
      continue
    } else {
      // else, write it into map
      localRecord[resource.Name] = resource.Md5
      downloaded[resource.Name] = resource.ObjectName
    }
    url := strings.Replace(urlFormat, "{o}", resource.ObjectName, 1)
    entries = append(entries, &downloader.Entry{
      Url:          url,
      SaveFileName: resource.Name,
      Type:         "rs",
    })
  }
  if len(entries) == 0 {
    rich.Warning("Entry list is empty, nothing to download.")
    return
  }
  dler.Entries = entries
  // download!
  if err := dler.DownloadAll(); err != nil {
    panic(err)
  }

  // deobfuscate
  for _, entry := range entries {
    if entry.Type != "ab" {
      continue
    }
    rawFs, err := os.Open(path.Join(OCTO_RAW_DIR, entry.SaveFileName))
    if err != nil {
      panic(err)
    }
    signature := make([]byte, 5)
    if _, err := rawFs.Read(signature); err != nil {
      panic(err)
    }
    // if not obfuscated
    if slices.Equal(signature, []byte("Unity")) {
      rich.Info("Assetbundle %q seems like is not obfuscated, simply copy it to the destination.", entry.SaveFileName)
      abFs, err := os.Create(path.Join(OCTO_ASSET_DIR, entry.SaveFileName))
      if err != nil {
        panic(err)
      }
      if _, err := io.Copy(abFs, rawFs); err != nil {
        panic(err)
      }
      if !keepRaw {
        rawFs.Close()
        if err := os.Remove(path.Join(OCTO_RAW_DIR, entry.SaveFileName)); err != nil {
          panic(err)
        }
      }
      abFs.Close()
      continue
    }
    if _, err := rawFs.Seek(0, 0); err != nil {
      panic(err)
    }
    // deobfuscate
    abBytes := Deobfuscate(rawFs, entry.SaveFileName)
    if !slices.Equal(abBytes[:5], []byte("Unity")) {
      rich.Warning("Assetbundle %q seems have a invalid signature, please check if this is intentional.", entry.SaveFileName)
    } else {
      rich.Info("Successfully deobfuscated %q.", entry.SaveFileName)
    }
    if err := os.WriteFile(path.Join(OCTO_ASSET_DIR, entry.SaveFileName), abBytes, 0644); err != nil {
      panic(err)
    }
    if !keepRaw {
      rawFs.Close()
      if err := os.Remove(path.Join(OCTO_RAW_DIR, entry.SaveFileName)); err != nil {
        panic(err)
      }
    }
    rawFs.Close()
  }
  // save this time downloads record
  m.saveLocalJson(downloaded, OCTO_DL_RECORD_PATH)
  // save local record
  m.saveLocalJson(localRecord, OCTO_LOCAL_RECORD_PATH)
}
