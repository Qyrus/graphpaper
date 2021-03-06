package graphpaper

import (
	"os"
	"io"
	"image"
	"math"
	"image/png"
	draw2d "draw2d.googlecode.com/hg/draw2d"
)

func DrawTable(r io.Writer, t *DataTable, w uint, h uint) os.Error {

	var minVal float64
	var maxVal float64

	for _, v := range t.Values {
		if v != nil {
			f64 := float64(v.Float64Value())
			minVal = math.Fmin(minVal, f64)
			maxVal = math.Fmax(maxVal, f64)
		}
	}

	i := image.NewRGBA(int(w), int(h))
	gc := draw2d.NewGraphicContext(i)
	dx := float64(w-1) / (float64(t.End) - float64(t.Start))
	dy := float64(h-1) / (maxVal - minVal)

	gc.SetLineWidth(1)
	for i, v := range t.Values {
		if v == nil {
			continue
		}
		x := dx * float64(int64(i)*t.Resolution)
		y := dy * (float64(v.Float64Value()) - minVal)
		gc.MoveTo(x+0.45, y+0.5)
		gc.ArcTo(x+0.5, y+0.5, 0.1, 0.1, 0, -math.Pi*2)
		gc.Stroke()
	}

	return png.Encode(r, i)
}
