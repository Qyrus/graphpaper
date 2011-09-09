package main

import (
  "graphpaper"
)

func main() {

  mc := make(graphpaper.NodeMetricMeasurementChannel, 10)

  go graphpaper.CollectdListener(mc)
  go graphpaper.FileWriter(mc)
  graphpaper.Watch("./data/raw.1h")

  quit := make(chan bool)
  <-quit

}
