package rpc

import "net/http"

var (
  loginHeader = &http.Header{
    "x-app-version":        {"1.0.2"},
    "x-platform":           {"Android"},
    "x-device-name":        {"Pixel 7"},
    "x-os-version":         {"Android OS 13 / API-33 (TQ3A.230901.001.C2)"},
    "accept-language":      {"ja"},
    "x-a-ise":              {"False"},
    "x-a-isda":             {"False"},
    "x-a-isop":             {"False"},
    "x-a-isr":              {"False"},
    "x-a-s":                {"adb9f0cb6bf3aa06f99cc15826dbbca7251f88369d38e1b4da5c5bc46ef3e951"},
    "x-ad-id":              {"11ca2dbbd1d87224a04500e6ddfcc03f"},
    "x-platform-user-id":   {"00000000-0000-0000-0000-000000000000"},
    "te":                   {"trailers"},
    "content-type":         {"application/grpc"},
    "user-agent":           {"grpc-csharp/2.37.0-dev grpc-c/15.0.0 (android; chttp2)"},
    "grpc-accept-encoding": {"identity,deflate,gzip"},
    "accept-encoding":      {"identity,gzip"},
    "grpc-timeout":         {"10S"},
  }
)
