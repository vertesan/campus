package analyser

import (
  "bufio"
  "io"
  "os"
  "regexp"
  "slices"
  "strings"
  "vertesan/campus/rich"
)

type ProtoTree struct {
  prefix   string
  name     string
  level    int
  category Category
  children map[string]*ProtoTree
}

type Category int

const (
  Enum Category = 1 + iota
  Common
  Master
  Api
  Transaction
  Nested
  Root
)

var (
  dumpFilePath = "cache/dump.cs"

  rootClassPtnStr = `: Campus\.Common\.Proto\.Client\.$which\npublic sealed class (?<className>\w+)[\s\S]+?\n}\n\n// Namespace`
  commonClassPtn  = regexp.MustCompile(strings.Replace(rootClassPtnStr, "$which", "Common", 1))
  masterClassPtn  = regexp.MustCompile(strings.Replace(rootClassPtnStr, "$which", "Master", 1))
  // apiClassPtn         = regexp.MustCompile(strings.Replace(generalClassString, "$which", "Api", 1))
  // transactionClassPtn = regexp.MustCompile(strings.Replace(generalClassString, "$which", "Transaction", 1))

  nestedClassPtnPrefixStr = `: Campus\.Common\.Proto\.Client\.$which\n`
  nestedClassPtnStr       = `public sealed class $nestedClassName : I[\s\S]+?\n}\n\n// Namespace`

  generalColumnPtn = regexp.MustCompile(`public const int \w+ = (?<columnVal>\d+);[\s\S]*?private (readonly )?(?<type>[\w<>\.]+) (?<name>\w+);`)

  generalClassTemplate = "message $className {\n"

  generalColumnTemplate = "$type $columnName = $decimal;\n"

  commonHeader = "syntax = \"proto3\";\npackage pcommon;\noption go_package = \"vertesan/campus/proto/pcommon\";\nimport \"penum.proto\";\n\n"
  masterHeader = "syntax = \"proto3\";\npackage pmaster;\noption go_package = \"vertesan/campus/proto/pmaster\";\nimport \"penum.proto\";\nimport \"pcommon.proto\";\n\n"

  outEnumPath   = "cache/GeneratedProto/penum.proto"
  outCommonPath = "cache/GeneratedProto/pcommon.proto"
  outMasterPath = "cache/GeneratedProto/pmaster.proto"

  typeMap = map[string]string{
    "int":        "int32",
    "long":       "int64",
    "string":     "string",
    "double":     "double",
    "float":      "float",
    "uint":       "uint32",
    "ulong":      "uint64",
    "bool":       "bool",
    "ByteString": "bytes",
  }

  indentMap = map[int]string{
    0: "",
    1: "  ",
    2: "    ",
    3: "      ",
    4: "        ",
    5: "          ",
    6: "            ",
    7: "              ",
    8: "                ",
  }

  enumList   []string
  commonList []string
  masterList []string
)

// intelligently combine tree.prefix[.]tree.name,
func getClassPath(tree *ProtoTree, needSuffixDot bool) string {
  path := ""
  if tree.prefix != "" {
    path = tree.prefix + "." + tree.name
  } else if tree.name != "" {
    path = tree.name
  }
  if path != "" && needSuffixDot {
    path += "."
  }
  return path
}

func constructRoot(entireContent *string, category Category) *ProtoTree {
  var classPtn *regexp.Regexp
  switch category {
  case Common:
    classPtn = commonClassPtn
  case Master:
    classPtn = masterClassPtn
  default:
    rich.ErrorThenThrow("Unkown class type: %v", category)
  }
  root := &ProtoTree{
    prefix:   "",
    name:     "",
    level:    0,
    category: Root,
    children: make(map[string]*ProtoTree),
  }
  contents := classPtn.FindAllStringSubmatch(*entireContent, -1)
  for _, oneClass := range contents {
    className := oneClass[1]
    attachChild(className, root, category)
    switch category {
    case Common:
      commonList = append(commonList, className)
    case Master:
      masterList = append(masterList, className)
    default:
      rich.ErrorThenThrow("Unkown class type: %v", category)
    }
  }
  return root
}

func attachChild(classPath string, parentTree *ProtoTree, category Category) {
  trimPath := getClassPath(parentTree, true)
  remnantClassPath := strings.TrimPrefix(classPath, trimPath)
  nameSlice := strings.Split(remnantClassPath, ".")
  currentTree := parentTree
  for _, name := range nameSlice {
    if _, ok := currentTree.children[name]; !ok {
      if currentTree.children == nil {
        currentTree.children = make(map[string]*ProtoTree)
      }
      var level int
      level = currentTree.level + 1
      currentTree.children[name] = &ProtoTree{
        prefix:   getClassPath(currentTree, false),
        name:     name,
        level:    level,
        category: category,
      }
    }
    currentTree = currentTree.children[name]
  }
}

