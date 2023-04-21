// go: generate goversioninfo -icon = kuruicon.ico
package main

import (
	// "fmt"
	"embed" 
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	

	rs "system64.net/KurumiGo/randomStuff"

	g "github.com/AllenDang/giu"
	c "github.com/atotto/clipboard"
	"github.com/ncruces/zenity"

	"system64.net/KurumiGo/kurumi"
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
var Context *kurumi.Synth

var myString = "Temporary string"
var operatorsGUIs = make([]*g.TabItemWidget, 0)
var waveStrs = []string{"255", "255", "255", "255"}
var morphStrs = []string{"255", "255", "255", "255"}
var envStrs = []string{"255", "255", "255", "255"}
var phaseStrs = []string{"0", "0", "0", "0"}
var comboBoxes = make([]*g.ComboWidget, 0)

func loop() {
	g.MainMenuBar().Layout(
		g.Menu("File").Layout(
			g.MenuItem("Open"),
			g.Separator(),
			g.MenuItem("Save"),
			g.Menu("Export").Layout(
				g.MenuItem("Furnace Wavetable").OnClick(func ()  {
					kurumi.CreateFUW()
				}),
				g.Menu("Dn-FamiTracker").Layout(
					g.MenuItem("Export FTI (N163)").OnClick(
						func() {
							kurumi.CreateFTIN163(false)
						},
					),
					g.MenuItem("Export FTI with sequence (N163)").OnClick(
						func() {
							kurumi.CreateFTIN163(true)
						},
					),
				),
				g.Menu("WAV").Layout(
					g.MenuItem("Export 16-Bits WAV").OnClick(
						func() {
							kurumi.SaveFile(false, true)
						},
					),
					g.MenuItem("Export Seq. as 16-Bits WAV").OnClick(
						func() {
							kurumi.SaveFile(true, true)
						},
					),
					g.MenuItem("Export 8-Bits WAV").OnClick(
						func() {
							kurumi.SaveFile(false, false)
						},
					),
					g.MenuItem("Export Seq. as 8-Bits WAV").OnClick(
						func() {
							kurumi.SaveFile(true, false)
						},
					),
				),
				g.Menu("TXT").Layout(
					g.MenuItem("Export TXT").OnClick(
						func() {
							kurumi.SaveTxt(false)
						},
					),
					g.MenuItem("Export Sequence as TXT").OnClick(
						func() {
							kurumi.SaveTxt(true)
						},
					),
				),
				g.Menu("RAW").Layout(
					g.MenuItem("Normalized 8-Bits RAW").OnClick(func() {
						kurumi.SaveRaw(false, 0)
					}),
					g.MenuItem("Non-Normalized 8-Bits RAW").OnClick(func() {
						kurumi.SaveRaw(false, 1)
					}),
					g.MenuItem("4-Bits RAW").OnClick(func() {
						kurumi.SaveRaw(false, 2)
					}),

					g.MenuItem("Sequence as Normalized 8-Bits RAW").OnClick(func() {
						kurumi.SaveRaw(true, 0)
					}),
					g.MenuItem("Sequence as Non-Normalized 8-Bits RAW").OnClick(func() {
						kurumi.SaveRaw(true, 1)
					}),
					g.MenuItem("Sequence as 4-Bits RAW").OnClick(func() {
						kurumi.SaveRaw(true, 2)
					}),
				),
				g.MenuItem("Deflemask DMW").OnClick(func() {
					openbrowser("https://github.com/tildearrow/furnace/releases/tag/v0.6pre4-hotfix")
				}),
				
			),
			g.Separator(),
			g.MenuItem("Load default patch").OnClick(func() {
				err := zenity.Warning("Are you sure you want to load default patch?\nAll unsaved data will be lost!!",
				zenity.Modal(), zenity.Title("WARNING!!"), zenity.ExtraButton("Yes, destroy everything!"), zenity.ExtraButton("NOOOOO!!!"))
				
				if(err == nil) {
					Context = kurumi.ConstructSynth()
					kurumi.SynthContext = Context
					initWaveStr()
					kurumi.Synthesize()
					kurumi.Synthesize()
					kurumi.InitAudio()
				}
			}),
			g.MenuItem("Exit"),
		),
		g.Menu("Misc").Layout(
			g.Checkbox("Enable Me", &checked),
			g.Button("Button"),
		),
		g.Checkbox("Enable sound", &Context.SongPlaying),
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

	g.Window("Operators").Pos(340, 370).Size(650, 400).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoResize).Layout(
		g.TabBar().TabItems(
			// g.TabItem("EEE").Layout(
			// 	g.Combo("", kurumi.Waveforms[Context.Operators[0].WaveformId], kurumi.Waveforms, &Context.Operators[0].WaveformId).Size(256),
			// ),
			buildOperator(0),
			buildOperator(1),
			buildOperator(2),
			buildOperator(3),
			// operatorsGUIs[0],
			// operatorsGUIs[1],
			// operatorsGUIs[2],
			// operatorsGUIs[3],
			// g.TabItem("Operator 2").Layout(),
			// g.TabItem("Operator 3").Layout(),
			// g.TabItem("Operator 4").Layout(),
		),
	)

	g.Window("Filter").Pos(1000, 30).Layout(
		g.Checkbox("Enable filter", &Context.FilterEnabled).OnChange(func() {
			kurumi.Synthesize()
		}),
		g.Tooltip("If enabled, the output will be filtered"),
		g.Style().SetDisabled(!Context.FilterEnabled).To(
			g.Combo("Filter Type", kurumi.FilterTypes[Context.FilterType], kurumi.FilterTypes, &Context.FilterType).OnChange(func() {
				kurumi.Synthesize()
			}),
			g.Row(
				g.SliderFloat(&Context.Cutoff, 0, 1).OnChange(func ()  {
					// filterCutoff = Math.min(sample_rate/2,filterCutoff);calcBiquadFactors();});
	
					kurumi.Synthesize()
				}),
				g.Tooltip("Cutoff frequency"),
				g.Label("Cutoff"),
			),
	
			g.Row(
				g.SliderInt(&Context.Pitch, 0, 96).OnChange(func ()  {
					// filterCutoff = Math.min(sample_rate/2,filterCutoff);
					kurumi.Synthesize()
				}),
				g.Tooltip("Simulates a particular pitch"),
				g.Label("Pitch"),

			),
	
			g.Row(
				g.SliderFloat(&Context.Resonance, 0.25, 4).OnChange(func ()  {
					kurumi.Synthesize()
				}),
				g.Tooltip("Resonance"),
				g.Label("Q"),

			),

			g.Checkbox("Enable ADSR", &Context.FilterAdsrEnabled).OnChange(func() {
				kurumi.Synthesize()
			}),
			g.Tooltip("If enabled, cutoff will be affected through time"),

			g.Style().SetDisabled(!Context.FilterAdsrEnabled).To(
				g.Row(
					g.SliderFloat(&Context.FilterStart, 0, 1).OnChange(func ()  {
						
						kurumi.Synthesize()
					}),
					g.Label("Start Cutoff"),
	
				),
				g.Row(
					g.SliderInt(&Context.FilterAttack, 0, 256).OnChange(func ()  {
						
						kurumi.Synthesize()
					}),
					g.Label("Attack"),
	
				),
				g.Row(
					g.SliderInt(&Context.FilterDecay, 0, 256).OnChange(func ()  {
						
						kurumi.Synthesize()
					}),
					g.Label("Decay"),
	
				),
				g.Row(
					g.SliderFloat(&Context.FilterSustain, 0, 1).OnChange(func ()  {
						
						kurumi.Synthesize()
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
											//adsr := Context.Operators[opId].Adsr
											tl := Context.Cutoff
											// Draw attack
											canvas.AddLine(pos.Add(image.Pt(0, int(64.0-Context.FilterStart*16))), pos.Add(image.Pt(int(Context.FilterAttack), int(64.0-tl*16))), color, 2)
											// canvas.AddLine(pos, pos.Add(image.Pt(256, 64)), color, 1)
											// Draw Decay
											canvas.AddLine(pos.Add(image.Pt(int(Context.FilterAttack), int(64.0-Context.Cutoff*16))), pos.Add(image.Pt(int(Context.FilterAttack+Context.FilterDecay), int(64.0-Context.FilterSustain*16))), color, 2)
											// Draw SUStain
											canvas.AddLine(pos.Add(image.Pt(int(Context.FilterAttack+Context.FilterDecay), int(64.0-Context.FilterSustain*16))), pos.Add(image.Pt(256, int(64.0-Context.FilterSustain*16))), color, 2)
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

	// ),

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
							g.Checkbox("", &Context.ModMatrix[0][0]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 1 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[0][1]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 2 modulates OP 1"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[0][2]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 3 modulates OP 1"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[0][3]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 4 modulates OP 1"),
						),
					),
					g.TableRow(
						g.Label("OP 2"),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][0]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 1 modulates OP 2"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][1]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 2 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][2]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 3 modulates OP 2"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[1][3]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 4 modulates OP 2"),
						),
					),
					g.TableRow(
						g.Label("OP 3"),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][0]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 1 modulates OP 3"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][1]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 2 modulates OP 3"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][2]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 3 modulates itself"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[2][3]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 4 modulates OP 3"),
						),
					),
					g.TableRow(
						g.Label("OP 4"),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][0]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 1 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][1]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 2 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][2]).OnChange(func() {
								
								kurumi.Synthesize()
							}),
							g.Tooltip("OP 3 modulates OP 4"),
						),
						g.Column(
							g.Checkbox("", &Context.ModMatrix[3][3]).OnChange(func() {
								
								kurumi.Synthesize()
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
						g.SliderFloat(&Context.OpOutputs[i], 0, 4).Size(256).OnChange(func() {
							
							kurumi.Synthesize()
						}),
						g.Label(str),

					)
				}),
				// g.Row(
				// 	g.Label("OP 1 :"),
				// 	g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				// ),
				// g.Row(
				// 	g.Label("OP 2 :"),
				// 	g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				// ),
				// g.Row(
				// 	g.Label("OP 3 :"),
				// 	g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				// ),
				// g.Row(
				// 	g.Label("OP 4 :"),
				// 	g.SliderFloat(myFloatPtr, 0, 4).Size(256),
				// ),
			),
			g.TabItem("Presets").Layout(
				g.Tooltip("Presets of classic FM algorithms"),
				g.Column(
					g.Row(
						g.ImageButtonWithRgba(*algImages[0]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(0)
							
							kurumi.Synthesize()
						}),
						g.ImageButtonWithRgba(*algImages[1]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(1)
							
							kurumi.Synthesize()
						}),
					),
					g.Row(
						g.ImageButtonWithRgba(*algImages[2]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(2)
							
							kurumi.Synthesize()
						}),
						g.ImageButtonWithRgba(*algImages[3]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(3)
							
							kurumi.Synthesize()
						}),
					),
					g.Row(
						g.ImageButtonWithRgba(*algImages[4]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(4)
							
							kurumi.Synthesize()
						}),
						g.ImageButtonWithRgba(*algImages[5]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(5)
							
							kurumi.Synthesize()
						}),
					),
					g.Row(
						g.ImageButtonWithRgba(*algImages[6]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(6)
							
							kurumi.Synthesize()
						}),
						g.ImageButtonWithRgba(*algImages[7]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(7)
							
							kurumi.Synthesize()
						}),
					),
					g.Row(
						g.ImageButtonWithRgba(*algImages[8]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(8)
							
							kurumi.Synthesize()
						}),
						g.ImageButtonWithRgba(*algImages[9]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(9)
							
							kurumi.Synthesize()
						}),
					),
					g.Row(
						g.ImageButtonWithRgba(*algImages[10]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(10)
							
							kurumi.Synthesize()
						}),
						g.ImageButtonWithRgba(*algImages[11]).Size(128, 64).OnClick(func() {
							kurumi.ApplyAlg(11)
							
							kurumi.Synthesize()
						}),
					),
				),
			),
		),
		// g.Label("This is the modulation Matrix"),

	)

	g.Window("Wave preview").Flags(g.WindowFlagsNoResize|g.WindowFlagsNoCollapse).Pos(340, 20).Size(600, 350).Layout(
		g.Table().Size(520, 260).Rows(
			g.TableRow(
				g.Custom(
					func() {
						canvas := g.GetCanvas()
						pos := g.GetCursorScreenPos()
						color := color.RGBA{200, 75, 75, 255}
						wt := kurumi.WaveOutput
						for i := 0; i < len(wt); i++ {
							x1 := int(math.Floor(float64(i) * 512.0 / float64(len(wt))))
							x2 := int(math.Ceil((float64(i) * 512.0 / float64(len(wt))) + (512.0 / float64(len(wt)))))
							sample := int(float64(wt[i])*(255.0/float64(Context.WaveHei)) + (float64(Context.WaveHei)/2)*(255.0/float64(Context.WaveHei)))
							// sample := -int(math.Round(float64(wt[i]) * 255.0 / float64(Context.WaveHei)))
							canvas.AddRectFilled(pos.Add(image.Pt(x1, 128)),
								pos.Add(image.Pt(x2, int((float64(-sample+383))))), color, 0, 0)
						}

					},
				),
			).MinHeight(256).BgColor(color.RGBA{200, 75, 75, 0}),
		).InnerWidth(520),
		g.Row(
			g.SliderInt(&Context.WaveLen, 1, 256).Size(512).OnChange(func() {
				
				kurumi.Synthesize()
			}),
			g.Tooltip("Changes length of wavetable"),
			g.Label("Length"),

		),
		g.Row(
			g.SliderInt(&Context.WaveHei, 1, 255).Size(512).OnChange(func() {
				
				kurumi.Synthesize()
			}),
			g.Tooltip("Changes height of wavetable"),
			g.Label("Height"),

		),
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

	g.Window("General settings").Size(340, 400).Pos(0, 220).Layout(
		g.Column(
			g.Checkbox("Normalize", &Context.Normalize).OnChange(func() {
				kurumi.Synthesize()
			}),
			g.Tooltip("Normalize the waveform"),
			g.Style().SetDisabled(Context.Normalize).To(
				g.Row(
					g.SliderFloat(&Context.Gain, 0, 4).OnChange(func() {
						
						kurumi.Synthesize()
					}),
					g.Tooltip("Amplifies the output"),
					g.Label("Gain"),
	
				),
			),
			g.Row(
				g.SliderInt(&Context.SmoothWin, 0, 128).OnChange(func() {
					
					kurumi.Synthesize()
				}),
				g.Tooltip("Smoothes the output"),
				g.Label("Avg. Filter Win."),

			),
			g.Row(
				g.SliderInt(&Context.MacLen, 1, 256).OnChange(func() {
					
					kurumi.Synthesize()
				}),
				g.Tooltip("How many frames the sequence has"),
				g.Label("Seq. Lenght"),

			),
			g.Row(
				g.SliderInt(&Context.Macro, 0, Context.MacLen-1).OnChange(func() {
					
					kurumi.Synthesize()
				}),
				g.Tooltip("The current sequence frame"),
				g.Label("Wav. Seq. Index"),

			),
			g.Row(
				g.SliderInt(&Context.Oversample, 1, 32).OnChange(func() {
					
					kurumi.Synthesize()
				}),
				g.Tooltip("Changes the oversample.\nOversample of 2x means everything is processed 2 times longer than the wavetable size,\nthen downsampled to its original size."),
				g.Label("Oversample"),

			),
			g.Label("Wave output :"),
			g.Row(
				g.InputText(kurumi.GenerateWaveStr()).Size(256).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAutoSelectAll),
				g.Button("Copy").OnClick(func() {
					c.WriteAll(*kurumi.GenerateWaveStr())
				}),
				g.Tooltip("Copy bitstring to clipboard"),
			),
			g.Label("Wave sequence output :"),
			g.InputTextMultiline(&kurumi.WaveSeqStr).Size(256, 128).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAutoSelectAll),
			g.Row(
				g.Button("Generate sequence").OnClick(func() {
					kurumi.GenerateWaveSeqStr()
				}),
				g.Tooltip("Generates the sequence"),
				g.Button("Copy").OnClick(func() {
					c.WriteAll(kurumi.WaveSeqStr)
				}),
				g.Tooltip("Copy bitstring sequence to clipboard"),
			),
		),
	)

	if showWindow2 {
		g.Window("Window 2").IsOpen(&showWindow2).Flags(g.WindowFlagsNone).Pos(250, 30).Size(200, 100).Layout(
			g.Label("I'm a label in window 2"),
			g.Button("Hide me").OnClick(onHideWindow2),
		)
	}
}

