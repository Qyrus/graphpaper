package web

import (
  web "github.com/hoisie/web.go"
  "graphpaper"
  "time"
  "log"
)

func node(ctx *web.Context, nodename string) {
  n := graphpaper.Node{nodename} // todo: should be a real method
  metrics, err := graphpaper.MetricList(time.Seconds(), n)
  if err != nil {
    log.Println("error: failed to fetch metrics", err)
    ctx.Abort(500, "Error")
  } else {
    data := struct {
      graphpaper.Node
      Metrics *[]graphpaper.Metric
    }{n, metrics}
    render(ctx, "node", data)
  }
}
