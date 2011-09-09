package main

import (
  "graphpaper"
  "os"
  "fmt"
  goopt "github.com/droundy/goopt"
)

func main() {
  goopt.Summary = "graphpaper file info"
  goopt.Parse(nil)

  f, err := graphpaper.OpenFile(goopt.Args[0])
  if err != nil { fmt.Println("Failed to open file", err); os.Exit(1); }
  fmt.Printf("GPD Version: %d\n", f.Version)
  if f.IsRaw() {
    fmt.Printf("Type:        raw\n")
  } else {
    fmt.Printf("Type:        summary\n")
    fmt.Printf("Resolution:  %ds\n", f.Resolution / 1000000000)
    fmt.Printf("Start Time:  %s\n", graphpaper.FormatTimeLocal(f.StartTime))
  }
  fmt.Printf("Metric Type: %v\n", f.ValueType)
  columns := f.Columns()
  fmt.Printf("Columns:     %d %v\n", len(columns), columns)
  fmt.Println()
}