package kurumi

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

func EncodeJson() []byte {
	// synthJson, err := json.Marshal(SynthContext)
	// if err != nil {
	//     fmt.Println("Error:", err)
	// }

	data := map[string]interface{}{
		"Format": "vampire",
		"Synth":  SynthContext,
	}

	synthJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return synthJson
}

func SaveJson() error {

	path, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{"Kurumi Vampire Patch files", []string{"*.kvp"}, false},
		})
	if errZen == zenity.ErrCanceled {
		return errZen
	}
	if !strings.HasSuffix(path, ".kvp") {
		path += ".kvp"
	}

	data := EncodeJson()
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err2 := file.Write(data)
	if err2 != nil {
		return err2
	}
	return nil
}

func LoadJson() error {
	path, errZen := zenity.SelectFile(
		zenity.FileFilters{
			{"Kurumi Vampire Patch files", []string{"*.kvp"}, false},
		})

	if errZen == zenity.ErrCanceled {
		return errZen
	}

	jsonData, err := os.Open(path)
	defer jsonData.Close()

	decoder := json.NewDecoder(jsonData)

	// Create a map to decode the JSON into
	var data map[string]interface{}

	// Decode the JSON
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	// Extract the Synth object from the map
	synth := &Synth{}
	format := data["Format"].(string)

	if format != "vampire" {
		panic("Invalid format")
	}

	synthMap, ok := data["Synth"].(map[string]interface{})
	if ok {
		// Decode the Synth map into a Synth struct
		bytes, err := json.Marshal(synthMap)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bytes, synth)
		if err != nil {
			panic(err)
		}
	}

	SynthContext.ModMatrix = synth.ModMatrix

	SynthContext.Cutoff = synth.Cutoff
	SynthContext.FilterAdsrEnabled = synth.FilterAdsrEnabled
	SynthContext.FilterStart = synth.FilterStart
	SynthContext.FilterAttack = synth.FilterAttack
	SynthContext.FilterDecay = synth.FilterDecay
	SynthContext.FilterSustain = synth.FilterSustain
	SynthContext.FilterType = synth.FilterType
	SynthContext.FilterEnabled = synth.FilterEnabled
	SynthContext.Pitch = synth.Pitch
	SynthContext.Resonance = synth.Resonance

	SynthContext.Operators = synth.Operators
	SynthContext.OpOutputs = synth.OpOutputs

	SynthContext.Normalize = synth.Normalize

	SynthContext.WaveLen = synth.WaveLen
	SynthContext.WaveHei = synth.WaveHei

	SynthContext.MacLen = synth.MacLen
	SynthContext.Macro = synth.Macro

	SynthContext.SmoothWin = synth.SmoothWin

	SynthContext.Gain = synth.Gain

	SynthContext.Oversample = synth.Oversample

	Synthesize()

	fmt.Printf("Format: %s\n", data["Format"])
	fmt.Printf("Synth: %+v\n", synth)

	return nil
}

