package rich

import (
  "fmt"
  "time"

  "github.com/fatih/color"
)

func getLocalTime() string {
  now := time.Now()
  return now.Format("2006-01-02 15:04:05.000")
}

func Info(text string, a ...any) {
  fmt.Println(color.BlueString(fmt.Sprintf(">>> [%s] [Info]", getLocalTime())), fmt.Sprintf(text, a...))
}

func Error(text string, a ...any) {
  fmt.Println(color.RedString(fmt.Sprintf(">>> [%s] [Error]", getLocalTime())), fmt.Sprintf(text, a...))
}

func Warning(text string, a ...any) {
  fmt.Println(color.YellowString(fmt.Sprintf(">>> [%s] [Warning]", getLocalTime())), fmt.Sprintf(text, a...))
}

func ErrorThenThrow(text string, a ...any) {
  Error(text, a...)
  panic("Exiting program due to the aforementioned reasons.")
}

func PanicError(err error, text string, a ...any) {
  Error(text, a...)
  panic(err)
}
