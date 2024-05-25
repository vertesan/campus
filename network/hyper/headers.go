package hyper

import "net/http"

var (
  masterHeader = http.Header{
    "User-Agent":      {"UnityPlayer/2022.3.21f1 (UnityWebRequest/1.0, libcurl/8.5.0-DEV)"},
    "Accept":          {"*/*"},
    "Accept-Encoding": {"deflate, gzip"},
    "X-Unity-Version": {"2022.3.21f1"},
  }
)
