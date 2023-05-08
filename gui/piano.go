package gui

import (
	g "github.com/AllenDang/giu"
	"github.com/system64MC/Kurumi-Go/kurumi"
)

func drawPiano() {
	g.Window("Piano").Pos(340, 780).Flags(g.WindowFlagsAlwaysAutoResize|g.WindowFlagsNoResize).Layout(
		g.Checkbox("Use Sequence", &kurumi.PianState.UseSequence),
		g.Row(
			g.SliderInt(&kurumi.PianState.Octave, 0, 7),
			g.Tooltip("Changes the octave of the notes"),
			g.Label("Octave"),
		),
		g.Row(
			g.Button("A").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 0
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("As").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 1
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("B").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 2
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("C").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 3
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("Cs").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 4
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("D").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 5
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("Ds").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 6
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("E").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 7
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("F").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 8
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("Fs").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 9
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("G").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 10
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
			g.Button("Gs").Size(30, 40).OnClick(func() {
				kurumi.PianState.Key = 11
				kurumi.PianState.IsPressed = true
				if kurumi.PianState.UseSequence {
					kurumi.SynthContext.Macro = 0
					kurumi.Synthesize()
				}
			}),
		),
	)
}
