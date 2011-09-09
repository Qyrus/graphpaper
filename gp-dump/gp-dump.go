package main

import (
  "graphpaper"
  "os"
  "fmt"
  goopt "github.com/droundy/goopt"
  "sort"
)

func main() {
  goopt.Summary = "graphpaper file dump"
  goopt.Parse(nil)

  f, err := graphpaper.OpenFile(goopt.Args[0])
  if err != nil { fmt.Println("Failed to open file", err); os.Exit(1); }

  table, err := f.ReadMeasurements()
  fmt.Println(table.Columns())

  fmt.Printf("%19.19s", "local time")
  for _, c := range table.Columns() {
    fmt.Printf("%10.10s ", c.StatisticalFunction)
  }
  fmt.Printf("\n")
  fmt.Printf("%19.19s", "")
  for _, c := range table.Columns() {
    fmt.Printf("%10.10s ", "("+c.ValueType.String()+")")
  }
  fmt.Printf("\n")

  data := table.Data()
  sort.Sort(data)

  for _, v := range data {

    fmt.Printf(graphpaper.FormatTime(v.NanoTime, false, "2006-01-02 15:04:05"))
    for i, c := range table.Columns() {
      switch c.ValueType {
      case 1:
        fmt.Printf("%10d ", v.Values[i])
      case 2:
        fmt.Printf("%10d ", v.Values[i])
      case 3:
        fmt.Printf("%10.3f ", v.Values[i])
      }
    }
    fmt.Printf("\n")

  }

}