func CreateFTIN163(macro bool) error {
	tmpLen := SynthContext.WaveLen
	tmpHei := SynthContext.WaveHei

	if tmpLen > 240 {
		SynthContext.WaveLen = 240
	}
	if tmpHei > 15 {
		SynthContext.WaveHei = 15
	}
	Synthesize()

	fpath, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".fti files", []string{"*.fti"}, false},
		})
	if errZen == zenity.ErrCanceled {
		return errZen
	}
	if !strings.HasSuffix(fpath, ".fti") {
		fpath += ".fti"
	}
	file, err := os.Create(fpath)
	name := filepath.Base(fpath)
	name = name[:len(name)-len(filepath.Ext(fpath))]
	if len(name) > 127 {
		name = name[:127]
	}

	// It's time to build the FUW file!!
	output := []byte{
		'F', 'T', 'I', '2', '.', '4', // Header
		0x05,                     // Instrument type (N163 here)
		byte(len(name)), 0, 0, 0, // Length of name string
	}

	for _, character := range name {
		output = append(output, byte(character&0xFF))
	}

	output = append(output, 0x05) // I dunno what it is

	output = append(output, 0x00) // We disable volume envelope...
	output = append(output, 0x00)
	output = append(output, 0x00)
	output = append(output, 0x00)

	waveMacroEnabled := byte(0)
	if macro {
		waveMacroEnabled = 1
	}
	output = append(output, waveMacroEnabled) // We enable waveform envelope!

	if macro {
		macroLen := int(math.Min(float64(SynthContext.MacLen), 64))
		output = append(output, byte(macroLen))

		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		for i := 0; i < macroLen; i++ {
			output = append(output, byte(i))
		}
		output = append(output, byte(SynthContext.WaveLen))
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		output = append(output, byte(macroLen))
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		tmpMacro := SynthContext.Macro
		for m := int32(0); m < int32(macroLen); m++ {
			SynthContext.Macro = m
			Synthesize()

			for _, sample := range WaveOutput {
				output = append(output, byte(sample&0x0F))
			}
		}
		SynthContext.Macro = tmpMacro
		Synthesize()
	} else {
		output = append(output, byte(SynthContext.WaveLen))
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x01)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		for _, sample := range WaveOutput {
			output = append(output, byte(sample&0x0F))
		}
	}
	SynthContext.WaveLen = tmpLen
	SynthContext.WaveHei = tmpHei
	Synthesize()

	i, err := file.Write(output)
	i = int(i)
	if err != nil {
		return err
	}
	return nil
}

const FURNACE_FORMAT_VER uint16 = 143

func CreateFUW() error {
	fpath, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".fuw files", []string{"*.fuw"}, false},
		})
	if errZen == zenity.ErrCanceled {
		return errZen
	}
	if !strings.HasSuffix(fpath, ".fuw") {
		fpath += ".fuw"
	}
	file, err := os.Create(fpath)

	var size uint32 = 1 + 4 + 4 + 4 + uint32(4*len(WaveOutput))
	const HEADER_SIZE = 16 + 2 + 2 + 4 + 4 + 1 + 4 + 4 + 4

	// It's time to build the FUW file!!
	output := []byte{
		'-', 'F', 'u', 'r', 'n', 'a', 'c', 'e', ' ', 'w', 'a', 'v', 'e', 't', 'a', '-', // Header, 16 bytes
		byte(FURNACE_FORMAT_VER & 0xFF), byte(FURNACE_FORMAT_VER >> 8), // Format version, 2 bytes
		'0', '0', // Reserved, 2 bytes
		'W', 'A', 'V', 'E', // WAVE chunk, 4 bytes
		byte(size & 0xFF), byte((size >> 8) & 0xFF), byte((size >> 16) & 0xFF), byte((size >> 24)), // Size of chunk, 4 bytes
		0,                                                                                                                                                          //empty string, 1 byte
		byte(SynthContext.WaveLen & 0xFF), byte((SynthContext.WaveLen >> 8) & 0xFF), byte((SynthContext.WaveLen >> 16) & 0xFF), byte((SynthContext.WaveLen >> 24)), // Wave length, 4 bytes
		0, 0, 0, 0, // Reserved, 4 bytes
		byte(SynthContext.WaveHei & 0xFF), byte((SynthContext.WaveHei >> 8) & 0xFF), byte((SynthContext.WaveHei >> 16) & 0xFF), byte((SynthContext.WaveHei >> 24)), // Wave height, 4 bytes
	}

	// Appending Data
	for _, sample := range WaveOutput {
		output = append(output, byte(sample&0xFF))
		output = append(output, byte((sample>>8)&0xFF))
		output = append(output, byte((sample>>16)&0xFF))
		output = append(output, byte(sample>>24))
	}

	i, err := file.Write(output)
	i = int(i)
	if err != nil {
		return err
	}
	return nil
}

