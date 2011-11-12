package web

import (
  web "github.com/hoisie/web.go"
  "log"
  "template"
)

func render(ctx *web.Context, templateName string, data interface{}) {
  err := templates.Execute(ctx, templateName, data)
  if err != nil {
    log.Println("error:", err)
    ctx.Abort(500, "Error parsing template")
  }
}

var templates = template.SetMust(template.ParseTemplateGlob("graphpaper/tmpl/*"))

// WebServer is the bare minimum web ui to show that the graphs are recoding data.
func Server() {
  web.Get("/", index)
  web.Get("/n/([^/]+)", node)
  web.Get("/n/([^/]+)/([^/]+).png", graph)

  web.Run("0.0.0.0:9999")
}
