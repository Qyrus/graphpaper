package graphpaper

import (
  "os"
  "fmt"
  "time"
  "path/filepath"
)

type DataTable struct {
  Start int64
  End int64
  Resolution int64
  LineDefinitions []LineDefinition
  Values []Value
}

type LineDefinition struct {
  Node
  Property
  StatisticalFunction statisticalFunction
  ValueType           valueType
}

func (d LineDefinition) String() string {
  return d.Node.Name + string(d.Property) + d.StatisticalFunction.String()
}

func (t DataTable) size() int {
  return len(t.Values) / len(t.LineDefinitions)
}

// todo: maybe this should take a named index not a numeric one?
func (t DataTable) Columnify(i int) (r DataTable) {
  definitions := t.LineDefinitions[i:i+1]
  defLen := len(t.LineDefinitions)
  data := make([]Value, t.size())
  for j, _ := range data {
    data[j] = t.Values[i + (j * defLen)]
  }
  return DataTable{t.Start, t.End, t.Resolution, definitions, data}
}

type Node struct {
  Name string
}

func (n Node) String() string {
  return n.Name
}

// todo: this should accept some notion of resolution
func (m Metric) GetMeasurements(start int64, end int64) (table DataTable, err os.Error){

  i := end
  for i > start {
    file, err := m.file(i)
    if err == nil {
      defer file.Close()
      summary, err := file.ReadAggregatedMeasurements()
      // todo: don't hardcode a 1 here
      table = summary.DataTable(m).Columnify(1)
      // todo: do something with err here
      _ = err
      i = (file.StartTime / 1000000000) - 1
    }
  }

  return table, err
}

type Property string

type Metric struct {
  Node
  Property
}

func GetMetric(n string, p string) (m Metric, err os.Error){
  // todo: this should check if the node actually has that metric
  return Metric{Node{n}, Property(p)}, nil
}

func (m Metric) file(t int64) (f *File, err os.Error) {
  // todo: this should take resolution as an argument
  date := time.SecondsToUTC(t).Format("2006-01-02")
  name := fmt.Sprintf("data/5m.1d/%s/%s/%s.gpr", date, m.Node, m.Property)
  return OpenFile(name)
}

func MetricList(t int64, n Node) (l *[]Metric, err os.Error) {
  // todo: dedupe this, move file path operations into shared code
  date := time.SecondsToUTC(t).Format("2006-01-02")
  glob := fmt.Sprintf("data/5m.1d/%s/%s/*.gpr", date, n.Name)
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
  date := time.SecondsToUTC(t).Format("2006-01-02")
  glob := fmt.Sprintf("data/5m.1d/%s/*", date)
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
