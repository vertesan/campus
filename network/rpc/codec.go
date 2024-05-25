package rpc

import (
  "fmt"

  "google.golang.org/protobuf/proto"
  "google.golang.org/protobuf/protoadapt"
)

type QuaCodec struct{}

func (g *QuaCodec) Marshal(v any) ([]byte, error) {
  vv := messageV2Of(v)
  if vv == nil {
    return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
  }
  // default codec
  protoBytes, err := proto.Marshal(vv)
  if err != nil {
    return nil, err
  }
  // customized codec
  encBytes := Serialize(protoBytes)
  return encBytes, nil
}

func (g *QuaCodec) Unmarshal(data []byte, v any) error {
  vv := messageV2Of(v)
  if vv == nil {
    return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
  }
  // customized codec
  protoBytes := Deserialize(data)
  // default codec
  return proto.Unmarshal(protoBytes, vv)
}

func messageV2Of(v any) proto.Message {
  switch v := v.(type) {
  case protoadapt.MessageV1:
    return protoadapt.MessageV2Of(v)
  case protoadapt.MessageV2:
    return v
  }
  return nil
}

func (g *QuaCodec) Name() string {
  return "qua"
}
