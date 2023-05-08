package kurumi

// typedef unsigned char Uint8;
// typedef unsigned short Uint16;
// typedef short Int16;
// void Wavetable(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"math"
	"reflect"
	"unsafe"

	g "github.com/AllenDang/giu"
	"github.com/veandco/go-sdl2/sdl"
)

var phase = 0.0

func phaseAcc(len int) float64 {
	freqTable := float64(sampleHz) / float64(len)
	playfreq := pianoKeys[((PianState.Octave)*12)+int32(PianState.Key)] / freqTable
	phase = math.Mod((phase + playfreq), float64(len))
	return math.Mod(phase, float64(len))
}

//export Wavetable
func Wavetable(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Int16)(unsafe.Pointer(&hdr))

	for i := 0; i < (n / 2); i++ {
		sample := float64(0)
		ind := 0
		if !SynthContext.SongPlaying {
			buf[i] = (0)
			continue
		}
		if SynthContext.SongPlaying && len(WaveOutput) > 0 && WaveOutput != nil {
			ind = int(phaseAcc(len(WaveOutput)))
		}
		if SynthContext.SongPlaying && len(WaveOutput) > 0 && WaveOutput != nil {
			sample = float64(WaveOutput[ind%len(WaveOutput)])
		}

		sample = sample * (255 / float64(SynthContext.WaveHei))
		s2 := C.Int16(sample) - 128
		buf[i] = C.Int16((s2) << 4)
	}

	if PianState.UseSequence && PianState.IsPressed {
		SynthContext.Macro++
		if SynthContext.Macro > SynthContext.MacLen {
			SynthContext.Macro = SynthContext.MacLen - 1
			PianState.IsPressed = false
		}
		Synthesize()
		g.Update()
	}
}

func InitAudio() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		println(err)
		return
	}

	spec := &sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_S16,
		Channels: 1,
		Samples:  512,
		Callback: sdl.AudioCallback(C.Wavetable),
	}
	if err := sdl.OpenAudio(spec, nil); err != nil {
		println(err)
		return
	}
	sdl.PauseAudio(false)
}
