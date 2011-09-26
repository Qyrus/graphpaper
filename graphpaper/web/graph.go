package web

import (
  web "github.com/hoisie/web.go"
  "graphpaper"
  "time"
  "log"
)

func graph(ctx *web.Context, nodename string, property string) {
  m := graphpaper.Metric{graphpaper.Node{nodename}, graphpaper.Property(property)} // todo: should be a real method
  // todo: this is crying out to be wrapped up in a function
  file, err := m.File(time.Seconds())
  if err != nil {
    log.Println("error: failed to fetch metrics", err)
    ctx.Abort(500, "Error")
  } else {
    table, err := file.ReadMeasurements()
    if err != nil {
      log.Println("error: failed to fetch metrics", err)
      ctx.Abort(500, "Error")
    } else {
      ctx.SetHeader("Content-type", "image/png", true)
      err = graphpaper.DrawGraph(ctx, table)
      if err != nil {
        log.Println("error: failed to draw graph", err)
        ctx.Abort(500, "Error")
      }
    }
  }
}