func createWavNew(path string, macro bool, bits16 bool) error {
	file, err := os.Create(path)

	var frames = 1
	if macro {
		frames = int(SynthContext.MacLen)
	}

	var bits byte = 8
	if bits16 {
		bits = 16
	}

	chunkSize := 0
	if bits16 {
		chunkSize = 36 + (len(WaveOutput)*frames)*2
	} else {
		chunkSize = 36 + (len(WaveOutput) * frames)
	}

	subchunkSize := 0
	if bits16 {
		subchunkSize = (len(WaveOutput) * frames) * 2
	} else {
		subchunkSize = (len(WaveOutput) * frames)
	}

	sampleRate := getSampleRate()
	byteRate := 0
	if bits16 {
		byteRate = (sampleRate * 16) / 8
	} else {
		byteRate = sampleRate
	}

	intBuffer := []byte{
		0x52, 0x49, 0x46, 0x46, // ChunkID: "RIFF" in ASCII form, big endian
		byte(chunkSize & 0xFF), byte((chunkSize >> 8) & 0xFF), byte((chunkSize >> 16) & 0xFF), byte(chunkSize >> 24), // ChunkSize - will be filled later,
		0x57, 0x41, 0x56, 0x45, // Format: "WAVE" in ASCII form
		0x66, 0x6d, 0x74, 0x20, // Subchunk1ID: "fmt " in ASCII form
		0x10, 0x00, 0x00, 0x00, // Subchunk1Size: 16 for PCM
		0x01, 0x00, // AudioFormat: PCM = 1
		0x01, 0x00, // NumChannels: Mono = 1
		byte(sampleRate & 0xFF), byte((sampleRate >> 8) & 0xFF), byte((sampleRate >> 16) & 0xFF), byte(sampleRate >> 24), // SampleRate: 44100 Hz - little endian
		byte(byteRate & 0xFF), byte((byteRate >> 8) & 0xFF), byte((byteRate >> 16) & 0xFF), byte(byteRate >> 24), // ByteRate: 44100 * 1 * 16 / 8 - little endian
		byte(bits / 8), 0x00, // BlockAlign: 1 * 16 / 8 - little endian
		bits, 0x00, // BitsPerSample: 16 bits per sample
		0x64, 0x61, 0x74, 0x61, // Subchunk2ID: "data" in ASCII form
		byte(subchunkSize & 0xFF), byte((subchunkSize >> 8) & 0xFF), byte((subchunkSize >> 16) & 0xFF), byte(subchunkSize >> 24), // Subchunk2Size - will be filled later
	}

	_, err = file.Write(intBuffer)
	if err != nil {
		return err
	}

	if macro {
		tmpMac := SynthContext.Macro
		for i := 0; i < int(SynthContext.MacLen); i++ {
			SynthContext.Macro = int32(i)
			Synthesize()

			for _, sample := range WaveOutput {
				var tmp float64
				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
				}
				if bits16 {

					myOut := int16(math.Round((tmp - 1) * float64((1<<(16-1))-1)))
					b1 := byte(myOut & 0xFF)
					b2 := byte(myOut >> 8)
					file.Write([]byte{b1, b2})
					continue
				}

				myOut := int16(math.Round((tmp - 0) * float64((1<<(8-1))-1)))
				file.Write([]byte{byte(myOut)})

			}
		}
		SynthContext.Macro = tmpMac
		Synthesize()
	} else {
		for _, sample := range WaveOutput {
			var tmp float64
			if (SynthContext.WaveHei & 0x0001) == 1 {
				tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
			} else {
				tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
			}
			if bits16 {

				myOut := int16(math.Round((tmp - 1) * float64((1<<(16-1))-1)))
				b1 := byte(myOut & 0xFF)
				b2 := byte(myOut >> 8)
				file.Write([]byte{b1, b2})
				continue
			}

			myOut := int16(math.Round((tmp - 0) * float64((1<<(8-1))-1)))
			file.Write([]byte{byte(myOut)})
		}
	}
	return nil
}

