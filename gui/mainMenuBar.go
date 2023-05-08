package gui

import (
	"fmt"
	"log"
	"runtime"

	g "github.com/AllenDang/giu"

	"os/exec"

	"github.com/ncruces/zenity"
	"github.com/system64MC/Kurumi-Go/kurumi"
)

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

func drawMainMenuBar() {
	g.MainMenuBar().Layout(
		g.Menu("File").Layout(
			g.MenuItem("Open").OnClick(func() {
				kurumi.LoadJson()
				InitWaveStr()
			}),
			g.Separator(),
			g.MenuItem("Save").OnClick(
				func() {
					kurumi.SaveJson()
				},
			),
			g.Menu("Export").Layout(
				g.MenuItem("Furnace Wavetable").OnClick(func() {
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

				if err == nil {
					kurumi.SynthContext = kurumi.ConstructSynth()
					InitWaveStr()
					kurumi.Synthesize()
					kurumi.Synthesize()
					kurumi.GenerateWaveStr()
					kurumi.InitAudio()
				}
			}),
			g.MenuItem("Exit"),
		),
		g.Checkbox("Enable sound", &kurumi.SynthContext.SongPlaying),
	).Build()
}
