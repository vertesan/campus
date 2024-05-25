package hyper

import (
  "net/http"
  "time"
)

var (
  transport = &http.Transport{
    DisableKeepAlives: true,
  }
  client = &http.Client{
    // Timeout seconds for downloading a single file.
    // As of 2024.2, the largest file size is 363MB (one of the feslive videos),
    // make sure this value is large enough for downloading large files.
    Timeout:   1200 * time.Second,
    Transport: transport,
  }
)

func prepareGetRequest(url string, header http.Header) *http.Request {
  request, err := http.NewRequest("GET", url, nil)
  if err != nil {
    panic(err)
  }
  request.Header = header
  return request
}
