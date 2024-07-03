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

  tsEnumClassTemplate  = "export enum $className {\n"
  tsEnumColumnTemplate = "$columnName = $decimal,\n"
  tsEnumHeader         = "// Generated code. DO NOT EDIT!\n\n"
)

func AnalyzeEnum(entireContent *string) (*strings.Builder, *strings.Builder) {
  sb := new(strings.Builder)
  tsSb := new(strings.Builder)
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
    tsLine := strings.Replace(tsEnumClassTemplate, "$className", className, 1)
    tsSb.WriteString(tsLine)

    // each loop is a message
    for _, subMatches := range enumColumnPtn.FindAllStringSubmatch(content, -1) {
      line := enumColumnTemplate
      columnName := subMatches[1]
      columnVal := subMatches[2]
      line = strings.Replace(line, "$className", className, 1)
      line = strings.Replace(line, "$columnName", columnName, 1)
      line = strings.Replace(line, "$decimal", columnVal, 1)
      sb.WriteString("  " + line)

      tsLine := tsEnumColumnTemplate
      tsLine = strings.Replace(tsLine, "$columnName", columnName, 1)
      tsLine = strings.Replace(tsLine, "$decimal", columnVal, 1)
      tsSb.WriteString("  " + tsLine)
    }
    sb.WriteString("}\n")
    tsSb.WriteString("}\n")
  }
  return sb, tsSb
}
