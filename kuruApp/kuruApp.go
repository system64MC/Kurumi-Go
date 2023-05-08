package kuruApp

import (
	kuruGui "github.com/system64MC/Kurumi-Go/gui"
	"github.com/system64MC/Kurumi-Go/kurumi"
)

func Init() {
	kurumi.SynthContext = kurumi.ConstructSynth()
	kuruGui.InitWaveStr()
	kurumi.Synthesize()
	kurumi.GenerateWaveStr()
	kurumi.InitAudio()

	kurumi.EncodeJson()

	kuruGui.InitGui().Run(func() {
		kuruGui.Draw()
	})
}