func createWav(path string, macro bool, bits16 bool) error {
	file, err := os.Create(path)
	var bufLen int
	if macro {
		bufLen = len(WaveOutput) * int(SynthContext.MacLen)
	} else {
		bufLen = len(WaveOutput)
	}

	var frames = 1
	if macro {
		frames = int(SynthContext.MacLen)
	}

	var bits byte = 8
	if bits16 {
		bits = 16
	}

	chunkSize := 0
	if bits16 {
		chunkSize = 36 + (bufLen*frames)*2
	} else {
		chunkSize = 36 + (bufLen * frames)
	}

	subchunkSize := 0
	if bits16 {
		subchunkSize = (bufLen * frames) * 2
	} else {
		subchunkSize = (bufLen * frames)
	}

	sampleRate := getSampleRate()
	byteRate := 0
	if bits16 {
		byteRate = (sampleRate * 16) / 8
	} else {
		byteRate = sampleRate
	}

	intBuffer := []byte{
		0x52, 0x49, 0x46, 0x46, // ChunkID: "RIFF" in ASCII form, big endian
		byte(chunkSize & 0xFF), byte((chunkSize >> 8) & 0xFF), byte((chunkSize >> 16) & 0xFF), byte(chunkSize >> 24), // ChunkSize - will be filled later,
		0x57, 0x41, 0x56, 0x45, // Format: "WAVE" in ASCII form
		0x66, 0x6d, 0x74, 0x20, // Subchunk1ID: "fmt " in ASCII form
		0x10, 0x00, 0x00, 0x00, // Subchunk1Size: 16 for PCM
		0x01, 0x00, // AudioFormat: PCM = 1
		0x01, 0x00, // NumChannels: Mono = 1
		byte(sampleRate & 0xFF), byte((sampleRate >> 8) & 0xFF), byte((sampleRate >> 16) & 0xFF), byte(sampleRate >> 24), // SampleRate: 44100 Hz - little endian
		byte(byteRate & 0xFF), byte((byteRate >> 8) & 0xFF), byte((byteRate >> 16) & 0xFF), byte(byteRate >> 24), // ByteRate: 44100 * 1 * 16 / 8 - little endian
		byte(bits / 8), 0x00, // BlockAlign: 1 * 16 / 8 - little endian
		bits, 0x00, // BitsPerSample: 16 bits per sample
		0x64, 0x61, 0x74, 0x61, // Subchunk2ID: "data" in ASCII form
		byte(subchunkSize & 0xFF), byte((subchunkSize >> 8) & 0xFF), byte((subchunkSize >> 16) & 0xFF), byte(subchunkSize >> 24), // Subchunk2Size - will be filled later
	}

	var output []uint8
	if macro {
		tmpMac := SynthContext.Macro
		for i := 0; i < int(SynthContext.MacLen); i++ {
			SynthContext.Macro = int32(i)
			Synthesize()
			for _, sample := range WaveOutput {
				var tmp float64
				if bits16 {

					if (SynthContext.WaveHei & 0x0001) == 1 {
						tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
					} else {
						tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
					}
					myOut := int16(math.Round((tmp - 1) * float64((1<<(16-1))-1)))
					output = append(output, byte(myOut&0xFF))
					output = append(output, byte(myOut>>8))
					continue
				}
				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
				}
				myOut := int16(math.Round((tmp - 0) * float64((1<<(8-1))-1)))
				output = append(output, byte(myOut))
			}
		}
		SynthContext.Macro = tmpMac
		Synthesize()
	} else {
		for _, sample := range WaveOutput {
			var tmp float64
			if bits16 {

				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
				}
				myOut := int16(math.Round((tmp - 1) * float64((1<<(16-1))-1)))
				output = append(output, byte(myOut&0xFF))
				output = append(output, byte(myOut>>8))
				continue
			}
			if (SynthContext.WaveHei & 0x0001) == 1 {
				tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
			} else {
				tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
			}
			myOut := int16(math.Round((tmp - 0) * float64((1<<(8-1))-1)))
			output = append(output, byte(myOut))
		}
	}

	for _, sample := range output {
		intBuffer = append(intBuffer, sample)
	}

	i, err := file.Write(intBuffer)
	i = int(i)
	if err != nil {
		return err
	}
	return nil
}

