package web

import (
  web "github.com/hoisie/web.go"
  "graphpaper"
  "time"
  "log"
)

func graph(ctx *web.Context, nodename string, property string) {

  m, err := graphpaper.GetMetric(nodename, property)
  if err != nil {
    log.Println("error: failed to get metric", err)
    ctx.Abort(500, "Error")
    return
  }

  table, err := m.GetMeasurements(time.Seconds() - 3600, time.Seconds())
  if err != nil {
    log.Println("error: failed to fetch metrics", err)
    ctx.Abort(500, "Error")
    return
  }

  ctx.SetHeader("Content-type", "image/png", true)
  err = graphpaper.DrawTable(ctx, table)
  if err != nil {
    log.Println("error: failed to draw graph", err)
    ctx.Abort(500, "Error")
  }
}
