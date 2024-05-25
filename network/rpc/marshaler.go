package rpc

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"time"
	"vertesan/campus/crypto"
	"vertesan/campus/utils/rich"
)

const KEY = "Kb31v85u"
const COMPRESS_THRESHOLD = 2048
const UNIX_1970 = 621355968000000000

var (
  keyHash = md5.Sum([]byte(KEY))
)

func getTicks() uint64 {
  now := uint64(time.Now().UnixMilli())
  ticks := now*10000 + UNIX_1970
  return ticks
}

func encrypt(data []byte, iv []byte) []byte {
  src := bytes.NewReader(data)
  dst := &bytes.Buffer{}
  crypto.Encrypt(keyHash[:], iv, src, dst)
  return dst.Bytes()
}

func decrypt(data []byte, iv []byte) []byte {
  src := bytes.NewReader(data)
  dst := &bytes.Buffer{}
  crypto.Decrypt(keyHash[:], iv, src, dst)
  return dst.Bytes()
}

func Serialize(protoBytes []byte) []byte {
  if len(protoBytes) == 0 {
    return protoBytes
  }

  buf := &bytes.Buffer{}

  if len(protoBytes) > COMPRESS_THRESHOLD {
    rich.Warning("Message should be compressed.")
  }

  // compressing flag
  isCompressed := []byte{0x00}

  // header
  headerBytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(headerBytes, getTicks())

  // body
  iv := append([]byte(KEY), headerBytes...)
  ivHash := md5.Sum(iv)
  bodyBytes := encrypt(protoBytes, ivHash[:])

  // message length
  messageLength := uint32(len(bodyBytes) + len(headerBytes) + 4)
  messageLengthBytes := make([]byte, 4)
  binary.BigEndian.PutUint32(messageLengthBytes, messageLength)

  // write
  if _, err := buf.Write(isCompressed); err != nil {
    panic(err)
  }
  if _, err := buf.Write(messageLengthBytes); err != nil {
    panic(err)
  }
  if err := buf.WriteByte(0x0a); err != nil {
    panic(err)
  }
  if _, err := buf.Write([]byte{0x00, 0x00}); err != nil {
    panic(err)
  }
  if err := buf.WriteByte(0x08); err != nil {
    panic(err)
  }
  if _, err := buf.Write(headerBytes); err != nil {
    panic(err)
  }
  if _, err := buf.Write(bodyBytes); err != nil {
    panic(err)
  }
  return buf.Bytes()
}

func Deserialize(raw []byte) []byte {
  if len(raw) == 0 {
    return raw
  }
  
  buf := bytes.NewBuffer(raw)

  // header length + 2
  if _, err := buf.ReadByte(); err != nil {
    panic(err)
  }

  // 00 00
  nothing := make([]byte, 2)
  if _, err := buf.Read(nothing); err != nil {
    panic(err)
  }

  // header length
  headerLengthByte, err := buf.ReadByte()
  if err != nil {
    panic(err)
  }
  headerLength := uint8(headerLengthByte)

  // header
  headerBytes := make([]byte, headerLength)
  if _, err := buf.Read(headerBytes); err != nil {
    panic(err)
  }

  // encrypted body
  encBytes, err := io.ReadAll(buf)
  if err != nil {
    panic(err)
  }

  // decrypt body
  iv := append([]byte(KEY), headerBytes...)
  ivHash := md5.Sum(iv)
  bodyBytes := decrypt(encBytes, ivHash[:])
  return bodyBytes
}
