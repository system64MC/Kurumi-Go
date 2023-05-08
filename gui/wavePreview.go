package gui

import (
	"image"
	"image/color"
	"math"

	g "github.com/AllenDang/giu"
	"github.com/system64MC/Kurumi-Go/kurumi"
)

func drawWavePreview() {
	g.Window("Wave preview").Flags(g.WindowFlagsNoResize|g.WindowFlagsNoCollapse).Pos(340, 20).Size(600, 350).Layout(
		g.Table().Size(520, 260).Rows(
			g.TableRow(
				g.Custom(
					func() {
						canvas := g.GetCanvas()
						pos := g.GetCursorScreenPos()
						color := color.RGBA{200, 75, 75, 255}
						wt := kurumi.WaveOutput

						if kurumi.SynthContext.WaveLen > 4096 {
							step := float64(kurumi.SynthContext.WaveLen) / 512.0
							len := float64(len(wt))
							for i := 0.0; i < len; i += step {
								x1 := int((float64(i) / step))
								x2 := int((float64(i) / step) + 1)
								sample := int(float64(wt[int(i)])*(255.0/float64(kurumi.SynthContext.WaveHei)) + (float64(kurumi.SynthContext.WaveHei)/2)*(255.0/float64(kurumi.SynthContext.WaveHei)))
								// sample := -int(math.Round(float64(wt[i]) * 255.0 / float64(kurumi.SynthContext.WaveHei)))
								canvas.AddRectFilled(pos.Add(image.Pt(x1, 128)),
									pos.Add(image.Pt(x2, int((float64(-sample+383))))), color, 0, 0)
							}
							return
						}

						for i := 0; i < len(wt); i++ {
							x1 := int(math.Floor(float64(i) * 512.0 / float64(len(wt))))
							x2 := int(math.Ceil((float64(i) * 512.0 / float64(len(wt))) + (512.0 / float64(len(wt)))))
							sample := int(float64(wt[i])*(255.0/float64(kurumi.SynthContext.WaveHei)) + (float64(kurumi.SynthContext.WaveHei)/2)*(255.0/float64(kurumi.SynthContext.WaveHei)))
							// sample := -int(math.Round(float64(wt[i]) * 255.0 / float64(kurumi.SynthContext.WaveHei)))
							canvas.AddRectFilled(pos.Add(image.Pt(x1, 128)),
								pos.Add(image.Pt(x2, int((float64(-sample+383))))), color, 0, 0)
						}

					},
				),
			).MinHeight(256).BgColor(color.RGBA{200, 75, 75, 0}),
		).InnerWidth(520),
		g.Row(
			g.SliderInt(&kurumi.SynthContext.WaveLen, 1, 256).Size(512).OnChange(func() {

				kurumi.Synthesize()
				kurumi.GenerateWaveStr()
			}),
			g.Tooltip("Changes length of wavetable"),
			g.Label("Length"),
		),
		g.Row(
			g.SliderInt(&kurumi.SynthContext.WaveHei, 1, 255).Size(512).OnChange(func() {

				kurumi.Synthesize()
				kurumi.GenerateWaveStr()
			}),
			g.Tooltip("Changes height of wavetable"),
			g.Label("Height"),
		),
	)
}
