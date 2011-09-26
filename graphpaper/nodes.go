package graphpaper

import (
  "os"
  "fmt"
  "time"
  "path/filepath"
)

type Node struct {
  Name string
}

func (n Node) String() string {
  return n.Name
}

type Property string

type Metric struct {
  Node
  Property
}

func (m Metric) File(t int64) (f *File, err os.Error) {
  // todo: this should take resolution as an argument
  date := time.SecondsToUTC(t).Format("2006-01-02-15")
  name := fmt.Sprintf("data/raw.1h/%s/%s/%s.gpr", date, m.Node, m.Property)
  return OpenFile(name)
}

func MetricList(t int64, n Node) (l *[]Metric, err os.Error) {
  // todo: dedupe this, move file path operations into shared code
  date := time.SecondsToUTC(t).Format("2006-01-02-15")
  glob := fmt.Sprintf("data/raw.1h/%s/%s/*.gpr", date, n.Name)
  // todo: notion of data dir?
  files, err := filepath.Glob(glob)
  if err != nil {
    return nil, err
  }
  metrics := make([]Metric, len(files))
  for i, f := range files {
    base := filepath.Base(f)
    property := base[:(len(base) - 4)]
    metrics[i] = Metric{n, Property(property)}
  }
  return &metrics, nil
}

func NodeList(t int64) (r *[]Node, err os.Error) {
  // todo: dedupe this, move file path operations into shared code
  date := time.SecondsToUTC(t).Format("2006-01-02-15")
  glob := fmt.Sprintf("data/raw.1h/%s/*", date)
  // todo: notion of data dir?
  files, err := filepath.Glob(glob)
  if err != nil {
    return nil, err
  }
  nodes := make([]Node, len(files))
  for i, s := range files {
    nodes[i] = Node{filepath.Base(s)}
  }
  return &nodes, nil
}