func buildOperator(a int) *g.TabItemWidget {

	// waveStrs[opId] = &kurumi.Waveforms[Context.Operators[opId].WaveformId]
	e := a
	str := "Operator " + strconv.Itoa(e+1)
	opId := e
	// out := fmt.Sprintf(str, i)
	// println(opId)
	return g.TabItem(str).Layout(
		g.ContextMenu().Layout(
			g.Button("Copy operator").OnClick(func() {
				// println("Copying operator" + strconv.Itoa(a))
				kurumi.CopyOp(&Context.Operators[a])
				g.CloseCurrentPopup()
			}),
			g.Button("Paste operator").OnClick(func() {
				kurumi.PasteOp(&Context.Operators[a])
				initWaveStr()
				// println("Pasting operator" + strconv.Itoa(a))
				g.CloseCurrentPopup()
			}),
		),
		g.Row(
			g.Column(
				g.Row(
					g.Label("Mod. Mode"),
					g.Combo("", kurumi.ModModes[Context.Operators[opId].ModMode], kurumi.ModModes, &Context.Operators[opId].ModMode).Size(128).OnChange(func() {
						
						kurumi.Synthesize()
					}),
					g.Tooltip("Changes the modulation mode.\nFor example, if it is set to MUL, the operator will take the input and then multiply it with current oscillator state.\nDefault is FM"),
				),
				g.Row(g.Label("")),

				g.Row(
					g.SliderFloat(&Context.Operators[opId].Tl, 0, 4).Size(256).OnChange(func() {
						
						kurumi.Synthesize()
					}),
					g.Tooltip("Changes Mod. Depth / Volume"),
					g.Label("Mod. Depth"),

				),
				g.Row(
					g.SliderFloat(&Context.Operators[opId].Feedback, 0, 4).Size(256).OnChange(func() {
						
						kurumi.Synthesize()
					}),
					g.Tooltip("Changes the Feedback of operator"),
					g.Label("Feedback"),

				),
				g.Row(
					g.SliderInt(&Context.Operators[opId].Mult, 1, 32).Size(256).OnChange(func() {
						
						kurumi.Synthesize()
					}),
					g.Tooltip("Changes the frequency multiplier"),
					g.Label("Mult"),

				),
				g.Row(
					g.SliderFloat(&Context.Operators[opId].Phase, 0, 1).Size(256).OnChange(func() {
						
						kurumi.Synthesize()
					}),
					g.Tooltip("Changes the phase of the waveform"),
					g.Label("Phase"),

				),

				g.Row(g.Label("")),

				g.Checkbox("Use envelope on volume", &Context.Operators[opId].IsEnvelopeEnabled).OnChange(func() {
					
					kurumi.Synthesize()
				}),
				g.Tooltip("If enabled, operator's volume / Mod. depth will be affected by envelope"),

				g.Style().SetDisabled(!Context.Operators[opId].IsEnvelopeEnabled).To(
					g.Row(
						g.Checkbox("Use custom volume envelope instead of ADSR", &Context.Operators[opId].UseCustomVolEnv).OnChange(func() {
							
							kurumi.Synthesize()
						}),
					),
					g.Tooltip("Disables ADSR and enables user-defined envelope"),
					g.Style().SetDisabled(Context.Operators[opId].UseCustomVolEnv).To(
						g.Column(
							g.Row(
								g.SliderInt(&Context.Operators[opId].Adsr.Attack, 0, 256).Size(256).OnChange(func() {
									
									kurumi.Synthesize()
								}),
								g.Label("Attack"),
	
							),
							g.Row(
								g.SliderInt(&Context.Operators[opId].Adsr.Decay, 0, 256).Size(256).OnChange(func() {
									
									kurumi.Synthesize()
								}),
								g.Label("Decay"),
	
							),
							g.Row(
								g.SliderFloat(&Context.Operators[opId].Adsr.Sustain, 0, 4).Size(256).OnChange(func() {
									
									kurumi.Synthesize()
								}),
								g.Label("Sustain"),
	
							),
						),
					),
				),

				
				g.Row(g.Label("")),
				// g.Row(
				// 	g.Checkbox("Phase modulation", &Context.Operators[opId].PhaseMod).OnChange(func() {
				// 		kurumi.ResetFB()
				// 		kurumi.Synthesize()
				// 		kurumi.Synthesize()
				// 	}),
				// ),
				g.Column(
					g.Row(
						g.SliderInt(&Context.Operators[opId].Detune, -32, 32).Size(256).OnChange(func() {
							
							kurumi.Synthesize()
						}),
						g.Tooltip("If not zero, the phase will change according to sequence"),
						g.Label("Detune"),

					),
					// g.Row(
					// 	g.Checkbox("Reverse phase", &Context.Operators[opId].PhaseRev).OnChange(func() {
					// 		kurumi.ResetFB()
					// 		kurumi.Synthesize()
					// 		kurumi.Synthesize()
					// 	}),
					// ),
					g.Row(
						g.Checkbox("Use custom phase envelope", &Context.Operators[opId].CustomPhaseEnv).OnChange(func() {
							
							kurumi.Synthesize()
						}),
						g.Tooltip("If enabled, you can control the phase with a custom envelope.\nIt can be used to imitate vibrato."),
					),
				),
				g.Row(g.Label("")),
				g.Checkbox("Enable wavetable morphing", &Context.Operators[opId].Morphing).OnChange(func() {
					
					kurumi.Synthesize()
				}),
				g.Tooltip("If enabled, the custom waveform will be morphed to another one"),
				g.Style().SetDisabled(!Context.Operators[opId].Morphing).To(
					g.Row(
						g.SliderInt(&Context.Operators[opId].MorphTime, 1, 256).Size(256).OnChange(func() {
							
							kurumi.Synthesize()
						}),
						g.Tooltip("How long the morphing takes"),
						g.Label("Morph Time"),

					),
				),

				g.Row(
					g.SliderFloat(&Context.Operators[opId].DutyCycle, 0, 1).Size(256).OnChange(func ()  {
						
						kurumi.Synthesize()
					}),
					g.Tooltip("Changes the duty cycle of the Pulse waveform.\nDefault is 0.5 (50%)"),
					g.Label("Duty Cycle"),
	
				),

				g.Checkbox("Enable PWM ADSR", &Context.Operators[opId].PwmAdsrEnabled).OnChange(func() {
					
					kurumi.Synthesize()
				}),
				g.Tooltip("If enabled, the pulse width will change through time."),

				g.Style().SetDisabled(!Context.Operators[opId].PwmAdsrEnabled).To(
					g.Row(
						g.SliderFloat(&Context.Operators[opId].PwmAdsr.Start, 0, 1).Size(256).OnChange(func ()  {
							
							kurumi.Synthesize()
						}),
						g.Label("Start Duty"),
		
					),
					g.Row(
						g.SliderInt(&Context.Operators[opId].PwmAdsr.Attack, 0, 256).Size(256).OnChange(func ()  {
							
							kurumi.Synthesize()
						}),
						g.Label("Attack"),
		
					),
					g.Row(
						g.SliderInt(&Context.Operators[opId].PwmAdsr.Decay, 0, 256).Size(256).OnChange(func ()  {
							
							kurumi.Synthesize()
						}),
						g.Label("Decay"),
		
					),
					g.Row(
						g.SliderFloat(&Context.Operators[opId].PwmAdsr.Sustain, 0, 1).Size(256).OnChange(func ()  {
							
							kurumi.Synthesize()
						}),
						g.Label("Final Duty"),
					),
				),

				
			),
			g.Column(

				g.Row(
					g.Label("Waveform :"),
					// g.Custom(func() {
					// 	combo := g.Combo(kurumi.Waveforms[Context.Operators[opId].WaveformId], kurumi.Waveforms[Context.Operators[opId].WaveformId], kurumi.Waveforms, &Context.Operators[opId].WaveformId).Size(256)
					// 	combo = combo
					// }),
					g.Combo("", kurumi.Waveforms[Context.Operators[opId].WaveformId], kurumi.Waveforms, &Context.Operators[opId].WaveformId).Size(170).OnChange(func() {
						
						kurumi.Synthesize()
					}),
				),
				g.Checkbox("Reverse waveform", &Context.Operators[opId].Reverse).OnChange(func() {
					
						kurumi.Synthesize()
				}),
				g.Tooltip("Flip down the waveform"),

				g.Table().Size(260, 68).Rows(
					g.TableRow(
						g.Custom(
							func() {
								canvas := g.GetCanvas()
								pos := g.GetCursorScreenPos()
								color := color.RGBA{75, 75, 255, 255}
								for i := 0; i < 256; i++ {
									y := 0.0
									if(Context.Operators[opId].Reverse) {
										y = -kurumi.WaveFuncs[Context.Operators[opId].WaveformId](&Context.Operators[opId], float64(i)/256.0)*-32 + 1

									} else {
										y = kurumi.WaveFuncs[Context.Operators[opId].WaveformId](&Context.Operators[opId], float64(i)/256.0)*-32 + 1

									}
									canvas.AddRectFilled(pos.Add(image.Pt(int(math.Floor(float64(i))), 32)),
										pos.Add(image.Pt(int(math.Ceil((float64(i))+1)), int(math.Abs(float64(-32-y))))), color, 0, 0)
								}

							},
						),
					).MinHeight(64).BgColor(color.RGBA{200, 75, 75, 0}),
				).InnerWidth(256),
				g.Style().SetDisabled(Context.Operators[opId].WaveformId < int32(len(kurumi.Waveforms)-5)).To(
					g.Row(
						g.Label("Interpolation :"),
						g.Combo("", kurumi.Interpolations[Context.Operators[opId].Interpolation], kurumi.Interpolations, &Context.Operators[opId].Interpolation).Size(160).OnChange(func() {
							
							kurumi.Synthesize()
						}),
					),
					g.InputText(&waveStrs[opId]).Size(256).OnChange(func() {
						// println(waveStrs[opId])
						kurumi.ApplyStringToWaveform(opId, waveStrs[opId], kurumi.DestWave)
						
						kurumi.Synthesize()
					}),
					g.Style().SetDisabled(!Context.Operators[opId].Morphing).To(
						g.Label("Waveform to morph to :"),
						g.InputText(&morphStrs[opId]).Size(256).OnChange(func() {
						kurumi.ApplyStringToWaveform(opId, envStrs[opId], kurumi.DestMorph)
							
							kurumi.Synthesize()
						}),
					),
				),

				g.Style().SetDisabled(!Context.Operators[opId].IsEnvelopeEnabled).To(
					g.Label("ADSR / Volume envelope :"),
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
										if !Context.Operators[opId].UseCustomVolEnv {
											adsr := Context.Operators[opId].Adsr
											tl := Context.Operators[opId].Tl
											// Draw attack
											canvas.AddLine(pos.Add(image.Pt(0, 64)), pos.Add(image.Pt(int(adsr.Attack), int(64.0-tl*16))), color, 2)
											// canvas.AddLine(pos, pos.Add(image.Pt(256, 64)), color, 1)
											// Draw Decay
											canvas.AddLine(pos.Add(image.Pt(int(adsr.Attack), int(64.0-tl*16))), pos.Add(image.Pt(int(adsr.Attack+adsr.Decay), int(64.0-adsr.Sustain*16))), color, 2)
											// Draw SUStain
											canvas.AddLine(pos.Add(image.Pt(int(adsr.Attack+adsr.Decay), int(64.0-adsr.Sustain*16))), pos.Add(image.Pt(256, int(64.0-adsr.Sustain*16))), color, 2)
											return
										}
										// For drawing custom envelope
										env := Context.Operators[opId].VolEnv
										if len(env) <= 0 {
											return
										}
										for x := 0; x < 256; x++ {
											s1 := float64(env[int(math.Min(float64(x), float64(len(env))-1))]) / 4.0
											s2 := float64(env[int(math.Min(float64(x+1), float64(len(env))-1))]) / 4.0
											canvas.AddLine(pos.Add(image.Pt(x, int(64-s1))), pos.Add(image.Pt(x+1, int(64-s2))), color2, 2)
										}
									},
								),
							),
						),
					),
				),

				

				g.Style().SetDisabled(!Context.Operators[opId].UseCustomVolEnv).To(
					g.InputText(&envStrs[opId]).Size(256).OnChange(func() {
						kurumi.ApplyStringToWaveform(opId, envStrs[opId], kurumi.DestVolEnv)
						
						kurumi.Synthesize()
					}),
				),

				g.Style().SetDisabled(!Context.Operators[opId].CustomPhaseEnv || !Context.Operators[opId].PhaseMod).To(
					g.Label("Phase envelope :"),
					g.Table().Size(260, 68).Rows(
						g.TableRow(
							g.Custom(
								func() {
									canvas := g.GetCanvas()
									pos := g.GetCursorScreenPos()
									color := color.RGBA{200, 75, 75, 255}

									env := Context.Operators[opId].PhaseEnv
									if len(env) <= 0 {
										return
									}
									for x := 0; x < 256; x++ {
										s1 := float64(env[int(math.Min(float64(x), float64(len(env))-1))]) / 4.0
										s2 := float64(env[int(math.Min(float64(x+1), float64(len(env))-1))]) / 4.0
										canvas.AddLine(pos.Add(image.Pt(x, int(64-s1))), pos.Add(image.Pt(x+1, int(64-s2))), color, 2)
									}

									// canvas.AddLine(pos, pos.Add(image.Pt(256, 64)), color, 1)
									// for i := 0; i < len(sineTable); i++ {
									// 	sample := -int(sineTable[i]) + 128
									// 	canvas.AddRectFilled(pos.Add(image.Pt(int(math.Floor(float64(i)*512.0/float64(len(sineTable)))), 128)),
									// 		pos.Add(image.Pt(int(math.Ceil((float64(i)*512.0/float64(len(sineTable)))+(512.0/float64(len(sineTable))))), int(math.Abs(float64(-128-sample))))), color, 0, 0)
									// }

								},
							),
						).MinHeight(64).BgColor(color.RGBA{200, 75, 75, 0}),
					),
					g.InputText(&phaseStrs[opId]).Size(256).OnChange(func() {
						kurumi.ApplyStringToWaveform(opId, phaseStrs[opId], kurumi.DestPhaseEnv)
						
						kurumi.Synthesize()
					}),
				),


				g.Style().SetDisabled(!Context.Operators[opId].PwmAdsrEnabled).To(
					g.Label("PWM envelope :"),
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

										pwmAdsr := Context.Operators[opId].PwmAdsr
	
										// For drawing ADSR
										{
											//adsr := Context.Operators[opId].Adsr
											tl := Context.Operators[opId].DutyCycle
											// Draw attack
											canvas.AddLine(pos.Add(image.Pt(0, int(64.0-pwmAdsr.Start*64))), pos.Add(image.Pt(int(pwmAdsr.Attack), int(64.0-tl*64))), color, 2)
											// canvas.AddLine(pos, pos.Add(image.Pt(256, 64)), color, 1)
											// Draw Decay
											canvas.AddLine(pos.Add(image.Pt(int(pwmAdsr.Attack), int(64.0-tl*64))), pos.Add(image.Pt(int(pwmAdsr.Attack+pwmAdsr.Decay), int(64.0-pwmAdsr.Sustain*64))), color, 2)
											// Draw SUStain
											canvas.AddLine(pos.Add(image.Pt(int(pwmAdsr.Attack+pwmAdsr.Decay), int(64.0-pwmAdsr.Sustain*64))), pos.Add(image.Pt(256, int(64.0-pwmAdsr.Sustain*64))), color, 2)
											return
										}
									},
								),
							),
						),
					),
				),
			),
		),
	)

}

