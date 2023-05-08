package gui

import (
	"image"
	"image/color"
	"math"
	"strconv"

	g "github.com/AllenDang/giu" // giu package is imported here.
	"github.com/system64MC/Kurumi-Go/kurumi"
)

func drawOperators() {
	g.Window("Operators").Pos(340, 370).Size(650, 400).Flags(g.WindowFlagsNoCollapse).Layout(
		g.TabBar().TabItems(
			buildOperator(0),
			buildOperator(1),
			buildOperator(2),
			buildOperator(3),
		),
	)
}

// Generate imports.
func buildOperator(a int) *g.TabItemWidget {

	// waveStrs[opId] = &kurumi.Waveforms[kurumi.SynthContext.Operators[opId].WaveformId]
	e := a
	str := "Operator " + strconv.Itoa(e+1)
	opId := e
	// out := fmt.Sprintf(str, i)
	// println(opId)
	return g.TabItem(str).Layout(
		g.ContextMenu().Layout(
			g.Button("Copy operator").OnClick(func() {
				// println("Copying operator" + strconv.Itoa(a))
				kurumi.CopyOp(&kurumi.SynthContext.Operators[a])
				g.CloseCurrentPopup()
			}),
			g.Button("Paste operator").OnClick(func() {
				kurumi.PasteOp(&kurumi.SynthContext.Operators[a])
				InitWaveStr()
				// println("Pasting operator" + strconv.Itoa(a))
				g.CloseCurrentPopup()
			}),
		),
		g.Row(
			g.Column(
				g.Row(
					g.Label("Mod. Mode"),
					g.Combo("", kurumi.ModModes[kurumi.SynthContext.Operators[opId].ModMode], kurumi.ModModes, &kurumi.SynthContext.Operators[opId].ModMode).Size(128).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Tooltip("Changes the modulation mode.\nFor example, if it is set to MUL, the operator will take the input and then multiply it with current oscillator state.\nDefault is FM"),
				),
				g.Row(g.Label("")),

				g.Row(
					g.SliderFloat(&kurumi.SynthContext.Operators[opId].Tl, 0, 4).Size(256).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Tooltip("Changes Mod. Depth / Volume"),
					g.Label("Mod. Depth"),
				),
				g.Row(
					g.SliderFloat(&kurumi.SynthContext.Operators[opId].Feedback, 0, 4).Size(256).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Tooltip("Changes the Feedback of operator"),
					g.Label("Feedback"),
				),
				g.Row(
					g.SliderInt(&kurumi.SynthContext.Operators[opId].Mult, 0, 32).Size(256).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Tooltip("Changes the frequency multiplier"),
					g.Label("Mult"),
				),
				g.Row(
					g.SliderFloat(&kurumi.SynthContext.Operators[opId].Phase, 0, 1).Size(256).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Tooltip("Changes the phase of the waveform"),
					g.Label("Phase"),
				),

				g.Row(g.Label("")),

				g.Checkbox("Use envelope on volume", &kurumi.SynthContext.Operators[opId].IsEnvelopeEnabled).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("If enabled, operator's volume / Mod. depth will be affected by envelope"),

				g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].IsEnvelopeEnabled).To(
					g.Row(
						g.Checkbox("Use custom volume envelope instead of ADSR", &kurumi.SynthContext.Operators[opId].UseCustomVolEnv).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
					),
					g.Tooltip("Disables ADSR and enables user-defined envelope"),
					g.Style().SetDisabled(kurumi.SynthContext.Operators[opId].UseCustomVolEnv).To(
						g.Column(
							g.Row(
								g.SliderInt(&kurumi.SynthContext.Operators[opId].Adsr.Attack, 0, 256).Size(256).OnChange(func() {

									kurumi.Synthesize()
									kurumi.GenerateWaveStr()
								}),
								g.Label("Attack"),
							),
							g.Row(
								g.SliderInt(&kurumi.SynthContext.Operators[opId].Adsr.Decay, 0, 256).Size(256).OnChange(func() {

									kurumi.Synthesize()
									kurumi.GenerateWaveStr()
								}),
								g.Label("Decay"),
							),
							g.Row(
								g.SliderFloat(&kurumi.SynthContext.Operators[opId].Adsr.Sustain, 0, 4).Size(256).OnChange(func() {

									kurumi.Synthesize()
									kurumi.GenerateWaveStr()
								}),
								g.Label("Sustain"),
							),
						),
					),
				),

				g.Row(g.Label("")),
				// g.Row(
				// 	g.Checkbox("Phase modulation", &kurumi.SynthContext.Operators[opId].PhaseMod).OnChange(func() {
				// 		kurumi.ResetFB()
				// 		kurumi.Synthesize()
				// 		kurumi.Synthesize()
				// 	}),
				// ),
				g.Column(
					g.Row(
						g.SliderInt(&kurumi.SynthContext.Operators[opId].Detune, -32, 32).Size(256).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Tooltip("If not zero, the phase will change according to sequence"),
						g.Label("Detune"),
					),
					// g.Row(
					// 	g.Checkbox("Reverse phase", &kurumi.SynthContext.Operators[opId].PhaseRev).OnChange(func() {
					// 		kurumi.ResetFB()
					// 		kurumi.Synthesize()
					// 		kurumi.Synthesize()
					// 	}),
					// ),
					g.Row(
						g.Checkbox("Use custom phase envelope", &kurumi.SynthContext.Operators[opId].CustomPhaseEnv).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Tooltip("If enabled, you can control the phase with a custom envelope.\nIt can be used to imitate vibrato."),
					),
				),
				g.Row(g.Label("")),
				g.Checkbox("Enable wavetable morphing", &kurumi.SynthContext.Operators[opId].Morphing).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("If enabled, the custom waveform will be morphed to another one"),
				g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].Morphing).To(
					g.Row(
						g.SliderInt(&kurumi.SynthContext.Operators[opId].MorphTime, 1, 256).Size(256).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Tooltip("How long the morphing takes"),
						g.Label("Morph Time"),
					),
				),

				g.Row(
					g.SliderFloat(&kurumi.SynthContext.Operators[opId].DutyCycle, 0, 1).Size(256).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Tooltip("Changes the duty cycle of the Pulse waveform.\nDefault is 0.5 (50%)"),
					g.Label("Duty Cycle"),
				),

				g.Checkbox("Enable PWM ADSR", &kurumi.SynthContext.Operators[opId].PwmAdsrEnabled).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("If enabled, the pulse width will change through time."),

				g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].PwmAdsrEnabled).To(
					g.Row(
						g.SliderFloat(&kurumi.SynthContext.Operators[opId].PwmAdsr.Start, 0, 1).Size(256).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Label("Start Duty"),
					),
					g.Row(
						g.SliderInt(&kurumi.SynthContext.Operators[opId].PwmAdsr.Attack, 0, 256).Size(256).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Label("Attack"),
					),
					g.Row(
						g.SliderInt(&kurumi.SynthContext.Operators[opId].PwmAdsr.Decay, 0, 256).Size(256).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Label("Decay"),
					),
					g.Row(
						g.SliderFloat(&kurumi.SynthContext.Operators[opId].PwmAdsr.Sustain, 0, 1).Size(256).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
						g.Label("Final Duty"),
					),
				),
			),
			g.Column(

				g.Row(
					g.Label("Waveform :"),
					// g.Custom(func() {
					// 	combo := g.Combo(kurumi.Waveforms[kurumi.SynthContext.Operators[opId].WaveformId], kurumi.Waveforms[kurumi.SynthContext.Operators[opId].WaveformId], kurumi.Waveforms, &kurumi.SynthContext.Operators[opId].WaveformId).Size(256)
					// 	combo = combo
					// }),
					g.Combo("", kurumi.Waveforms[kurumi.SynthContext.Operators[opId].WaveformId], kurumi.Waveforms, &kurumi.SynthContext.Operators[opId].WaveformId).Size(170).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
				),
				g.Checkbox("Reverse waveform", &kurumi.SynthContext.Operators[opId].Reverse).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
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
									if kurumi.SynthContext.Operators[opId].Reverse {
										y = -kurumi.WaveFuncs[kurumi.SynthContext.Operators[opId].WaveformId](&kurumi.SynthContext.Operators[opId], float64(i)/256.0)*-32 + 1

									} else {
										y = kurumi.WaveFuncs[kurumi.SynthContext.Operators[opId].WaveformId](&kurumi.SynthContext.Operators[opId], float64(i)/256.0)*-32 + 1

									}
									canvas.AddRectFilled(pos.Add(image.Pt(int(math.Floor(float64(i))), 32)),
										pos.Add(image.Pt(int(math.Ceil((float64(i))+1)), int(math.Abs(float64(-32-y))))), color, 0, 0)
								}

							},
						),
					).MinHeight(64).BgColor(color.RGBA{200, 75, 75, 0}),
				).InnerWidth(256),
				g.Style().SetDisabled(kurumi.SynthContext.Operators[opId].WaveformId < int32(len(kurumi.Waveforms)-4)).To(
					g.Row(
						g.Label("Interpolation :"),
						g.Combo("", kurumi.Interpolations[kurumi.SynthContext.Operators[opId].Interpolation], kurumi.Interpolations, &kurumi.SynthContext.Operators[opId].Interpolation).Size(160).OnChange(func() {

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
					),
					g.InputText(&waveStrs[opId]).Size(256).OnChange(func() {
						// println(waveStrs[opId])
						kurumi.ApplyStringToWaveform(opId, waveStrs[opId], kurumi.DestWave)

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].Morphing).To(
						g.Label("Waveform to morph to :"),
						g.InputText(&morphStrs[opId]).Size(256).OnChange(func() {
							kurumi.ApplyStringToWaveform(opId, morphStrs[opId], kurumi.DestMorph)

							kurumi.Synthesize()
							kurumi.GenerateWaveStr()
						}),
					),
				),

				g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].IsEnvelopeEnabled).To(
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
										if !kurumi.SynthContext.Operators[opId].UseCustomVolEnv {
											adsr := kurumi.SynthContext.Operators[opId].Adsr
											tl := kurumi.SynthContext.Operators[opId].Tl
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
										env := kurumi.SynthContext.Operators[opId].VolEnv
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

				g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].UseCustomVolEnv).To(
					g.InputText(&envStrs[opId]).Size(256).OnChange(func() {
						kurumi.ApplyStringToWaveform(opId, envStrs[opId], kurumi.DestVolEnv)

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
				),

				g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].CustomPhaseEnv && kurumi.SynthContext.Operators[opId].Detune == 0).To(
					g.Label("Phase envelope :"),
					g.Table().Size(260, 68).Rows(
						g.TableRow(
							g.Custom(
								func() {
									canvas := g.GetCanvas()
									pos := g.GetCursorScreenPos()
									color := color.RGBA{200, 75, 75, 255}

									env := kurumi.SynthContext.Operators[opId].PhaseEnv
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
						kurumi.GenerateWaveStr()
					}),
				),

				g.Style().SetDisabled(!kurumi.SynthContext.Operators[opId].PwmAdsrEnabled).To(
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

										pwmAdsr := kurumi.SynthContext.Operators[opId].PwmAdsr

										// For drawing ADSR
										{
											//adsr := kurumi.SynthContext.Operators[opId].Adsr
											tl := kurumi.SynthContext.Operators[opId].DutyCycle
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
