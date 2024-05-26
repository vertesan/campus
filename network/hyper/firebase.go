package hyper

import (
  "bytes"
  "errors"
  "io"
  "net/http"
  "vertesan/campus/utils/rich"

  "github.com/goccy/go-json"
)

const FIREBASE_ENDPOINT = "https://securetoken.googleapis.com/v1/token?key=AIzaSyCe_vKRW5Pc0rTXFksur-ZCDb_kRCxNhng"

func GetFirebaseToken(refreshToken string) (string, string) {
  if refreshToken == "" {
    panic(errors.New("refreshToken must be given"))
  }

  headers := http.Header{
    "Content-Type":        {"application/json"},
    "X-Android-Package":   {"com.bandainamcoent.idolmaster_gakuen"},
    "X-Android-Cert":      {"D05C4DC398804FEB2ADF30E8854500D2834EAFED"},
    "Accept-Language":     {"en-US"},
    "X-Client-Version":    {"Android/Fallback/X22003001/FirebaseCore-Android"},
    "X-Firebase-GMPID":    {"1:544792984478:android:9124c0c2fb92a23780b48f"},
    "X-Firebase-Client":   {"H4sIAAAAAAAAAKtWykhNLCpJSk0sKVayio7VUSpLLSrOzM9TslIyUqoFAFyivEQfAAAA"},
    "X-Firebase-AppCheck": {"eyJlcnJvciI6IlVOS05PV05fRVJST1IifQ=="},
    "User-Agent":          {"Dalvik/2.1.0 (Linux; U; Android 13; Pixel 7 Build/TQ3A.230901.001.C2)"},
    // "Accept-Encoding":     {"gzip"},
  }
  reqBody, err := json.Marshal(map[string]string{
    "grantType":    "refresh_token",
    "refreshToken": refreshToken,
  })
  if err != nil {
    panic(err)
  }
  reqBodyReader := bytes.NewReader(reqBody)
  res, err := SendRequest(FIREBASE_ENDPOINT, "POST", headers, reqBodyReader, 10, 1)
  if err != nil {
    panic(err)
  }
  defer res.Body.Close()

  buf := &bytes.Buffer{}
  if _, err := io.Copy(buf, res.Body); err != nil {
    panic(err)
  }

  resMap := make(map[string]string)
  if err := json.Unmarshal(buf.Bytes(), &resMap); err != nil {
    panic(err)
  }

  refToken, ok := resMap["refresh_token"]
  if !ok {
    rich.ErrorThenThrow("refresh_token is absent in firebase response")
  }

  idToken, ok := resMap["id_token"]
  if !ok {
    rich.ErrorThenThrow("id_token is absent in firebase response")
  }

  return idToken, refToken
}
