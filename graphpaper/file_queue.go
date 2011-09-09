package graphpaper

import (
  "os"
  "fmt"
  "time"
)

// todo: better name please!
// todo: in the process of renaming, reconsider the interface between collectd and this queue
type NodeMetricMeasurement struct {
  Node string
  Metric string // todo: should be a type in itself?
  Measurement
}

type NodeMetricMeasurementChannel chan NodeMetricMeasurement

func (m *NodeMetricMeasurement) filename() (string) {
  // todo: sanity check node name and metric name before using them in the filesystem
  seconds := m.NanoTime / 1000000000
  fmt.Println(seconds, m.NanoTime)
  date := time.SecondsToUTC(seconds).Format("2006-01-02-15")
  return fmt.Sprintf("data/raw.1h/%s/%s/%s.gpr", date, m.Node, m.Metric)
}

// fileQueue is a naive implementation of a write queue for measurements that
// need to be written to disk.
// todo: should filequeue take filename / measurements pairs?
type fileQueue struct {
  data map[string] []Measurement
  queue []string
  waiting chan bool
  size int
}

// newFileQueue creates and initializes a new FileQueue. If data is waiting to
// be written then true will be sent to the waiting channel.
func newFileQueue(waiting chan bool) *fileQueue {
  q := new(fileQueue)
  // todo: simplify this
  q.data = make(map[string] []Measurement)
  q.queue = make([]string, 0)
  q.waiting = waiting
  return q
}

// setWaiting sends true to the waiting channel, unless the channel is blocked
func (q *fileQueue) setWaiting() {
  select {
  case q.waiting <- true:
  default:
  }
}

// push adds a measurement to the queue
func (q *fileQueue) push(m NodeMetricMeasurement) {
  // todo: queue should check and store types too
  filename := m.filename()
  _, exists := q.data[filename]
  q.size ++
  if exists {
    q.data[filename] = append(q.data[filename], m.Measurement)
  } else {
    q.data[filename] = []Measurement{}
    q.data[filename] = append(q.data[filename], m.Measurement)
    q.queue = append(q.queue, filename)
    q.setWaiting()
  }
}

// shift removes a list of measurements from the queue and returns them
func (q *fileQueue) shift() (string, []Measurement) {
  f := q.queue[0]
  q.queue = q.queue[1:]

  d := q.data[f]
  q.data[f] = []Measurement{}, false
  q.size = q.size - (len(d) / 16)
  if len(q.queue) > 0 {
    q.setWaiting()
  }
  return f, d
}

func writeMeasurement(filename string, l []Measurement) {
  if(len(l) > 0){
    first := l[0]
    valueType := first.Value.Type()
    fmt.Println("opening", valueType, filename)
    rawfile, err := CreateOrOpenFile(filename, valueType, 0, 0, 0)
    fmt.Println("done opening", valueType, filename)
    if err != nil { fmt.Println("Failed to open file", err); os.Exit(1); }
    defer rawfile.Close()

    err = rawfile.appendRawMeasurements(l)
    if err != nil { fmt.Println("Failed to write measurements", err); os.Exit(1); }
  }
}

// Filewriter recieves measurements from the channel and writes them to disk
// using a queue.
func FileWriter(c NodeMetricMeasurementChannel) {
  waiting := make(chan bool, 1)
  queue := newFileQueue(waiting)
  for {
    select {
    case m := <- c:
      queue.push(m)
    case <- waiting:
      filename, l := queue.shift()
      writeMeasurement(filename, l)
    }
  }
}
//*/