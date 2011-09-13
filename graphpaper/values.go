package graphpaper

import (
  "encoding/binary"
  "bytes"
  "io"
  "os"
  "fmt"
  "log"
)

type Value interface {
  Bytes() []byte
  Type() valueType
  Float64Value() float64Value
}

const (
  uint64Type  = 1
  int64Type   = 2
  float64Type = 3
)

type valueType int32

func (t valueType) String() string {
  switch t {
  case uint64Type:
    return "uint64"
  case int64Type:
    return "int64"
  case float64Type:
    return "float64"
  }
  return fmt.Sprintf("unexpected type (%d)", t)
}

func ReadValue(t valueType, f io.Reader) (v Value, err os.Error) {
  switch t {
  case uint64Type:
    var u64 uint64Value
    err = binary.Read(f, binary.BigEndian, &u64)
    return u64, err
  case int64Type:
    var i64 int64Value
    err = binary.Read(f, binary.BigEndian, &i64)
    return i64, err
  case float64Type:
    var f64 float64Value
    err = binary.Read(f, binary.BigEndian, &f64)
    return f64, err
  }
  return nil, os.NewError("unexpected type")
}

type collectdValue struct {
  collectdType uint8
  bytes        []byte
}

func (v collectdValue) Float64Value() float64Value {
  b := bytes.NewBuffer(v.Bytes())
  f, err := ReadValue(v.Type(), b)
  if err != nil {
    log.Fatalln("fatal: Failed to convert value", err)
  }
  return f.Float64Value()
}

func (v collectdValue) Bytes() []byte {
  b := v.bytes
  // float is little endian in collectd (why!?) and needs to be reversed
  if v.collectdType == 1 {
    for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
      b[i], b[j] = b[j], b[i]
    }
  }
  return b
}

func (v collectdValue) Type() (r valueType) {
  switch v.collectdType {
  case 0: // counter, unsigned int
    r = uint64Type
  case 1: // guage, float
    r = float64Type
  case 2: // derive, signed int
    r = int64Type
  case 3: // absolute, unsigned int
    r = uint64Type
  }
  return r
}

type uint64Value uint64

func (v uint64Value) Float64Value() float64Value {
  return float64Value(v)
}

func (v uint64Value) Bytes() []byte {
  b := new(bytes.Buffer)
  binary.Write(b, binary.BigEndian, v)
  return b.Bytes()
}

func (v uint64Value) Type() (r valueType) {
  return valueType(uint64Type)
}

type int64Value int64

func (v int64Value) Float64Value() float64Value {
  return float64Value(v)
}

func (v int64Value) Bytes() []byte {
  b := new(bytes.Buffer)
  binary.Write(b, binary.BigEndian, v)
  return b.Bytes()
}

func (v int64Value) Type() (r valueType) {
  return valueType(int64Type)
}

type float64Value float64

func (v float64Value) Float64Value() float64Value {
  return v
}

func (v float64Value) Bytes() []byte {
  b := new(bytes.Buffer)
  binary.Write(b, binary.BigEndian, v)
  return b.Bytes()
}

func (v float64Value) Type() (r valueType) {
  return valueType(float64Type)
}

type ValueSlice []Value

// sort interface
func (s ValueSlice) Len() int {
  return len(s)
}

func (s ValueSlice) Less(i, j int) bool {
  // todo: this would probably be more efficient (and accurate?) if we didn't use floats
  return s[i].Float64Value() < s[j].Float64Value()
}

func (s ValueSlice) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
