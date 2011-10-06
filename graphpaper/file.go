package graphpaper

import (
  "os"
  "io"
  "bytes"
  "encoding/binary"
  "path/filepath"
)

/*
File format is relatively simple. It owes a lot to Artur Bergman's work on 
timesplicedb, but allows you to store values at either regular intervals or
arbitrary points in time.

40 bytes for a header:
  8 bytes "\x085grphppr"
  4 bytes for version number big endian
  4 bytes for type of measurements stored:
    1 - big endian unsigned 64 bit integer
    2 - big endian signed 64 bit integer
    3 - big endian double
  8 bytes for resolution in nanoseconds
  8 bytes for start time in nanoseconds since unix epoch
  8 bytes for column bit mask

Then for each data measurement:
  8 bytes for each column, in format specified by "type" and "function"

If resolution is zero and column bit mask is 0 then:
  start date is ignored
  8 bytes for date in nanoseconds since unix epoch
  8 bytes for value

*/

type fileHeader struct {
  Magic      uint64
  Version    int32
  ValueType  valueType
  Resolution int64
  StartTime  int64
  Functions  int64
}

const fileMagic uint64 = 0x8567727068616767

type File struct {
  *os.File
  *fileHeader
}

func (f *File) ColumnCount() int64 {
  // Kernighan's bit counting method
  c := int64(0)
  v := f.Functions
  for c = 0; v != 0; c++ {
    v = v & (v - 1)
  }
  return c
}

func (f *fileHeader) columnTypes() []valueType {
  // Todo: for now we assume coumns == 63. Fix that.
  switch f.ValueType {
  case uint64Type:
    return []valueType{int64Type, uint64Type, uint64Type, float64Type, float64Type, uint64Type}
  case int64Type:
    return []valueType{int64Type, int64Type, int64Type, float64Type, float64Type, int64Type}
  case float64Type:
    return []valueType{int64Type, float64Type, float64Type, float64Type, float64Type, float64Type}
  }
  panic("unexpected type")
}

func (f *fileHeader) Columns() []ColumnType {
  if f.IsRaw() {
    return []ColumnType{
      ColumnType{rawFunc, f.ValueType},
    }
  } else {
    // Todo: for now we assume functions == 63. Fix that.
    return []ColumnType{
      ColumnType{countFunc, int64Type},
      ColumnType{minFunc, f.ValueType},
      ColumnType{maxFunc, f.ValueType},
      ColumnType{sumFunc, float64Type},
      ColumnType{meanFunc, float64Type},
      ColumnType{medianFunc, f.ValueType},
    }
  }
  panic("unreachable code")
}

func (f *File) TimeOffset(time int64) int64 {
  return 40 + ((time-f.StartTime)/f.Resolution)*f.ColumnCount()*8
}

func (f *fileHeader) Bytes() (header []byte) {
  b := bytes.NewBuffer(header)
  binary.Write(b, binary.BigEndian, f)
  return b.Bytes()
}

func (h *fileHeader) IsRaw() bool {
  return (h.Resolution == 0 && h.StartTime == 0 && h.Functions == 0)
}

// CreateAggFile will open the specified file, returning an error if the file
// exists
func CreateFile(path string, valueType valueType, startTime int64, resolution int64, functions int64) (file *File, err os.Error) {
  // make directory
  dirpath, _ := filepath.Split(path)
  os.MkdirAll(dirpath, uint32(0755))

  // open file
  // todo: O_APPEND is a bad idea here?
  fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0666)
  if err != nil {
    return nil, err
  }

  header := fileHeader{fileMagic, 1, valueType, resolution, startTime, functions}
  _, err = fd.Write(header.Bytes())
  // todo: is this the right response to this error?
  if err != nil {
    return nil, err
  }

  return &File{fd, &header}, nil
}

// OpenAggFile will open the specified file, returning an error if the file
// does not exist.
func OpenFile(path string) (file *File, err os.Error) {

  fd, err := os.OpenFile(path, os.O_RDWR, 0666)
  if err != nil {
    return nil, err
  }

  var header fileHeader
  err = binary.Read(fd, binary.BigEndian, &header)
  if err != nil {
    return nil, err
  }
  if header.Magic != fileMagic || header.Version != 1 {
    return nil, os.NewError("Not a valid graphpaper file")
  }

  file = &File{fd, &header}
  return file, nil
}

func CreateOrOpenFile(path string, valueType valueType, startTime int64, resolution int64, functions int64) (file *File, err os.Error) {
  file, err = OpenFile(path)
  if e, ok := err.(*os.PathError); ok && e.Error == os.ENOENT {
    // todo: this has race conditions if run concurrently
    file, err = CreateFile(path, valueType, startTime, resolution, functions)
  }
  if err != nil {
    return nil, err
  }

  if file.ValueType != valueType && file.Resolution != resolution && file.StartTime != startTime && file.Functions != functions {
    return nil, os.NewError("File has unexpected values")
  }
  return file, nil

}

func (r *File) WriteSummary(s Summary) (err os.Error) {
  if r.Functions != s.Functions {
    return os.NewError("Incorrect functiondef")
  }

  if r.ValueType != s.ValueType {
    return os.NewError("Incorrect value type")
  }

  if r.Resolution != s.Resolution {
    return os.NewError("Incorrect interval")
  }

  for start, values := range s.Intervals {

    // todo: this doesn't actually need to be a buffer, could just be []byte
    b := bytes.NewBuffer([]byte{})
    for _, value := range values {
      b.Write(value.Bytes())
    }

    offset := r.TimeOffset(start)
    _, err := r.WriteAt(b.Bytes(), offset)
    // todo: check we've written the whole thing
    _ = err //todo: check error

  }
  return nil
}

func (r *File) appendRawMeasurements(l []Measurement) (e os.Error) {
  _, e = r.Seek(0, 2)
  if e != nil {
    return e
  }

  // todo: could be better
  b := []byte{}
  for _, m := range l {
    b = append(b, m.bytes()...)
  }
  _, e = r.Write(b)
  return e
}

func (r *File) ReadRawMeasurements() (l *MeasurementList, e os.Error) {
  // todo: should check if file is raw file
  byte_array := [16]byte{}
  ml := MeasurementList{}
  t := int64(0)
  for {
    _, err := io.ReadFull(r, byte_array[:])
    switch err {
    case os.EOF:
      return &ml, nil
    case nil:
      buf := bytes.NewBuffer(byte_array[:])
      binary.Read(buf, binary.BigEndian, &t)
      v, err := ReadValue(r.ValueType, buf)
      if err != nil {
        return nil, err
      }
      ml = append(ml, Measurement{v, t})
    default:
      return nil, err
    }
  }
  panic("unreachable code")
}

func (r *File) ReadAggregatedMeasurements() (s Summary, e os.Error) {
  s = Summary{r.ValueType, r.Functions, r.Resolution, map[int64][]Value{}}

  columns := r.columnTypes()

  for j := int64(0); ; j++ {
    point := make([]Value, len(columns))
    for i, column := range columns {
      v, err := ReadValue(column, r)
      if err == os.EOF {
        return s, nil
      }
      // todo: way more error checking here
      point[i] = v
    }
    s.Intervals[r.StartTime+(r.Resolution*j)] = point
  }
  panic("unreachable code")
}

func (r *File) ReadMeasurements() (Table, os.Error) {
  if r.Resolution == 0 && r.StartTime == 0 && r.Functions == 0 {
    list, err := r.ReadRawMeasurements()
    return list, err
  } else {
    summary, err := r.ReadAggregatedMeasurements()
    return summary, err
  }
  panic("unreachable code")
}
