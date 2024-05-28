package downloader

import (
  "bufio"
  "bytes"
  "context"
  "errors"
  "io"
  "net/http"
  "os"
  "path"
  "time"
  "vertesan/campus/utils/rich"

  "golang.org/x/sync/semaphore"
)

type Entry struct {
  Url          string
  SaveFileName string
  Header       *http.Header
  Type         string
}

type Downloader struct {
  Client         *http.Client
  Header         *http.Header
  Counter        *SafeCounter
  SaveDir        string
  MaxConcurrency int64
  MaxSingleRetry int
  Entries        []*Entry
}

func NewDownloader(timeout int, header *http.Header, saveDir string, maxThread int64) *Downloader {
  d := &Downloader{
    Client: &http.Client{
      Timeout: time.Duration(timeout) * time.Second,
      Transport: &http.Transport{
        DisableKeepAlives: true,
      },
    },
    Header:         header,
    Counter:        &SafeCounter{},
    SaveDir:        saveDir,
    MaxConcurrency: maxThread,
    MaxSingleRetry: 3,
    Entries:        nil,
  }
  return d
}

func (d *Downloader) SetEntries(ent []*Entry) {
  d.Entries = ent
}

func (d *Downloader) DownloadAll() error {
  if len(d.Entries) == 0 {
    return errors.New("entries is empty, nothing to download")
  }
  if err := os.MkdirAll(d.SaveDir, 0755); err != nil {
    panic(err)
  }
  sem := semaphore.NewWeighted(d.MaxConcurrency)
  quantity := len(d.Entries)
  ctx := context.Background()

  rich.Info("Start to download %d objects.", quantity)
  for _, entry := range d.Entries {
    if err := sem.Acquire(ctx, 1); err != nil {
      panic(err)
    }
    go d.downloadOneAsync(entry, sem, quantity)
  }
  // wait all concurrencies completed
  if err := sem.Acquire(ctx, d.MaxConcurrency); err != nil {
    panic(err)
  }
  rich.Info("Successfully downloaded all objects.")
  return nil
}

func (d *Downloader) downloadOneAsync(
  entry *Entry,
  sem *semaphore.Weighted,
  total int,
) {
  if sem != nil {
    defer sem.Release(1)
  }
  d.downloadOne(entry, total)
}

func (d *Downloader) downloadOne(
  entry *Entry,
  total int,
) {
  // prepare request
  request, err := http.NewRequest("GET", entry.Url, nil)
  if err != nil {
    panic(err)
  }
  // set request header
  if entry.Header == nil {
    request.Header = *d.Header
  } else {
    request.Header = *entry.Header
  }

  for i := range d.MaxSingleRetry {
    res, err := d.Client.Do(request)
    if err != nil {
      rich.Error("%v", err)
      rich.Warning("An internal error was occurred when downloading %v (%s), retrying...(%d/%d)", entry.SaveFileName, entry.Url, i+1, d.MaxSingleRetry)
      continue
    }
    if res.StatusCode != 200 {
      rich.Error("Status code: %d, message: %v.", res.StatusCode, res.Status)
      rich.Warning("A HTTP error was occurred when downloading %v (%s), retrying...(%d/%d)", entry.Url, entry.SaveFileName, i+1, d.MaxSingleRetry)
      if err := res.Body.Close(); err != nil {
        panic(err)
      }
      continue
    }
    // save response stream to a file
    fs, err := os.Create(path.Join(d.SaveDir, entry.SaveFileName))
    if err != nil {
      panic(err)
    }
    defer fs.Close()
    bufw := bufio.NewWriter(fs)
    if _, err := bufw.ReadFrom(res.Body); err != nil {
      rich.Warning("An internal error was occurred when reading body from %v (%s), retrying...(%d/%d)", entry.Url, entry.SaveFileName, i+1, d.MaxSingleRetry)
      continue
    }
    if err := bufw.Flush(); err != nil {
      panic(err)
    }

    d.Counter.Increase()
    rich.Info("(%d/%d) Download completed: %q (%v).", d.Counter.Value(), total, entry.SaveFileName, entry.Url)
    res.Body.Close()
    return
  }
  // max retry exhausted
  rich.ErrorThenThrow("Max retries exhausted when downloading %v. Will be stopping process.", request.URL)
}

func (d *Downloader) DownloadToMem(entry *Entry, total int) *bytes.Buffer {
  // prepare request
  request, err := http.NewRequest("GET", entry.Url, nil)
  if err != nil {
    panic(err)
  }
  // set request header
  if entry.Header == nil {
    request.Header = *d.Header
  } else {
    request.Header = *entry.Header
  }

  for i := range d.MaxSingleRetry {
    res, err := d.Client.Do(request)
    if err != nil {
      rich.Error("%v", err)
      rich.Warning("An internal error was occurred when downloading %v, retrying...(%d/%d)", entry.Url, i+1, d.MaxSingleRetry)
      continue
    }
    if res.StatusCode != 200 {
      rich.Error("Status code: %d, message: %v.", res.StatusCode, res.Status)
      rich.Warning("A HTTP error was occurred when downloading %v, retrying...(%d/%d)", entry.Url, i+1, d.MaxSingleRetry)
      if err := res.Body.Close(); err != nil {
        panic(err)
      }
      continue
    }
    // save response stream to buffer
    buf := &bytes.Buffer{}
    if _, err := io.Copy(buf, buf); err != nil {
      panic(err)
    }

    d.Counter.Increase()
    rich.Info("(%d/%d) Download completed: %q(%v).", d.Counter.Value(), total, entry.SaveFileName, entry.Url)
    res.Body.Close()
    return buf
  }
  // max retry exhausted
  rich.ErrorThenThrow("Max retries exhausted when downloading %v. Will be stopping process.", request.URL)
  return nil
}
