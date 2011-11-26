package web

import (
	web "github.com/hoisie/web.go"
	"graphpaper"
	"time"
	"log"
	"strconv"
)

func graph(ctx *web.Context, nodename string, property string, width string, height string) {

	m, err := graphpaper.GetMetric(nodename, property)
	if err != nil {
		log.Println("error: failed to get metric", err)
		ctx.Abort(500, "Error")
		return
	}

	w, err := strconv.Atoui(width)
	if err != nil {
		log.Println("error: failed to convert height", err)
		ctx.Abort(500, "Error")
		return
	}

	h, err := strconv.Atoui(height)
	if err != nil {
		log.Println("error: failed to convert height", err)
		ctx.Abort(500, "Error")
		return
	}

	table, err := m.GetMeasurements(time.Seconds()-86400, time.Seconds())
	if err != nil {
		log.Println("error: failed to fetch metrics", err)
		ctx.Abort(500, "Error")
		return
	}

	ctx.SetHeader("Content-type", "image/png", true)
	err = graphpaper.DrawTable(ctx, table, w, h)
	if err != nil {
		log.Println("error: failed to draw graph", err)
		ctx.Abort(500, "Error")
	}
}
