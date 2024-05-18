package analyser

import (
  "regexp"
  "strings"
)

var (
  enumClassPtn       = regexp.MustCompile(`// Namespace: Campus\.Common\.Proto\.Client\.Enums\npublic enum (?<className>\w+)[\s\S]+?\n}\n\n`)
  enumColumnPtn      = regexp.MustCompile(`public const \w+ (?<columnName>\w+) = (?<columnVal>\d+);`)
  enumClassTemplate  = "enum $className {\n"
  enumColumnTemplate = "$className_$columnName = $decimal;\n"
  enumHeader         = "syntax = \"proto3\";\npackage penum;\noption go_package = \"vertesan/campus/proto/penum\";\n\n"
)

func AnalyzeEnum(entireContent *string) *strings.Builder {
  sb := new(strings.Builder)
  // match all classes
  contents := enumClassPtn.FindAllStringSubmatch(*entireContent, -1)

  // each loop is a class
  for _, oneClass := range contents {
    content := oneClass[0]
    className := oneClass[1]

    // append a className into corresponding list
    enumList = append(enumList, className)

    // write class prefix
    line := strings.Replace(enumClassTemplate, "$className", className, 1)
    sb.WriteString(line)

    // each loop is a message
    for _, subMatches := range enumColumnPtn.FindAllStringSubmatch(content, -1) {
      line := enumColumnTemplate
      columnName := subMatches[1]
      columnVal := subMatches[2]
      line = strings.Replace(line, "$className", className, 1)
      line = strings.Replace(line, "$columnName", columnName, 1)
      line = strings.Replace(line, "$decimal", columnVal, 1)
      sb.WriteString("  " + line)
    }
    sb.WriteString("}\n")
  }
  return sb
}
