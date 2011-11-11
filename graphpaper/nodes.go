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

func (d LineDefinition) Equals(n LineDefinition) bool {
  return d.Node.Name == n.Node.Name && d.Property == n.Property && d.StatisticalFunction == n.StatisticalFunction && d.ValueType == n.ValueType 
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

func (t *DataTable) Append(n *DataTable) (r *DataTable, err os.Error) {
  if t.Resolution != n.Resolution {
    return nil, os.NewError("Resolutions don't match")
  }
  for i, d := range(t.LineDefinitions) {
    if !d.Equals(n.LineDefinitions[i]) {
      return nil, os.NewError("Line definitions don't match")
    }
  }
  // todo: should handle gaps between data and/or overlaps. For now we'll skip that, but we shouldn't.
  if t.Start >= n.End {
    return nil, os.NewError("DataTable date calculations are unimplemented")
  }
  data := append(t.Values, n.Values...)
  r = &DataTable{t.Start, n.End, t.Resolution, t.LineDefinitions, data}
  return r, nil
}

type Node struct {
  Name string
}

func (n Node) String() string {
  return n.Name
}

// todo: this should accept some notion of resolution
func (m Metric) GetMeasurements(start int64, end int64) (table *DataTable, err os.Error){

  tables := []*DataTable{}

  names := m.filenames(start * 1000000000, end * 1000000000)
  for _, name := range names {
    file, err := OpenFile(name)
    if err == nil {
      defer file.Close()
      fullTable, err := file.ReadAggregatedDataTable(m)
      if err == nil {
        // todo: don't hardcode a 1 here
        table := fullTable.Columnify(1)
        tables = append(tables, &table)
      }
    }
  }
  if len(tables) == 0 {
    return nil, err
  }
  r := tables[0]
  for i, t := range(tables) {
    if (i > 0) {
      r, err = t.Append(r)
      if err != nil {
        return nil, err
      }
    }
  }
  return r, nil
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

func (m Metric) filenames(start int64, end int64) ([]string) {
  // todo: should take resolution as an argument
  fc := Config.Resolutions[0]
  filenames := make([]string, 0)
  i := start - (start % fc.Size)
  for (i < end) {
    date := FormatTime(i, true, fc.DateFmt)
    filenames = append(filenames, fmt.Sprintf("data/%s/%s/%s/%s.gpr", fc.Name, date, m.Node, m.Property))
    i += fc.Size
  }
  return filenames
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
