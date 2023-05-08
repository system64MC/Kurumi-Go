package gui

import (
	"embed"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"strconv"

	g "github.com/AllenDang/giu"
	"github.com/system64MC/Kurumi-Go/kurumi"
	rs "github.com/system64MC/Kurumi-Go/randomStuff"
)

var operatorsGUIs = make([]*g.TabItemWidget, 0)
var waveStrs = []string{"255", "255", "255", "255"}
var morphStrs = []string{"255", "255", "255", "255"}
var envStrs = []string{"255", "255", "255", "255"}
var phaseStrs = []string{"0", "0", "0", "0"}

func InitWaveStr() {
	for i := 0; i < 4; i++ {
		wavStr := ""
		morphStr := ""
		volStr := ""
		phaseStr := ""
		for _, val := range kurumi.SynthContext.Operators[i].Wavetable {
			wavStr += strconv.Itoa(int(val)) + " "
		}
		for _, val := range kurumi.SynthContext.Operators[i].MorphWave {
			morphStr += strconv.Itoa(int(val)) + " "
		}
		for _, val := range kurumi.SynthContext.Operators[i].VolEnv {
			volStr += strconv.Itoa(int(val)) + " "
		}
		for _, val := range kurumi.SynthContext.Operators[i].PhaseEnv {
			phaseStr += strconv.Itoa(int(val)) + " "
		}
		waveStrs[i] = wavStr
		morphStrs[i] = morphStr
		envStrs[i] = volStr
		phaseStrs[i] = phaseStr
	}
}

//go:embed assets
var f embed.FS
var algTextures = make([]*g.Texture, 12)
var algImages = make([]*image.Image, 12)

func LoadTextures() {

}

func Draw() {
	drawFilterWindow()
	drawGeneralSettings()
	drawMainMenuBar()
	drawOperators()
	drawMatrixWindow()
	drawOperators()
	drawWavePreview()
	drawPiano()
}

func loadImage(file fs.File) (*image.RGBA, error) {

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("LoadImage: error decoding png image: %w", err)
	}

	return g.ImageToRgba(img), nil
}
func loadImageOnly(file fs.File) (*image.Image, error) {

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("LoadImage: error decoding png image: %w", err)
	}

	return &img, nil
}

func InitGui() *g.MasterWindow {
	myIconF, _ := f.Open("assets/kuruicon.png")
	myIconImg, _ := loadImage(myIconF)

	wnd := g.NewMasterWindow("Kurumi 3 ~ "+rs.GetTitle(), 1280, 900, 0)
	wnd.SetIcon([]image.Image{myIconImg})
	// Load algorithms pictures //
	for i := 0; i < 12; i++ {
		e := i
		println("assets/algs/alg" + strconv.Itoa(e) + ".png")
		tmp, _ := f.Open("assets/algs/alg" + strconv.Itoa(e) + ".png")
		img, _ := loadImageOnly(tmp)
		algImages[i] = img
	}

	return wnd
}
