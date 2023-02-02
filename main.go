package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	g "github.com/AllenDang/giu"
	// glob "system64.net/KurumiGo/model/Globals"

	// op "system64.net/KurumiGo/model/Operator"
	"system64.net/KurumiGo/Kurumi"
	// synth "system64.net/KurumiGo/Kurumi/Kurumi"
)

var (
	showWindow2 bool
	checked     bool
)

func onShowWindow2() {
	showWindow2 = true
}

func onHideWindow2() {
	showWindow2 = false
}

var myBool = false
var myBoolPtr = &myBool

var myFloat = float32(0.0)
var myFloatPtr = &myFloat

var myLen int32 = 1
var myHei int32 = 1

var sineTable = make([]uint8, 0)
var Context *Kurumi.Synth

var myString = "Temporary string"

func loop() {
	g.MainMenuBar().Layout(
		g.Menu("File").Layout(
			g.MenuItem("Open"),
			g.Separator(),
			g.MenuItem("Save"),
			g.Menu("Export").Layout(
				g.MenuItem("Export WAV"),
				g.MenuItem("Export Seq. as WAV"),
			),
			g.Separator(),
			g.MenuItem("Load default patch"),
			g.MenuItem("Exit"),
		),
		g.Menu("Misc").Layout(
			g.Checkbox("Enable Me", &checked),
			g.Button("Button"),
		),
	).Build()

	// opLayout := g.Layout {
	// 	// g.Label(""),
	// 	g.TabBar().TabItems(
	// 		g.RangeBuilder("Operators", []interface{}{"Operator 1", "Operator 2", "Operator 3", "Operator 4"}, func (i int, v interface{}) g.Widget {
	// 			str := v.(string)
	// 			return g.TabItem(str),
	// 		}),
	// )
	// }

	g.Window("Operators").Pos(200, 200).Size(800, 400).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoResize).Layout(
		g.TabBar().TabItems(
			g.TabItem("Operator 1").Layout(
				g.Row(
					g.Column(
						g.Row(
							g.Label("TL volume"),
							g.SliderFloat(&Context.Operators[0].Tl, 0, 4).Size(256),
						),

						g.Row(
							g.Label("Feedback"),
							g.SliderFloat(&Context.Operators[0].Feedback, 0, 4).Size(256),
						),

						g.Row(
							g.Label("Mult"),
							g.SliderInt(&Context.Operators[0].Mult, 0, 32).Size(256),
						),
						g.Row(g.Label("")),
						g.Row(
							g.Checkbox("Phase modulation", &Context.Operators[0].PhaseMod),
						),
						g.Style().SetDisabled(!Context.Operators[0].PhaseMod).To(
							g.Column(
								g.Row(
									g.Label("Detune"),
									g.SliderInt(&Context.Operators[0].Detune, 1, 32).Size(256),
								),
								g.Row(
									g.Checkbox("Reverse phase", &Context.Operators[0].PhaseRev),
								),
								g.Row(
									g.Checkbox("Use custom phase envelope", &Context.Operators[0].CustomPhaseEnv),
								),
								g.Style().SetDisabled(Context.Operators[0].CustomPhaseEnv).To(
									g.Row(
										g.Label("Phase"),
										g.SliderFloat(&Context.Operators[0].Phase, 0, 1).Size(256),
									),
								),
							),
						),
					),
					g.Column(
						g.Label("Waveform :"),
						g.Combo("Waveform", Kurumi.Waveforms[Context.Operators[0].WaveformId], Kurumi.Waveforms, &Context.Operators[0].WaveformId).Size(256),
						g.Style().SetDisabled(Context.Operators[0].WaveformId < int32(len(Kurumi.Waveforms)-1)).To(
							g.Table().Size(256, 64).Rows(
								g.TableRow(
									g.Custom(
										func() {
											canvas := g.GetCanvas()
											pos := g.GetCursorScreenPos()
											color := color.RGBA{200, 75, 75, 255}
											canvas.AddLine(pos, pos.Add(image.Pt(256, 64)), color, 1)
											// for i := 0; i < len(sineTable); i++ {
											// 	sample := -int(sineTable[i]) + 128
											// 	canvas.AddRectFilled(pos.Add(image.Pt(int(math.Floor(float64(i)*512.0/float64(len(sineTable)))), 128)),
											// 		pos.Add(image.Pt(int(math.Ceil((float64(i)*512.0/float64(len(sineTable)))+(512.0/float64(len(sineTable))))), int(math.Abs(float64(-128-sample))))), color, 0, 0)
											// }

										},
									),
								).MinHeight(64).BgColor(color.RGBA{200, 75, 75, 0}),
							).InnerWidth(256),
							g.InputText(&myString).Size(256),
						),

						g.Style().SetDisabled(!Context.Operators[0].CustomPhaseEnv || !Context.Operators[0].PhaseMod).To(
							g.Label("Phase envelope :"),
							g.Table().Size(256, 64).Rows(
								g.TableRow(
									g.Custom(
										func() {
											canvas := g.GetCanvas()
											pos := g.GetCursorScreenPos()
											color := color.RGBA{200, 75, 75, 255}
											canvas.AddLine(pos, pos.Add(image.Pt(256, 64)), color, 1)
											// for i := 0; i < len(sineTable); i++ {
											// 	sample := -int(sineTable[i]) + 128
											// 	canvas.AddRectFilled(pos.Add(image.Pt(int(math.Floor(float64(i)*512.0/float64(len(sineTable)))), 128)),
											// 		pos.Add(image.Pt(int(math.Ceil((float64(i)*512.0/float64(len(sineTable)))+(512.0/float64(len(sineTable))))), int(math.Abs(float64(-128-sample))))), color, 0, 0)
											// }

										},
									),
								).MinHeight(64).BgColor(color.RGBA{200, 75, 75, 0}),
							).InnerWidth(256),
							g.InputText(&myString).Size(256),
						),
					),
				),
			),
			g.TabItem("Operator 2").Layout(),
			g.TabItem("Operator 3").Layout(),
			g.TabItem("Operator 4").Layout(),
		),
	)

	// ),

	g.Window("Modulation Matrix").Pos(30, 30).Size(320, 200).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoResize).Layout(
		g.TabBar().TabItems(
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
							g.Checkbox("", &Context.ModMatrix[0][0]),
							g.Tooltip("OP 1 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[0][1]),
							g.Tooltip("OP 2 modulates OP 1"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[0][2]),
							g.Tooltip("OP 3 modulates OP 1"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[0][3]),
							g.Tooltip("OP 4 modulates OP 1"),
						),
					),
					g.TableRow(
						g.Label("OP 2"),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][0]),
							g.Tooltip("OP 1 modulates OP 2"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][1]),
							g.Tooltip("OP 2 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][2]),
							g.Tooltip("OP 3 modulates OP 2"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][3]),
							g.Tooltip("OP 4 modulates OP 2"),
						),
					),
					g.TableRow(
						g.Label("OP 3"),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][0]),
							g.Tooltip("OP 1 modulates OP 3"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][1]),
							g.Tooltip("OP 2 modulates OP 3"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][2]),
							g.Tooltip("OP 3 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][3]),
							g.Tooltip("OP 4 modulates OP 3"),
						),
					),
					g.TableRow(
						g.Label("OP 4"),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][0]),
							g.Tooltip("OP 1 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][1]),
							g.Tooltip("OP 2 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][2]),
							g.Tooltip("OP 3 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][3]).OnChange(func() {
								println(*myBoolPtr)
							}),
							g.Tooltip("OP 4 modulates itself"),
						),
					),
				),
			),
			g.TabItem("Out. Levels").Layout(
				g.Tooltip("Here you can set the output of each operator to the exit."),
				g.Row(
					g.Label("OP 1 :"),
					g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				),
				g.Row(
					g.Label("OP 2 :"),
					g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				),
				g.Row(
					g.Label("OP 3 :"),
					g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				),
				g.Row(
					g.Label("OP 4 :"),
					g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				),
			),
			g.TabItem("Presets").Layout(
				g.ImageButton(nil).Size(100, 100),
			),
		),
		// g.Label("This is the modulation Matrix"),

	)

	g.Window("Wave preview").Flags(g.WindowFlagsNoResize|g.WindowFlagsNoCollapse).Pos(100, 100).Size(600, 600).Layout(
		g.Label("TEST PREVIEW WAVE"),

		g.Table().Size(520, 260).Rows(
			g.TableRow(
				g.Custom(
					func() {
						canvas := g.GetCanvas()
						pos := g.GetCursorScreenPos()
						color := color.RGBA{200, 75, 75, 255}
						for i := 0; i < len(sineTable); i++ {
							sample := -int(sineTable[i]) + 128
							canvas.AddRectFilled(pos.Add(image.Pt(int(math.Floor(float64(i)*512.0/float64(len(sineTable)))), 128)),
								pos.Add(image.Pt(int(math.Ceil((float64(i)*512.0/float64(len(sineTable)))+(512.0/float64(len(sineTable))))), int(math.Abs(float64(-128-sample))))), color, 0, 0)
						}

					},
				),
			).MinHeight(256).BgColor(color.RGBA{200, 75, 75, 0}),
		).InnerWidth(520),
		g.SliderInt(&myLen, 1, 256).Size(512),
		g.SliderInt(&myHei, 1, 255).Size(512),
	)

	// g.Window("Window 1").Pos(10, 30).Size(200, 100).Layout(
	// 	g.Label("I'm a label in window 1"),
	// 	g.Button("Show Window 2").OnClick(onShowWindow2),
	// 	g.Custom(func() {
	// 		canvas := g.GetCanvas()

	// 		pos := g.GetCursorScreenPos()
	// 		color := color.RGBA{200, 75, 75, 255}
	// 		canvas.AddLine(pos, pos.Add(image.Pt(100, 100)), color, 1)
	// 		canvas.AddRect(pos.Add(image.Pt(110, 0)), pos.Add(image.Pt(200, 100)), color, 5, g.DrawFlagsRoundCornersAll, 1)
	// 		canvas.AddRectFilled(pos.Add(image.Pt(220, 0)), pos.Add(image.Pt(320, 100)), color, 0, 0)

	// 		// canvas.AddTriangle(p1, p2, p3, color, 2)

	// 		// canvas.PathBezierCubicCurveTo(p2.Add(image.Pt(40, 0)), p3.Add(image.Pt(-50, 0)), p3, 0)
	// 		// canvas.PathStroke(color, false, 1)
	// 		// canvas.PathClear()
	// 	}),
	// )

	if showWindow2 {
		g.Window("Window 2").IsOpen(&showWindow2).Flags(g.WindowFlagsNone).Pos(250, 30).Size(200, 100).Layout(
			g.Label("I'm a label in window 2"),
			g.Button("Hide me").OnClick(onHideWindow2),
		)
	}
}

func main() {
	var len = 512
	for i := 0; i < len; i++ {
		sineTable = append(sineTable, uint8(math.Round((math.Sin((2*math.Pi*float64(i))/float64(len))+1)*127.5)))
	}
	Context = Kurumi.ConstructSynth()
	// glob.Context = Context
	fmt.Printf("%v", Context)
	wnd := g.NewMasterWindow("Kurumi 3 : The ultimate wavetable tool", 1280, 720, 0)
	wnd.Run(loop)
}
