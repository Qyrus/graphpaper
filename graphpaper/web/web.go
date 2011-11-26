package web

import (
  web "github.com/hoisie/web.go"
  "log"
  tmpl "template"
)

type template struct {
  *tmpl.Set
}

func(t *template) render(ctx *web.Context, data interface{}){
  err := t.Set.Execute(ctx, "layout", data)
  if err != nil {
    log.Println("error:", err)
    ctx.Abort(500, "Error rendering template")
  }
}

func parseTemplate(file string) template {
  return template{tmpl.SetMust(tmpl.ParseSetFiles("graphpaper/tmpl/layout.html", file))}
}

// WebServer is the bare minimum web ui to show that the graphs are recoding data.
func Server() {

  web.Get("/", index)
  web.Get("/n/([^/]+)", node)
  web.Get("/n/([^/]+)/([^/]+).png", graph)

  web.Run("0.0.0.0:9999")
}