func SaveRaw(macro bool, mode int) error {
	path, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".raw files", []string{"*.raw"}, false},
		})
	if errZen == zenity.ErrCanceled {
		return errZen
	}
	if !strings.HasSuffix(path, ".raw") {
		path += ".raw"
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if mode == 2 {
		if len(WaveOutput)&1 > 0 {
			zenity.Error("Only an even length is accepted for 4-bits export.")
			return nil
		}
	}

	bufLen := len(WaveOutput)
	if macro {
		bufLen = bufLen * int(SynthContext.MacLen)
	}
	var output []byte
	if macro {
		tmpMac := SynthContext.Macro
		for i := 0; i < int(SynthContext.MacLen); i++ {
			SynthContext.Macro = int32(i)
			Synthesize()
			if mode == 2 {
				for i := 0; i < len(WaveOutput); i += 2 {
					var tmp1 float64
					var tmp2 float64
					if (SynthContext.WaveHei & 0x0001) == 1 {
						tmp1 = float64(WaveOutput[i]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
						tmp2 = float64(WaveOutput[i+1]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
					} else {
						tmp1 = float64(WaveOutput[i]) / (float64(SynthContext.WaveHei) / 2.0)
						tmp2 = float64(WaveOutput[i+1]) / (float64(SynthContext.WaveHei) / 2.0)
					}

					myOut1 := int16(math.Round((tmp1 - 0) * float64((1<<(4-1))-1)))
					myOut2 := int16(math.Round((tmp2 - 0) * float64((1<<(4-1))-1)))
					output = append(output, byte((myOut1>>4)|myOut2&0xF))
				}
				continue
			}
			for _, sample := range WaveOutput {
				var tmp float64

				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
				}
				switch mode {
				case 0: // Normalized RAW
					myOut := int16(math.Round((tmp - 0) * float64((1<<(8-1))-1)))
					output = append(output, byte(myOut))
				case 1: // Non Normalized
					output = append(output, byte(sample))
				}

			}
		}
		SynthContext.Macro = tmpMac
		Synthesize()
	} else {
		if mode == 2 {
			for i := 0; i < len(WaveOutput); i += 2 {
				var tmp1 float64
				var tmp2 float64
				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp1 = float64(WaveOutput[i]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
					tmp2 = float64(WaveOutput[i+1]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp1 = float64(WaveOutput[i]) / (float64(SynthContext.WaveHei) / 2.0)
					tmp2 = float64(WaveOutput[i+1]) / (float64(SynthContext.WaveHei) / 2.0)
				}

				myOut1 := int16(math.Round((tmp1 - 0) * float64((1<<(4-1))-1)))
				myOut2 := int16(math.Round((tmp2 - 0) * float64((1<<(4-1))-1)))
				output = append(output, byte((myOut1<<4)|myOut2&0xF))
			}
		}
		for _, sample := range WaveOutput {
			var tmp float64

			if (SynthContext.WaveHei & 0x0001) == 1 {
				tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
			} else {
				tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
			}
			switch mode {
			case 0: // Normalized RAW
				myOut := int16(math.Round((tmp - 0) * float64((1<<(8-1))-1)))
				output = append(output, byte(myOut))
			case 1: // Non Normalized
				output = append(output, byte(sample))
			}

		}
	}

	_, err2 := file.Write(output)
	if err2 != nil {
		return err2
	}
	return nil
}

func SaveTxt(macro bool) error {

	path, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".txt files", []string{"*.txt"}, false},
		})
	if errZen == zenity.ErrCanceled {
		return errZen
	}
	if !strings.HasSuffix(path, ".txt") {
		path += ".txt"
	}

	str := ""
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if macro {
		GenerateWaveSeqStr()
		str = WaveSeqStr
	} else {
		GenerateWaveStr()
		str = WaveStr
	}
	_, err2 := file.WriteString(str)
	if err2 != nil {
		return err2
	}
	return nil
}

func SaveFile(macro bool, bits16 bool) {
	path, err := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".WAV files", []string{"*.wav"}, false},
		})
	if err == zenity.ErrCanceled {
		return
	}
	if !strings.HasSuffix(path, ".wav") {
		path += ".wav"
	}
	createWavNew(path, macro, bits16)
}