func analyzeTree(
  entireContent *string,
  rootCategory Category,
  parentTree *ProtoTree,
) *strings.Builder {
  sb := new(strings.Builder)

  for _, tree := range parentTree.children {
    classPath := getClassPath(tree, false)
    if classPath == "" {
      rich.ErrorThenThrow("Empty classPath.")
    }
    sb.WriteString(indentMap[tree.level - 1] + strings.Replace(generalClassTemplate, "$className", tree.name, 1))

    classSearchPtnStr := strings.Replace(nestedClassPtnStr, "$nestedClassName", classPath, 1)
    prefix := ""
    switch tree.category {
    case Common:
      prefix = strings.Replace(nestedClassPtnPrefixStr, "$which", "Common", 1)
    case Master:
      prefix = strings.Replace(nestedClassPtnPrefixStr, "$which", "Master", 1)
    }
    classSearchPtnStr = prefix + classSearchPtnStr
    classSearchPtn := regexp.MustCompile(classSearchPtnStr)

    contents := classSearchPtn.FindAllStringSubmatch(*entireContent, -1)
    // search for entire class context
    for _, oneClass := range contents {
      content := oneClass[0]
      // search for every single message
      for _, subMatches := range generalColumnPtn.FindAllStringSubmatch(content, -1) {
        line := generalColumnTemplate
        columnVal := subMatches[1]
        typeName := subMatches[3]
        columnName := subMatches[4]
        isRepeated := false
        // if typeName is a list, prune the redundant characters
        if strings.HasPrefix(typeName, "RepeatedField<") {
          typeName = strings.TrimPrefix(typeName, "RepeatedField<")
          typeName = strings.TrimSuffix(typeName, ">")
          isRepeated = true
        }
        if mappedType, ok := typeMap[typeName]; ok {
          // in case of primitive types
          typeName = mappedType
        } else {
          // in case of user defined types
          // first, check if it contains ".", if yes, it's highly likely a nested message
          if strings.Contains(typeName, ".") {
            attachChild(typeName, tree, Nested)
          } else {
            // if not, it can be an imported type
            if slices.Contains(enumList, typeName) {
              typeName = "penum." + typeName
            } else if slices.Contains(commonList, typeName) && rootCategory != Common {
              typeName = "pcommon." + typeName
            } else if slices.Contains(masterList, typeName) && rootCategory != Master {
              typeName = "pmaster." + typeName
            }
          }
        }
        if isRepeated {
          typeName = "repeated " + typeName
        }
        columnName = strings.TrimSuffix(columnName, "_")
        line = strings.Replace(line, "$type", typeName, 1)
        line = strings.Replace(line, "$columnName", columnName, 1)
        line = strings.Replace(line, "$decimal", columnVal, 1)

        sb.WriteString(indentMap[tree.level] + line)
      }
    }

    nestedSb := analyzeTree(entireContent, rootCategory, tree)
    sb.WriteString(nestedSb.String())
    sb.WriteString(indentMap[tree.level - 1] + "}\n")
  }
  return sb
}

func analyzeFile(
  entireContent *string,
  category Category,
  outPath string,
) {
  // create a to be generated file
  protoFile, err := os.Create(outPath)
  if err != nil {
    panic(err)
  }
  defer protoFile.Close()
  buf := bufio.NewWriter(protoFile)
  // write prefixes
  var header string
  switch category {
  case Enum:
    header = enumHeader
  case Common:
    header = commonHeader
  case Master:
    header = masterHeader
  default:
    rich.ErrorThenThrow("Unkown type of proto file: %v.", category)
  }
  buf.WriteString(header)

  var sb *strings.Builder
  if category == Enum {
    sb = AnalyzeEnum(entireContent)
  } else {
    root := constructRoot(entireContent, category)
    sb = analyzeTree(entireContent, category, root)
  }


  buf.WriteString(sb.String())
  // flush
  if err := buf.Flush(); err != nil {
    panic(err)
  }
}

func getWrap(category Category) string {
  switch category {
  case Common:
    return "message Common {\n"
  case Master:
    return "message Master {\n"
  default:
    return ""
  }
}

func Analyze() {
  f, err := os.Open(dumpFilePath)
  if err != nil {
    panic(err)
  }
  defer f.Close()
  bufr := bufio.NewReader(f)
  // read dump file to string builder
  sb := new(strings.Builder)
  _, err = io.Copy(sb, bufr)
  if err != nil {
    panic(err)
  }
  entireContent := sb.String()
  analyzeFile(&entireContent, Enum, outEnumPath)
  analyzeFile(&entireContent, Common, outCommonPath)
  analyzeFile(&entireContent, Master, outMasterPath)
}
