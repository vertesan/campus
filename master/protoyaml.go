package master

import (
  "github.com/goccy/go-yaml"
  "google.golang.org/protobuf/encoding/protojson"
  "google.golang.org/protobuf/proto"
)

func YamlMarshal(message proto.Message) ([]byte, error) {
  data, err := protojson.MarshalOptions{
    Multiline:         false,
    AllowPartial:      false,
    UseProtoNames:     true,
    UseEnumNumbers:    false,
    EmitUnpopulated:   true,
    EmitDefaultValues: false,
  }.Marshal(message)
  if err != nil {
    return nil, err
  }
  r, err := yaml.JSONToYAML(data)
  return r, err
}

// func YamlMarshal(message proto.Message) ([]byte, error) {
//   data, err := protojson.MarshalOptions{
//     Multiline:         false,
//     AllowPartial:      false,
//     UseProtoNames:     true,
//     UseEnumNumbers:    false,
//     EmitUnpopulated:   true,
//     EmitDefaultValues: false,
//   }.Marshal(message)
//   if err != nil {
//     return nil, err
//   }
//   var jsonVal interface{}
//   if err := json.Unmarshal(data, &jsonVal); err != nil {
//     return nil, err
//   }
//   jv, err := json.Marshal(jsonVal)
//   fmt.Printf("%v", jv)
//   // Write the JSON back out as YAML
//   buffer := &bytes.Buffer{}
//   encoder := yaml.NewEncoder(buffer)
//   if err := encoder.Encode(message); err != nil {
//     return nil, err
//   }
//   b, err := yaml.Marshal(jsonVal)
//   return b, nil
// }
