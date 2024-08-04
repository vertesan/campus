package analyser

import (
  "bufio"
  "io"
  "os"
  "regexp"
  "slices"
  "sort"
  "strings"
  "vertesan/campus/utils/rich"
)

const (
  Enum Category = 1 + iota
  Common
  Master
  Transaction
  ApiCommon
  Api
  Nested
  Root
)

type Category int

type ProtoTree struct {
  prefix    string
  name      string
  level     int
  category  Category
  children  map[string]*ProtoTree
  traversed bool
}

var (
  dumpFilePath = "cache/dump.cs"

  commonClassPtn      = regexp.MustCompile(strings.Replace(rootClassPtnStr, "$which", "Common", 1))
  masterClassPtn      = regexp.MustCompile(strings.Replace(rootClassPtnStr, "$which", "Master", 1))
  transactionClassPtn = regexp.MustCompile(strings.Replace(rootClassPtnStr, "$which", "Transaction", 1))
  apiCommonClassPtn   = regexp.MustCompile(strings.Replace(rootClassPtnStr, "$which", "Api.Common", 1))
  apiClassPtn         = regexp.MustCompile(strings.Replace(rootClassPtnStr, "$which", "Api", 1))

  outDir             = "cache/GeneratedProto"
  outEnumPath        = outDir + "/penum.proto"
  outCommonPath      = outDir + "/pcommon.proto"
  outMasterPath      = outDir + "/pmaster.proto"
  outTransactionPath = outDir + "/ptransaction.proto"
  outApiCommonPath   = outDir + "/papicommon.proto"
  outApiPath         = outDir + "/papi.proto"

  tsOutDir             = "cache/GeneratedTypeScript"
  tsOutEnumPath        = tsOutDir + "/penum.ts"
  tsOutCommonPath      = tsOutDir + "/pcommon.d.ts"
  tsOutMasterPath      = tsOutDir + "/pmaster.d.ts"
  tsOutTransactionPath = tsOutDir + "/ptransaction.d.ts"
  tsOutApiCommonPath   = tsOutDir + "/papicommon.d.ts"
  tsOutApiPath         = tsOutDir + "/papi.d.ts"

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

  tsTypeMap = map[string]string{
    "int":        "number",
    "long":       "string",
    "string":     "string",
    "double":     "number",
    "float":      "number",
    "uint":       "number",
    "ulong":      "string",
    "bool":       "boolean",
    "ByteString": "string",
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

  // top level entry list
  enumList        []string
  commonList      []string
  masterList      []string
  transactionList []string
  apiCommonList   []string
  apiList         []string

  // cross reference list
  xList []string

  mappingSb      = new(strings.Builder)
  mappingOutFile = "cache/GeneratedProto/mapping.go"
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

// split classPath to prefix and name
func getPrefixAndName(classPath string) (string, string) {
  idx := strings.LastIndex(classPath, ".")
  if idx == -1 {
    return "", classPath
  }
  prefix := classPath[:idx]
  name := classPath[idx+1:]
  return prefix, name
}

func constructRoot(entireContent *string, category Category) *ProtoTree {
  var classPtn *regexp.Regexp
  var packageName string
  var categoryString string
  switch category {
  case Common:
    classPtn = commonClassPtn
    packageName = "pcommon"
    categoryString = "Common"
  case Master:
    classPtn = masterClassPtn
    packageName = "pmaster"
    categoryString = "Master"
  case Transaction:
    classPtn = transactionClassPtn
    packageName = "ptransaction"
    categoryString = "Transaction"
  case ApiCommon:
    classPtn = apiCommonClassPtn
    packageName = "papicommon"
    categoryString = "ApiCommon"
  case Api:
    classPtn = apiClassPtn
    packageName = "papi"
    categoryString = "Api"
  default:
    rich.ErrorThenThrow("Unkown class type: %v", category)
  }
  root := &ProtoTree{
    prefix:    "",
    name:      "",
    level:     0,
    category:  Root,
    children:  make(map[string]*ProtoTree),
    traversed: false,
  }
  contents := classPtn.FindAllStringSubmatch(*entireContent, -1)

  for _, oneClass := range contents {
    className := oneClass[1]
    if category == Common || category == Master {
      mappingLine := mappingTemplate
      mappingLine = strings.ReplaceAll(mappingLine, "$className", className)
      mappingLine = strings.Replace(mappingLine, "$category", categoryString, 1)
      mappingLine = strings.Replace(mappingLine, "$package", packageName, 1)
      mappingSb.WriteString(mappingLine)
    }

    attachChild(className, root, category)
    switch category {
    case Common:
      commonList = append(commonList, className)
    case Master:
      masterList = append(masterList, className)
    case Transaction:
      transactionList = append(transactionList, className)
    case ApiCommon:
      apiCommonList = append(apiCommonList, className)
    case Api:
      apiList = append(apiList, className)
    default:
      rich.ErrorThenThrow("Unkown class type: %v", category)
    }
  }
  return root
}

func attachChild(classPath string, parentTree *ProtoTree, category Category) {
  trimPrefix := getClassPath(parentTree, true)
  remnantClassPath := strings.TrimPrefix(classPath, trimPrefix)
  nameSlice := strings.Split(remnantClassPath, ".")
  currentTree := parentTree
  for _, name := range nameSlice {
    if _, ok := currentTree.children[name]; !ok {
      if currentTree.children == nil {
        currentTree.children = make(map[string]*ProtoTree)
      }
      level := currentTree.level + 1
      currentTree.children[name] = &ProtoTree{
        prefix:    getClassPath(currentTree, false),
        name:      name,
        level:     level,
        category:  category,
        traversed: false,
      }
    }
    currentTree = currentTree.children[name]
  }
}

func analyzeTree(
  entireContent *string,
  rootCategory Category,
  parentTree *ProtoTree,
  rootTree *ProtoTree,
  reducedLevel int,
) (*strings.Builder, *strings.Builder) {
  sb := new(strings.Builder)
  tsSb := new(strings.Builder)

  // to iterate map alphabetically
  keys := make([]string, 0, len(parentTree.children))
  for key := range parentTree.children {
    keys = append(keys, key)
  }
  sort.Strings(keys)

  for _, key := range keys {
    tree := parentTree.children[key]
    if tree.traversed {
      continue
    }
    classPath := getClassPath(tree, false)
    if classPath == "" {
      rich.ErrorThenThrow("Empty classPath.")
    }
    tsClassPath := strings.ReplaceAll(classPath, ".Types.", ".")
    tsClassPath = strings.ReplaceAll(tsClassPath, ".", "_")

    // if current tree name equals "Types" exactly, ignore it and reduce the tree level
    if tree.name == "Types" {
      reducedLevel++
    } else {
      sb.WriteString(indentMap[tree.level-1-reducedLevel] + strings.Replace(generalClassTemplate, "$className", tree.name, 1))
      if tree.level <= 1 {
        // exported class
        tsSb.WriteString(strings.Replace(tsGeneralClassTemplate, "$className", tree.name, 1))
      } else {
        // non-exported class
        tsSb.WriteString(strings.Replace(tsInnerClassTemplate, "$className", tsClassPath, 1))
      }
    }

    classSearchPtnStr := strings.Replace(nestedClassPtnStr, "$nestedClassName", classPath, 1)
    prefix := ""
    switch tree.category {
    case Common:
      prefix = strings.Replace(nestedClassPtnPrefixStr, "$which", "Common", 1)
    case Master:
      prefix = strings.Replace(nestedClassPtnPrefixStr, "$which", "Master", 1)
    case Transaction:
      prefix = strings.Replace(nestedClassPtnPrefixStr, "$which", "Transaction", 1)
    case ApiCommon:
      prefix = strings.Replace(nestedClassPtnPrefixStr, "$which", "Api.Common", 1)
    case Api:
      prefix = strings.Replace(nestedClassPtnPrefixStr, "$which", "Api", 1)
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
        tsLine := tsGeneralColumnTemplate
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
        tsTypeName := ""
        if mappedType, ok := typeMap[typeName]; ok {
          // in case of primitive types
          tsTypeName = tsTypeMap[typeName]
          typeName = mappedType
        } else {
          // in case of user defined types
          // first, check if it contains ".", if yes, it's highly likely a nested message
          if strings.Contains(typeName, ".") {
            if strings.HasPrefix(typeName, classPath+".") {
              attachChild(typeName, tree, Nested)
            } else {
              xList = append(xList, typeName)
            }
          } else {
            // if not, it can be an imported type
            if slices.Contains(enumList, typeName) {
              typeName = "penum." + typeName
            } else if slices.Contains(commonList, typeName) && rootCategory != Common {
              typeName = "pcommon." + typeName
            } else if slices.Contains(masterList, typeName) && rootCategory != Master {
              typeName = "pmaster." + typeName
            } else if slices.Contains(transactionList, typeName) && rootCategory != Transaction {
              typeName = "ptransaction." + typeName
            } else if slices.Contains(apiCommonList, typeName) && rootCategory != ApiCommon {
              typeName = "papicommon." + typeName
            } else if slices.Contains(apiList, typeName) && rootCategory != Api {
              typeName = "papi." + typeName
            }
          }
        }
        if tsTypeName == "" {
          tsTypeName = typeName
        }
        if isRepeated {
          typeName = "repeated " + typeName
          tsTypeName = tsTypeName + "[]"
        }
        columnName = strings.TrimSuffix(columnName, "_")

        // proto typeName
        if strings.Contains(typeName, ".Types.") {
          typeName = strings.ReplaceAll(typeName, ".Types.", ".")
        }
        line = strings.Replace(line, "$type", typeName, 1)
        line = strings.Replace(line, "$columnName", columnName, 1)
        line = strings.Replace(line, "$decimal", columnVal, 1)
        sb.WriteString(indentMap[tree.level-reducedLevel] + line)

        // typescript typeName
        if strings.Contains(tsTypeName, ".Types.") {
          tsTypeName = strings.ReplaceAll(tsTypeName, ".Types.", ".")
        }
        sections := strings.Split(tsTypeName, ".")
        if len(sections) > 1 {
          if slices.Contains([]string{
            "penum", "pcommon", "pmaster", "ptransaction", "papicommon", "papi",
          }, sections[0]) {
            tsTypeName = sections[0] + "." + strings.Join(sections[1:], "_")
          } else {
            tsTypeName = strings.Join(sections, "_")
          }
        } else {
          tsTypeName = strings.ReplaceAll(tsTypeName, ".", "_")
        }
        tsLine = strings.ReplaceAll(tsLine, "$type", tsTypeName)
        tsLine = strings.Replace(tsLine, "$columnName", columnName, 1)
        tsSb.WriteString("  " + tsLine)
      }
    }
    var tsXSbs []*strings.Builder
    nestedSb, tsNestedSb := analyzeTree(entireContent, rootCategory, tree, rootTree, reducedLevel)
    if len(xList) > 0 {
      for _, fullname := range xList {
        pfx, _ := getPrefixAndName(fullname)
        if pfx == classPath {
          attachChild(fullname, tree, Nested)
          xSb, tsXSb := analyzeTree(entireContent, rootCategory, tree, rootTree, reducedLevel)
          sb.WriteString(xSb.String())
          tsXSbs = append(tsXSbs, tsXSb)
        }
      }
    }
    sb.WriteString(nestedSb.String())
    if tree.name != "Types" {
      sb.WriteString(indentMap[tree.level-1-reducedLevel] + "}\n")
      tsSb.WriteString("}\n")
    }
    // write after parent class enclosed
    tsSb.WriteString(tsNestedSb.String())
    for _, xsb := range tsXSbs {
      tsSb.WriteString(xsb.String())
    }
    tree.traversed = true
  }
  return sb, tsSb
}

func analyzeFile(
  entireContent *string,
  category Category,
  outPath string,
  tsOutPath string,
) {
  // create a to be generated file
  protoFile, err := os.Create(outPath)
  if err != nil {
    panic(err)
  }
  tsFile, err := os.Create(tsOutPath)
  if err != nil {
    panic(err)
  }
  defer protoFile.Close()
  defer tsFile.Close()
  buf := bufio.NewWriter(protoFile)
  tsBuf := bufio.NewWriter(tsFile)
  // write prefixes
  var header string
  var tsHeader string
  switch category {
  case Enum:
    header = enumHeader
    tsHeader = tsEnumHeader
  case Common:
    header = commonHeader
    tsHeader = tsCommonHeader
  case Master:
    header = masterHeader
    tsHeader = tsMasterHeader
  case Transaction:
    header = transactionHeader
    tsHeader = tsTransactionHeader
  case ApiCommon:
    header = apiCommonHeader
    tsHeader = tsApiCommonHeader
  case Api:
    header = apiHeader
    tsHeader = tsApiHeader
  default:
    rich.ErrorThenThrow("Unkown type of proto file: %v.", category)
  }
  buf.WriteString(header)
  tsBuf.WriteString(tsHeader)

  var sb *strings.Builder
  var tsSb *strings.Builder
  if category == Enum {
    sb, tsSb = AnalyzeEnum(entireContent)
  } else {
    root := constructRoot(entireContent, category)
    sb, tsSb = analyzeTree(entireContent, category, root, root, 0)
  }

  buf.WriteString(sb.String())
  tsBuf.WriteString(tsSb.String())
  // flush
  if err := buf.Flush(); err != nil {
    panic(err)
  }
  if err := tsBuf.Flush(); err != nil {
    panic(err)
  }
}

func Analyze() {
  rich.Info("Start analyzing.")
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
  // create directory if not exists
  os.MkdirAll(outDir, 0755)
  os.MkdirAll(tsOutDir, 0755)

  entireContent := sb.String()
  mappingSb.WriteString(mappingHeader)
  analyzeFile(&entireContent, Enum, outEnumPath, tsOutEnumPath)
  rich.Info("Analyze Enum completed.")
  analyzeFile(&entireContent, Common, outCommonPath, tsOutCommonPath)
  rich.Info("Analyze Common completed.")
  analyzeFile(&entireContent, Master, outMasterPath, tsOutMasterPath)
  rich.Info("Analyze Master completed.")
  analyzeFile(&entireContent, Transaction, outTransactionPath, tsOutTransactionPath)
  rich.Info("Analyze Transaction completed.")
  analyzeFile(&entireContent, ApiCommon, outApiCommonPath, tsOutApiCommonPath)
  rich.Info("Analyze Api.Common completed.")
  analyzeFile(&entireContent, Api, outApiPath, tsOutApiPath)
  rich.Info("Analyze Api completed.")
  mappingSb.WriteString(mappingTail)
  if err := os.WriteFile(mappingOutFile, []byte(mappingSb.String()), 0644); err != nil {
    panic(err)
  }
  rich.Info("Analysis completed.")
}
