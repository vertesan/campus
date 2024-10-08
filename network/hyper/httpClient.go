package hyper

import (
  "context"
  "fmt"
  "io"
  "net/http"
  "time"
  "vertesan/campus/utils/rich"
)

var (
  transport = &http.Transport{
    DisableKeepAlives: true,
  }
  client = &http.Client{
    // Timeout:   10 * time.Second,
    Transport: transport,
  }
)

func SendRequest(url string, method string, header *http.Header, body io.Reader, timeout int, maxTry int) (*http.Response, context.CancelFunc, error) {
  for i := range maxTry {
    res, cancel, err := sendRequestInternal(url, method, header, body, timeout)
    if err != nil {
      rich.Error("Failed to request %q, attempts(%d/%d).", url, i+1, maxTry)
      if i+1 >= maxTry {
        rich.Error("Max retries exhausted when requesting %q.", url)
        return nil, cancel, err
      }
      continue
    }
    return res, cancel, nil
  }
  return nil, nil, fmt.Errorf("max retries exhausted when requesting %q", url)
}

func sendRequestInternal(url string, method string, header *http.Header, body io.Reader, timeout int) (*http.Response, context.CancelFunc, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
  // if context get canceled, response becomes invalid
  // defer cancel()
  request, err := http.NewRequestWithContext(ctx, method, url, body)
  if err != nil {
    return nil, cancel, err
  }
  request.Header = *header
  res, err := client.Do(request)
  if err != nil {
    return nil, cancel, err
  }
  if res.StatusCode != 200 {
    if err := res.Body.Close(); err != nil {
      panic(err)
    }
    rich.Error("Get an abnormal status when requesting %q.", url)
    return nil, cancel, fmt.Errorf("status code: %d, message: %v", res.StatusCode, res.Status)
  }
  return res, cancel, nil
}
