package main

import (
  "graphpaper"
  goopt "github.com/droundy/goopt"
)


func main() {

  verbose := goopt.Flag([]string{"-v", "--verbose"}, []string{"--quiet"}, "output verbosely", "be quiet, instead")

  goopt.Summary = "The graphpaper metrics tool"
  goopt.Parse(nil)
  graphpaper.Debug = *verbose

  mc := make(graphpaper.NodeMetricMeasurementChannel, 10)

  go graphpaper.CollectdListener(mc)
  go graphpaper.FileWriter(mc)
  graphpaper.Watch("./data/raw.1h")

  quit := make(chan bool)
  <-quit

}
