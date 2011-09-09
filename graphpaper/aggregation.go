package graphpaper

import (
  "math"
)

func percentile(sorted []Value, percentile float64) Value {
  // todo: let's actually test this one
  i := int(math.Floor(percentile * float64(len(sorted))))
  return sorted[i]
}

func Aggregate(l *MeasurementList, interval int64, functions int64) (r Summary) {

  // todo: implement functions properly
  if functions != 63 {
    panic("unwritten code")
  }

  r.Functions = 63
  r.Resolution = interval
  r.ValueType = (*l)[0].Value.Type()
  r.Intervals = map[int64][]Value{}

  buckets := l.sortbucketize(interval)
  for dateStart, list := range buckets {
    count := int64Value(len(list))
    if count > 0 {
      min := list[0]
      max := list[count-1]
      sum := float64Value(0)
      for _, value := range list {
        // should be a method on value
        sum += value.Float64Value()
      }
      mean := sum / float64Value(count)
      // p25 := percentile(list, 0.25)
      median := percentile(list, 0.5)
      // p75 := percentile(list, 0.75)
      // p95 := percentile(list, 0.95)
      // p99 := percentile(list, 0.99)
      r.Intervals[dateStart] = []Value{count, min, max, sum, mean, median}
    }

  }
  return r
}
