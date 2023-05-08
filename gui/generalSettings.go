package gui

import (
	g "github.com/AllenDang/giu"
	c "github.com/atotto/clipboard"
	"github.com/system64MC/Kurumi-Go/kurumi"
)

func drawGeneralSettings() {
	g.Window("General settings").Size(340, 400).Pos(0, 220).Layout(
		g.Column(
			g.Checkbox("Normalize", &kurumi.SynthContext.Normalize).OnChange(func() {
				kurumi.Synthesize()
				kurumi.GenerateWaveStr()
			}),
			g.Tooltip("Normalize the waveform"),
			g.Style().SetDisabled(!kurumi.SynthContext.Normalize).To(
				g.Checkbox("Normalize after treatment", &kurumi.SynthContext.NewNormalizeBehavior).OnChange(func() {
					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("Normalize the waveform after treatment instead of before."),
			),
			g.Style().SetDisabled(kurumi.SynthContext.Normalize).To(
				g.Row(
					g.SliderFloat(&kurumi.SynthContext.Gain, 0, 4).OnChange(func() {

						kurumi.Synthesize()
						kurumi.GenerateWaveStr()
					}),
					g.Tooltip("Amplifies the output"),
					g.Label("Gain"),
				),
			),
			g.Row(
				g.SliderInt(&kurumi.SynthContext.SmoothWin, 0, 128).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("Smoothes the output"),
				g.Label("Avg. Filter Win."),
			),
			g.Row(
				g.SliderInt(&kurumi.SynthContext.MacLen, 1, 256).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("How many frames the sequence has"),
				g.Label("Seq. Lenght"),
			),
			g.Row(
				g.SliderInt(&kurumi.SynthContext.Macro, 0, kurumi.SynthContext.MacLen-1).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("The current sequence frame"),
				g.Label("Wav. Seq. Index"),
			),
			g.Row(
				g.SliderInt(&kurumi.SynthContext.Oversample, 1, 32).OnChange(func() {

					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
				}),
				g.Tooltip("Changes the oversample.\nOversample of 2x means everything is processed 2 times longer than the wavetable size,\nthen downsampled to its original size."),
				g.Label("Oversample"),
			),
			g.Label("Wave output :"),
			g.Row(
				g.InputText(&kurumi.WaveStr).Size(256).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAutoSelectAll),
				g.Button("Copy").OnClick(func() {
					c.WriteAll(*&kurumi.WaveStr)
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
}
