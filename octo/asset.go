package octo

import (
  "io"
)

func Deobfuscate(reader io.Reader, abName string) []byte {
  raw, err := io.ReadAll(reader)
  if err != nil {
    panic(err)
  }
  return cryptByString(raw, abName, 0, 0, 256)
}

func cryptByString(input []byte, maskString string, offset int, streamPos int, headerLength int) []byte {
  maskStringLength := len(maskString)
  bytesLength := maskStringLength << 1
  result := input
  maskBytes := stringToMaskBytes(maskString, maskStringLength, bytesLength)
  for i := 0; streamPos+i < headerLength; i++ {
    result[offset+i] ^= maskBytes[streamPos+i-(streamPos+i)/bytesLength*bytesLength]
  }
  return result
}

func stringToMaskBytes(maskString string, maskStringLength int, bytesLength int) []byte {
  maskBytes := make([]byte, bytesLength)
  if maskStringLength >= 1 {
    i := 0
    j := 0
    k := bytesLength - 1
    for maskStringLength != j {
      charJ := maskString[j]
      j += 1
      maskBytes[i] = charJ
      i += 2
      charJ = ^charJ & 0xFF
      maskBytes[k] = charJ
      k -= 2
    }
  }
  if bytesLength >= 1 {
    l := bytesLength
    var v13 byte = 0x9B
    m := bytesLength
    pointer := 0
    for m != 0 {
      v16 := maskBytes[pointer]
      pointer += 1
      m -= 1
      v13 = (((v13 & 1) << 7) | (v13 >> 1)) ^ v16
    }
    b := 0
    for l != 0 {
      l -= 1
      maskBytes[b] ^= v13
      b += 1
    }
  }
  return maskBytes
}
