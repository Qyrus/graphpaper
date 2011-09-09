package graphpaper

import (
  "fmt"
  "sort"
  "bytes"
  "encoding/binary"
)

type Measurement struct {
  Value    Value
  NanoTime int64
}

func (m *Measurement) bytes() []byte {
  b := new(bytes.Buffer)
  binary.Write(b, binary.BigEndian, m.NanoTime)
  b.Write(m.Value.Bytes())
  return b.Bytes()
}

type ColumnType struct {
  StatisticalFunction statisticalFunction
  ValueType           valueType
}

// a "measurement" is a single data point at a point in time
// the node and metric it is measured from to can be inferred by its context

type MeasurementList []Measurement

func (l MeasurementList) sortbucketize(interval int64) map[int64]ValueSlice {
  buckets := map[int64]ValueSlice{}

  r := map[int64]ValueSlice{}

  for _, m := range l {
    dateBucketStart := m.NanoTime - (m.NanoTime % interval)
    buckets[dateBucketStart] = append(buckets[dateBucketStart], m.Value)
  }

  for date, list := range buckets {
    sort.Sort(list)
    for _, i := range list {
      r[date] = append(r[date], i)
    }
  }
  return r
}

type TimeSplice struct {
  NanoTime int64
  Values   []Value
}

type TimeSpliceSlice []TimeSplice
// sort interface
func (s TimeSpliceSlice) Len() int {
  return len(s)
}
func (s TimeSpliceSlice) Less(i, j int) bool {
  return s[i].NanoTime < s[j].NanoTime
}
func (s TimeSpliceSlice) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}


func (l MeasurementList) Data() (d TimeSpliceSlice) {
  d = make(TimeSpliceSlice, len(l))
  for i, m := range l {
    d[i] = TimeSplice{m.NanoTime, []Value{m.Value}}
  }
  return d
}

func (l MeasurementList) Columns() []ColumnType {
  // Todo: fix the assumption that all measurements are the same type
  // Todo: fix assumption that a measurement list has at least one element
  return []ColumnType{
    ColumnType{rawFunc, l[0].Value.Type()},
  }
}

// a "function" is what we apply to measurements to get statistics - eg: mean, count, max, 75th percentile
// for performance reasons they're not actually implented as separate functions in the code
type statisticalFunction int32

const (
  rawFunc = iota
  countFunc
  minFunc
  maxFunc
  sumFunc
  meanFunc
  medianFunc
)

func (t statisticalFunction) String() string {
  switch t {
  case rawFunc:
    return "raw"
  case countFunc:
    return "count"
  case minFunc:
    return "min"
  case maxFunc:
    return "max"
  case sumFunc:
    return "sum"
  case meanFunc:
    return "mean"
  case medianFunc:
    return "median"
  }
  return fmt.Sprintf("unexpected type (%d)", t)
}

// a "summary" is a set of statistical functions applied to some measurements grouped by an interval
type Summary struct {
  ValueType  valueType
  Functions  int64
  Resolution int64
  Intervals  map[int64][]Value
}

type Table interface {
  Data() TimeSpliceSlice
  Columns() []ColumnType
}

func (s Summary) Data() (d TimeSpliceSlice) {
  d = make(TimeSpliceSlice, len(s.Intervals))
  i := 0
  for t, l := range s.Intervals {
    d[i] = TimeSplice{t, l}
    i++
  }
  return d
}

func (s Summary) Columns() []ColumnType {
  // Todo: for now we assume functions == 63. Fix that.
  return []ColumnType{
    ColumnType{countFunc, int64Type},
    ColumnType{minFunc, s.ValueType},
    ColumnType{maxFunc, s.ValueType},
    ColumnType{sumFunc, float64Type},
    ColumnType{meanFunc, float64Type},
    ColumnType{medianFunc, s.ValueType},
  }
}

func (s Summary) Bucketize(capacity int64) (l map[int64]Summary) {
  // todo: should check that capacity is divisible by interval
  l = map[int64]Summary{}
  for date, list := range s.Intervals {

    start := date - (date % capacity)
    _, exists := l[start]

    if !exists {
      l[start] = Summary{s.ValueType, s.Functions, s.Resolution, map[int64][]Value{}}
    }
    l[start].Intervals[date] = list
  }
  return l
}
