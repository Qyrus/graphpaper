package web

import (
  web "github.com/hoisie/web.go"
  mustache "github.com/hoisie/mustache.go"
  "path/filepath"
  "os"
  "fmt"
  "log"
)

func render(ctx *web.Context, templateName string, data interface{}) {
  // todo: validate template name?
  templatePath := fmt.Sprintf("templates/%s.mustache", templateName)
  templateFullPath := filepath.Join(os.Getenv("PWD"), templatePath)
  template, err := mustache.ParseFile(templateFullPath)

  if err != nil {
    log.Println("error:", err)
    ctx.Abort(500, "Error parsing template")
  } else {
    ctx.WriteString(template.Render(data))
  }
}

// WebServer is the bare minimum web ui to show that the graphs are recoding data.
func Server() {
  web.Get("/", index)
  web.Get("/n/([^/]+)", node)
  web.Get("/n/([^/]+)/([^/]+).png", graph)

  web.Run("0.0.0.0:9999")
}