func initWaveStr() {
	for i := 0; i < 4; i++ {
		wavStr := ""
		morphStr := ""
		volStr := ""
		phaseStr := ""
		for _, val := range Context.Operators[i].Wavetable {
			wavStr += strconv.Itoa(int(val)) + " "
		}
		for _, val := range Context.Operators[i].MorphWave {
			morphStr += strconv.Itoa(int(val)) + " "
		}
		for _, val := range Context.Operators[i].VolEnv {
			volStr += strconv.Itoa(int(val)) + " "
		}
		for _, val := range Context.Operators[i].PhaseEnv {
			phaseStr += strconv.Itoa(int(val)) + " "
		}
		waveStrs[i] = wavStr
		morphStrs[i] = morphStr
		envStrs[i] = volStr
		phaseStrs[i] = phaseStr
	}
}

func LoadImage(file fs.File) (*image.RGBA, error) {

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("LoadImage: error decoding png image: %w", err)
	}

	return g.ImageToRgba(img), nil
}
func LoadImageOnly(file fs.File) (*image.Image, error) {

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("LoadImage: error decoding png image: %w", err)
	}

	return &img, nil
}

//go:embed assets
var f embed.FS
var algTextures = make([]*g.Texture, 12)
var algImages = make([]*image.Image, 12)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	// var len = 512
	// for i := 0; i < len; i++ {
	// 	sineTable = append(sineTable, uint8(math.Round((math.Sin((2*math.Pi*float64(i))/float64(len))+1)*127.5)))
	// }
	Context = kurumi.ConstructSynth()
	kurumi.SynthContext = Context
	initWaveStr()
	kurumi.Synthesize()
	kurumi.Synthesize()
	kurumi.InitAudio()
	myIconF, _ := f.Open("assets/kuruicon.png")
	myIconImg, _ := LoadImage(myIconF)
	println(rs.GetTitle())

	wnd := g.NewMasterWindow("Kurumi 3 ~ "+rs.GetTitle(), 1280, 720, 0)
	wnd.SetIcon([]image.Image{myIconImg})
	// Load algorithms pictures
	for i := 0; i < 12; i++ {
		e := i
		println("assets/algs/alg" + strconv.Itoa(e) + ".png")
		tmp, _ := f.Open("assets/algs/alg" + strconv.Itoa(e) + ".png")
		img, _ := LoadImageOnly(tmp)
		algImages[i] = img
		// g.EnqueueNewTextureFromRgba(img, func(tex *g.Texture) {
		// 	algTextures[int(math.Min(float64(i), 11))] = tex
		// })
	}

	// beepStream := beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
	// 	return len(samples), true
	// })

	// go func() {
	// 	for {
	// 		for i := range beepStream.(*beep.StreamerFunc).Samples {}
	// 	}
	// }

	wnd.Run(loop)
	// openbrowser("https://www.youtube.com/watch?v=xvFZjo5PgG0")
}
