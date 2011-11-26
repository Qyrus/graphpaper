package web

import (
  web "github.com/hoisie/web.go"
  "graphpaper"
  "log"
)

var metricTemplate = parseTemplate("graphpaper/tmpl/metric.html")

func metric(ctx *web.Context, nodename string, property string) {
  m, err := graphpaper.GetMetric(nodename, property)
  if err != nil {
    log.Println("error: failed to fetch metric", err)
    ctx.Abort(500, "Error")
  } else {
    metricTemplate.render(ctx, m)
  }
}
