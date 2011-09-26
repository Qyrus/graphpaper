package graphpaper

import (
  "os"
  "io"
  "image"
  "math"
  "image/png"
  draw2d "draw2d.googlecode.com/hg/draw2d"
)

// DrawGraph draws a graph of the measurements as a PNG file and writes it to w.
func DrawGraph(w io.Writer, t Table) os.Error {

  var minDate int64
  var maxDate int64
  var minVal float64
  var maxVal float64

  for _, m := range t.Data() {
    if (minDate == 0 || m.NanoTime < minDate) { minDate = m.NanoTime }
    if (m.NanoTime > maxDate) { maxDate = m.NanoTime }
    f64 := float64(m.Values[0].Float64Value())
    minVal = math.Fmin(minVal, f64)
    maxVal = math.Fmax(maxVal, f64)
  }

  i := image.NewNRGBA(150, 100)
  gc := draw2d.NewGraphicContext(i)
  dx := 150 / (float64(maxDate) - float64(minDate))
  dy := 100 / (maxVal - minVal)

  gc.SetLineWidth(1)
  for _,m := range t.Data() {
    x := dx * float64(m.NanoTime - minDate)
    y := dy * (float64(m.Values[0].Float64Value()) - minVal)
    gc.MoveTo(x+0.45,y+0.5)
    gc.ArcTo(x+0.5,y+0.5,0.1,0.1,0,-math.Pi*2)
    gc.Stroke()
  }

  return png.Encode(w, i)
}

