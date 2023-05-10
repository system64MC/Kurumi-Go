package gui

import (
	g "github.com/AllenDang/giu"
	"github.com/system64MC/Kurumi-Go/kurumi"
)

func drawMatrixWindow() {
	g.Window("Modulation Matrix").Pos(0, 20).Size(320, 200).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoResize).Layout(
		g.TabBar().Flags(g.TabBarFlagsFittingPolicyScroll).TabItems(
			g.TabItem("Matrix").Layout(
				g.Tooltip("Here you can create your own connections between operators"),
				g.Table().Size(200, 120).Flags(g.TableFlagsSizingStretchSame|g.TableFlagsBorders).Rows(
					g.TableRow(
						g.Label("Y\\X"),
						g.Column(
							g.Label("OP 1"),
						),
						g.Column(
							g.Label("OP 2"),
						),
						g.Column(
							g.Label("OP 3"),
						),
						g.Column(
							g.Label("OP 4"),
						),
					),

					g.TableRow(
						g.Label("OP 1"),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[0][0]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 1 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[0][1]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 2 modulates OP 1"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[0][2]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 3 modulates OP 1"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[0][3]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 4 modulates OP 1"),
						),
					),
					g.TableRow(
						g.Label("OP 2"),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[1][0]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 1 modulates OP 2"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[1][1]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 2 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[1][2]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 3 modulates OP 2"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[1][3]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 4 modulates OP 2"),
						),
					),
					g.TableRow(
						g.Label("OP 3"),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[2][0]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 1 modulates OP 3"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[2][1]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 2 modulates OP 3"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[2][2]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 3 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[2][3]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 4 modulates OP 3"),
						),
					),
					g.TableRow(
						g.Label("OP 4"),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[3][0]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 1 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[3][1]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 2 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[3][2]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 3 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &kurumi.SynthContext.ModMatrix[3][3]).OnChange(func() {

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.Tooltip("OP 4 modulates itself"),
						),
					),
				),
			),
			g.TabItem("Out. Levels").Layout(
				g.Tooltip("Here you can set the output of each operator to the exit."),
				g.RangeBuilder("OpOutSliders", []interface{}{"OP 1 :", "OP 2 :", "OP 3 :", "OP 4 :"}, func(i int, v interface{}) g.Widget {
					str := v.(string)
					// println(i)
					return g.Row(
						g.SliderFloat(&kurumi.SynthContext.OpOutputs[i], 0, 4).Size(256).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Label(str),
					)
				}),
			),

			g.TabItem("Presets").Layout(
				g.Tooltip("Presets of classic FM algorithms"),
				g.Column(
					g.RangeBuilder("rows", []interface{}{nil, nil, nil, nil, nil, nil}, func(i int, v interface{}) g.Widget {
						return g.Row(
							g.ImageButtonWithRgba(*algImages[i*2]).Size(128, 64).OnClick(func() {
								kurumi.ApplyAlg(i * 2)

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
							g.ImageButtonWithRgba(*algImages[(i*2)+1]).Size(128, 64).OnClick(func() {
								kurumi.ApplyAlg((i * 2) + 1)

								kurumi.Synthesize()
								kurumi.GenerateWaveStr()
							}),
						)
					}),
				),
			),
		),
	)
}
