package analyser

import "regexp"

var (
  rootClassPtnStr         = `: Campus\.Common\.Proto\.Client\.$which\npublic sealed class (?<className>\w+)[\s\S]+?\n}\n\n// Namespace`
  nestedClassPtnPrefixStr = `: Campus\.Common\.Proto\.Client\.$which\n`
  nestedClassPtnStr       = `public sealed class $nestedClassName : I[\s\S]+?\n}\n\n// Namespace`

  generalColumnPtn = regexp.MustCompile(`public const int \w+ = (?<columnVal>\d+);[\s\S]*?private (readonly )?(?<type>[\w<>\.]+) (?<name>\w+);`)

  generalClassTemplate = "message $className {\n"

  generalColumnTemplate = "$type $columnName = $decimal;\n"

  commonHeader = `syntax = "proto3";
package pcommon;
option go_package = "vertesan/campus/proto/pcommon";
import "penum.proto";

`
  masterHeader = `syntax = "proto3";
package pmaster;
option go_package = "vertesan/campus/proto/pmaster";
import "penum.proto";
import "pcommon.proto";

`
  transactionHeader = `syntax = "proto3";
package ptransaction;
option go_package = "vertesan/campus/proto/ptransaction";
import "penum.proto";
import "pcommon.proto";

`
  apiCommonHeader = `syntax = "proto3";
package papicommon;
option go_package = "vertesan/campus/proto/papicommon";
import "penum.proto";
import "pcommon.proto";
import "ptransaction.proto";
import "pmaster.proto";

`
  apiHeader = `syntax = "proto3";
package client.api;
option go_package = "vertesan/campus/proto/papi";
import "penum.proto";
import "pcommon.proto";
import "pmaster.proto";
import "papicommon.proto";

`

  mappingHeader   = "// Generated code. DO NOT EDIT!\npackage mapping\n\nimport \"vertesan/campus/proto/pcommon\"\nimport \"vertesan/campus/proto/pmaster\"\nimport \"google.golang.org/protobuf/reflect/protoreflect\"\n\nvar (\n  ProtoMap = map[string]protoreflect.ProtoMessage{\n"
  mappingTemplate = "    \"$category.$className\": &$package.$className{},\n"
  mappingTail     = "  }\n)\n"
)
