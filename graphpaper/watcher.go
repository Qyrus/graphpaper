package graphpaper

import (
  "os"
  "fmt"
  "path/filepath"
  "time"
  "strings"
)

type visitor int64

func (v visitor) VisitDir(path string, f *os.FileInfo) bool {
    return true
}

func (v visitor) VisitFile(path string, f *os.FileInfo) {
  if (f.Mtime_ns > int64(v)) {
    parts := strings.Split(path, "/")
    file := parts[len(parts) - 1]
    dir := parts[len(parts) - 2]
    ext := filepath.Ext(file)

    if ext == ".gpr" {

      metric := file[:(len(file)-4)]
      node := dir

      rawfile, err := OpenFile(path)
      if err != nil { fmt.Println("Failed to open file", err); os.Exit(1); }
      defer rawfile.Close()

      list, err := rawfile.ReadRawMeasurements()
      s := Aggregate(list, 5 * 60 * 1000000000, 63)
      b := s.Bucketize(24 * 60 * 60 * 1000000000)

      for start, summary := range b {
        seconds := start / 1000000000
        
        date := time.SecondsToUTC(seconds).Format("2006-01-02")
        filename := fmt.Sprintf("data/5m.1d/%s/%s/%s.gpr", date, node, metric)
        file, err := CreateOrOpenFile(filename, summary.ValueType, start,  5 * 60 * 1000000000, summary.Functions)
        if err != nil { fmt.Println("Failed to open file", err); os.Exit(1); }
        defer file.Close()

        file.WriteSummary(s)

      }
    }

  }
}

// TODO: this doesn't check for any files that were written to before it starts
// up.
func Watch(dir string) {
  ticker := time.NewTicker(10 * 1000000000);
  last := int64(0)
  for {
    select {
    case t := <- ticker.C:
      if last > 0 {
        filepath.Walk(dir, visitor(last), nil)
      }
      last = t
    }
  }
}
