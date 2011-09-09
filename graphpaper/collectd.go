package graphpaper

import(
  "net"
  "os"
  "fmt"
  "bytes"
  "encoding/binary"
)

// CollectdListener listens for data sent from collectd, parses it into
// CollectdMeasurements, and sends those measurements to the channel.
func CollectdListener(c NodeMetricMeasurementChannel) {
  laddr, err := net.ResolveUDPAddr("udp","127.0.0.1:25827");
  if err != nil { fmt.Println("Failed to resolve address", err); os.Exit(1); }
  conn, err := net.ListenUDP("udp", laddr);
  if err != nil { fmt.Println("Failed to listen", err); os.Exit(1); }
  for {
    buf := make([]byte, 1452)
    n, err := conn.Read(buf[:])
    if err != nil { fmt.Println("Failed to recieve packet", err); os.Exit(1); }
    collectdParse(c, buf[0:n])
  }
}

// generateName converts the various names collectd has for a metric into a single string.
// todo: make this not a compiled in hack
func generateName(plugin string, pluginInstance string, pluginType string, pluginTypeInstance string, index int) (name string) {
  switch {
  case plugin == "interface":
    switch {
    case index == 0:
      name = fmt.Sprintf("%s_%s_tx", pluginType, pluginTypeInstance)
    case index == 1:
      name = fmt.Sprintf("%s_%s_rx", pluginType, pluginTypeInstance)
    }
  case plugin == "df":
    switch {
    case index == 0:
      name = fmt.Sprintf("%s_%s_used", pluginType, pluginTypeInstance)
    case index == 1:
      name = fmt.Sprintf("%s_%s_free", pluginType, pluginTypeInstance)
    }
  case plugin == "load":
    switch {
    case index == 0:
      name = "load1"
    case index == 1:
      name = "load5"
    case index == 2:
      name = "load15"
    }
  case plugin == "memory":
    name = fmt.Sprintf("memory_%s", pluginTypeInstance)
  default:
    name = fmt.Sprintf("%s_%s_%s_%s", plugin, pluginInstance, pluginType, pluginTypeInstance)
  }
  return name
}

// collectdParse parses a packet sent using the collectd binary protocol and
// sends the resulting measurements to the channel.
func collectdParse(channel NodeMetricMeasurementChannel, b []byte) {
  buf := bytes.NewBuffer(b)
  hostname := ""
  plugin := ""
  pluginInstance := ""
  pluginType := ""
  pluginTypeInstance := ""
  var time int64

  for buf.Len() > 0 {
    var partType uint16
    var partLength uint16
    binary.Read(buf, binary.BigEndian, &partType)
    binary.Read(buf, binary.BigEndian, &partLength)
    partBytes := buf.Next(int(partLength) - 4)
    partBuffer := bytes.NewBuffer(partBytes)
    switch {
    case partType == 0:
      str := partBuffer.String()
      hostname = str[0:len(str)-1]
    case partType == 1:
      var timeSeconds int64
      binary.Read(partBuffer, binary.BigEndian, &timeSeconds)
      time = timeSeconds * 1000000000
    case partType == 2:
      str := partBuffer.String()
      plugin = str[0:len(str)-1]
    case partType == 3:
      str := partBuffer.String()
      pluginInstance = str[0:len(str)-1]
    case partType == 4:
      str := partBuffer.String()
      pluginType = str[0:len(str)-1]
    case partType == 5:
      str := partBuffer.String()
      pluginTypeInstance = str[0:len(str)-1]
    case partType == 6:
      var valueCount16 uint16
      binary.Read(partBuffer, binary.BigEndian, &valueCount16)
      valueCount := int(valueCount16)

      for i:=0; i < valueCount; i++ {

        name := generateName(plugin, pluginInstance, pluginType, pluginTypeInstance, i)
        
        // collectd's protocol puts things in a weird order.
        var valueType uint8
        binary.Read(partBuffer, binary.BigEndian, &valueType)
        valueBytes := make([]byte, 8, 8)
        copy(valueBytes, partBytes[2 + valueCount + (i*8):2 + valueCount + 8 + (i*8)])

        channel <- NodeMetricMeasurement{hostname, name, Measurement{ collectdValue{valueType, valueBytes}, time }}
        // messyness ends here
      }
    case partType == 7:
      // interval, ignore
    case partType == 8:
      // high res time
      // todo: get a copy of collectd 5 and test this
    case partType == 9:
      // interval, ignore
    case partType == 0x100:
      // message (notifications), ignore
    case partType == 0x100:
      // severity, ignore
    case partType == 0x200:
      // Signature (HMAC-SHA-256), todo
    case partType == 0x210:
      // Encryption (AES-256/OFB/SHA-1), todo
    default:
      // todo: log unexpected type here
    }
  }
}
