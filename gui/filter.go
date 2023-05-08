package gui

import (
	"image"
	"image/color"

	g "github.com/AllenDang/giu"
	"github.com/system64MC/Kurumi-Go/kurumi"
)

func drawFilterWindow() {
	g.Window("Filter").Pos(1000, 30).Layout(
		g.Checkbox("Enable filter", &kurumi.SynthContext.FilterEnabled).OnChange(func() {
			kurumi.Synthesize()
			kurumi.GenerateWaveStr()
		}),
		g.Tooltip("If enabled, the output will be filtered"),
		g.Style().SetDisabled(!kurumi.SynthContext.FilterEnabled).To(
			g.Combo("Filter Type", kurumi.FilterTypes[kurumi.SynthContext.FilterType], kurumi.FilterTypes, &kurumi.SynthContext.FilterType).OnChange(func() {
				kurumi.Synthesize()
				kurumi.GenerateWaveStr()
			}),
			g.Row(
				g.SliderFloat(&kurumi.SynthContext.Cutoff, 0, 1).OnChange(func() {
					// filterCutoff = Math.min(sample_rate/2,filterCutoff);calcBiquadFactors();});

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("Cutoff frequency"),
				g.Label("Cutoff"),
			),

			g.Row(
				g.SliderInt(&kurumi.SynthContext.Pitch, 0, 96).OnChange(func() {
					// filterCutoff = Math.min(sample_rate/2,filterCutoff);
					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("Simulates a particular pitch"),
				g.Label("Pitch"),
			),

			g.Row(
				g.SliderFloat(&kurumi.SynthContext.Resonance, 0.25, 4).OnChange(func() {
					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("Resonance"),
				g.Label("Q"),
			),

			g.Checkbox("Enable ADSR", &kurumi.SynthContext.FilterAdsrEnabled).OnChange(func() {
				kurumi.Synthesize()
				kurumi.GenerateWaveStr()
			}),
			g.Tooltip("If enabled, cutoff will be affected through time"),

			g.Style().SetDisabled(!kurumi.SynthContext.FilterAdsrEnabled).To(
				g.Row(
					g.SliderFloat(&kurumi.SynthContext.FilterStart, 0, 1).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Label("Start Cutoff"),
				),
				g.Row(
					g.SliderInt(&kurumi.SynthContext.FilterAttack, 0, 256).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Label("Attack"),
				),
				g.Row(
					g.SliderInt(&kurumi.SynthContext.FilterDecay, 0, 256).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Label("Decay"),
				),
				g.Row(
					g.SliderFloat(&kurumi.SynthContext.FilterSustain, 0, 1).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Label("Sustain"),
				),

				g.Label("ADSR :"),
				g.Style().SetStyle(g.StyleVarFramePadding, 0, 0).To(
					g.Table().Size(260, 68).Rows(
						g.TableRow(
							g.Custom(
								func() {
									canvas := g.GetCanvas()
									pos := g.GetCursorScreenPos()
									color := color.RGBA{75, 200, 75, 255}
									color2 := color
									color2.R = 25
									color2.G = 255
									color2.B = 25

									// For drawing ADSR
									{
										//adsr := kurumi.SynthContext.Operators[opId].Adsr
										tl := kurumi.SynthContext.Cutoff
										// Draw attack
										canvas.AddLine(pos.Add(image.Pt(0, int(64.0-kurumi.SynthContext.FilterStart*16))), pos.Add(image.Pt(int(kurumi.SynthContext.FilterAttack), int(64.0-tl*16))), color, 2)
										// canvas.AddLine(pos, pos.Add(image.Pt(256, 64)), color, 1)
										// Draw Decay
										canvas.AddLine(pos.Add(image.Pt(int(kurumi.SynthContext.FilterAttack), int(64.0-kurumi.SynthContext.Cutoff*16))), pos.Add(image.Pt(int(kurumi.SynthContext.FilterAttack+kurumi.SynthContext.FilterDecay), int(64.0-kurumi.SynthContext.FilterSustain*16))), color, 2)
										// Draw SUStain
										canvas.AddLine(pos.Add(image.Pt(int(kurumi.SynthContext.FilterAttack+kurumi.SynthContext.FilterDecay), int(64.0-kurumi.SynthContext.FilterSustain*16))), pos.Add(image.Pt(256, int(64.0-kurumi.SynthContext.FilterSustain*16))), color, 2)
										return
									}
								},
							),
						),
					),
				),
			),
		),
	)
}
