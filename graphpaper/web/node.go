package web

import (
  web "github.com/hoisie/web.go"
  "graphpaper"
  "time"
  "log"
)

var nodeTemplate = parseTemplate("graphpaper/tmpl/node.html")

func node(ctx *web.Context, nodename string) {
  n := graphpaper.Node{nodename} // todo: should be a real method
  metrics, err := n.Metrics(time.Nanoseconds()-(86400*1000000000), time.Nanoseconds())
  if err != nil {
    log.Println("error: failed to fetch metrics", err)
    ctx.Abort(500, "Error")
  } else {
    data := struct {
      graphpaper.Node
      Metrics *[]graphpaper.Metric
    }{n, metrics}
    nodeTemplate.render(ctx, data)
  }
}
