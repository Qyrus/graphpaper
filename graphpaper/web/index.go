package web

import (
  web "github.com/hoisie/web.go"
  "graphpaper"
  "time"
  "log"
)

var indexTemplate = parseTemplate("graphpaper/tmpl/index.html")

func index(ctx *web.Context) {
  nodes, err := graphpaper.NodeList(time.Seconds())
  if err != nil {
    log.Println("error: failed to fetch nodes", err)
    ctx.Abort(500, "Error")
  } else {
    data := struct{ Nodes *[]graphpaper.Node }{nodes}
    indexTemplate.render(ctx, data)
  }
}
