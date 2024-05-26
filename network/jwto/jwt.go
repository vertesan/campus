package jwto

import (
  "errors"
  "time"

  "github.com/golang-jwt/jwt/v5"
)

func IsJwtExpired(jwtString string) bool {
  // Parse the token without verifying the signature
  token, _, err := new(jwt.Parser).ParseUnverified(jwtString, jwt.MapClaims{})
  if err != nil {
    panic(err)
  }

  // Extract the expiration time
  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    panic(errors.New("invalid token claims"))
  }
  exp, ok := claims["exp"].(float64)
  if !ok {
    panic(errors.New("invalid 'exp' claim in token"))
  }

  currentEpoch := time.Now().Unix()
  // let's dial forward the clock 5 minutes
  return currentEpoch+300 >= int64(exp)
}
