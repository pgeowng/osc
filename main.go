package main

import (
	"fmt"
	"math"

	"osc/imgui"

	m32 "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

const (
	sampleRate          = 44100
	maxSamples          = sampleRate / 8
	maxSamplesPerUpdate = sampleRate / 8
	channelCount        = 2
)

func run() (err error) {
	rl.InitWindow(1024, 768, "rl audio")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetMasterVolume(0.25)

	rl.SetAudioStreamBufferSizeDefault(maxSamples)
	stream := rl.LoadAudioStream(sampleRate, 32, 1)
	defer rl.UnloadAudioStream(stream)

	osc1buf := make([]float32, maxSamples)
	osc2buf := make([]float32, maxSamples)
	buffer := make([]float32, maxSamples)

	for i := 0; i < maxSamples; i++ {
		osc1buf[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
		osc2buf[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
		buffer[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
	}

	rl.PlayAudioStream(stream)

	// totalSamples := int32(maxSamples)
	// samplesLeft := int32(totalSamples)

	position := rl.NewVector2(0, 0)

	osc1 := Oscillator{gain: 0.5}
	osc2 := Oscillator{gain: 0.5}

	rl.SetTargetFPS(165)
	// frequency := float32(200.0)
	sampleDuration := float32(1.0) / sampleRate

	coef1 := float32(1.0)
	coef2 := float32(1.0)
	minfreq := float32(20)
	maxfreq := float32(20000)
	maxgain := float32(.5)
	mingain := float32(.0)
	tempo := float32(168.0)
	mintempo := float32(10)
	maxtempo := float32(999)

	bar := [16]int32{
		48, 50, 52, 53,
		55, 53, 52, 50,
		48, -1, -1, -1,
		-1, -1, -1, -1,
	}
	isPlaying := false

	framesPerNote := int32(tempo * 60.0 * sampleRate / 4)
	framesPast := int32(0)

	for !rl.WindowShouldClose() {
		if rl.IsAudioStreamProcessed(stream) {

			bufferPtr := 0
			for i := 0; i < maxSamplesPerUpdate; i++ {
				framesPast = (framesPast + int32(i)) % (int32(len(bar)) * framesPerNote)

				currNote := bar[(framesPast/framesPerNote)%int32(len(bar))]
				if framesPast%framesPerNote == 0 {
					osc1.Reset()
				}

				buffer[bufferPtr] = osc1.NextSample(imgui.NoteToFreq(currNote), sampleDuration)
				bufferPtr = (bufferPtr + 1) % maxSamples
			}

			// numSamples := int32(0)
			// if samplesLeft >= maxSamplesPerUpdate {
			// 	numSamples = maxSamplesPerUpdate
			// } else {
			// 	numSamples = samplesLeft
			// }
			// /imgui.NoteToFreq(note/25)
			// fmt.Println("numSamples", numSamples, samplesLeft, stream.Buffer)
			// osc1.UpdateSignal(osc1buf, frequency, sampleDuration)
			// osc2.UpdateSignal(osc2buf, frequency, sampleDuration)
			// for i := 0; i < maxSamplesPerUpdate; i++ {
			// 	buffer[i] = osc1buf[i] + osc2buf[i]
			// }

			rl.UpdateAudioStream(stream, buffer, maxSamples)

			// samplesLeft -= numSamples

			// Reset samples feeding (loop audio)
			// if samplesLeft <= 0 {
			// 	samplesLeft = totalSamples
			// }
		}

		startSample := int(sampleRate*rl.GetTime()) % maxSamples
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		// rl.DrawText("sine wave should be playing!", 240, 140, 20, rl.LightGray)
		// rl.DrawText(fmt.Sprintf("freq: %.2f", osc1.GetFrequency()), 0, 0, 20, rl.Red)
		// rl.DrawText(fmt.Sprintf("freq: %.2f", osc1.GetFrequency()), 0, 20, 20, rl.Yellow)

		minX := int(maxSamples)
		if minX > rl.GetScreenWidth() {
			minX = rl.GetScreenWidth()
		}
		for i := 0; i < int(minX); i++ {
			position.X = float32(i)
			position.Y = 250 + 50*osc1buf[(startSample+i)%maxSamples]
			rl.DrawPixelV(position, rl.Red)
		}

		for i := 0; i < int(minX); i++ {
			position.X = float32(i)
			position.Y = 450 + 50*osc2buf[(startSample+i)%maxSamples]
			rl.DrawPixelV(position, rl.Yellow)
		}

		for i := 0; i < int(minX); i++ {
			position.X = float32(i)
			position.Y = 650 + 50*buffer[(startSample+i)%maxSamples]
			rl.DrawPixelV(position, rl.Green)
		}

		height2 := float32(rl.GetScreenHeight() / 2)
		width2 := float32(rl.GetScreenWidth() / 2)
		for i := 1; i < int(minX); i++ {
			// newPos := rl.Vector2{
			// 	X: m32.Cos(float32(i))*height2*buffer[(startSample+i)%maxSamples] + width2,
			// 	Y: m32.Sin(float32(i))*height2*buffer[(startSample+i)%maxSamples] + height2,
			// }

			// oldPos := rl.Vector2{
			// 	X: m32.Cos(float32(i-1))*height2*buffer[(startSample+i+maxSamples-1)%maxSamples] + width2,
			// 	Y: m32.Sin(float32(i-1))*height2*buffer[(startSample+i+maxSamples-1)%maxSamples] + height2,
			// }

			// rl.DrawLine(int32(oldPos.X), int32(oldPos.Y), int32(newPos.X), int32(newPos.Y), rl.NewColor(200, 122, 255, uint8(i/minX*255)))
			position.X = m32.Cos(float32(i))*height2*buffer[(startSample+i)%maxSamples] + width2
			position.Y = m32.Sin(float32(i))*height2*buffer[(startSample+i)%maxSamples] + height2
			rl.DrawPixelV(position, rl.Purple)
		}

		imgui.NumericSlicer("coef1", &coef1, 1.0, nil, nil, "", 1)
		imgui.NumericSlicer("freq1", &osc1.freq, coef1, &minfreq, &maxfreq, "hz", 0.01)
		imgui.NumericSlicer("gain1", &osc1.gain, 1.0, &mingain, &maxgain, "", 0.01)
		imgui.IMColumn()

		imgui.NumericSlicer("coef2", &coef2, 1.0, nil, nil, "", 1)
		imgui.NumericSlicer("freq2", &osc2.freq, coef2, &minfreq, &maxfreq, "hz", 0.01)
		imgui.NumericSlicer("gain2", &osc2.gain, 1.0, &mingain, &maxgain, "", 0.01)

		imgui.IMColumn()

		imgui.NumericSlicer("tempo", &tempo, 1.0, &mintempo, &maxtempo, "bpm", 1)

		imgui.IMColumn()

		togglePlaying := false
		if isPlaying {
			togglePlaying = imgui.Button("stop")
		} else {
			togglePlaying = imgui.Button("play")
		}

		if togglePlaying {
			// if isPlaying {
			// 	stop := true
			// } else {
			// 	nextPlay := true
			// 	playbackTime = 0
			// }
			isPlaying = !isPlaying
		}

		imgui.IMColumn()

		for i := 0; i < len(bar); i++ {
			imgui.Note(&bar[i])
		}

		imgui.IMColumn()

		imgui.Label(fmt.Sprintf("%d; %.2f", framesPast, float32(framesPast)/sampleRate))

		// val := int32(note / 25)
		// imgui.Note(&val)

		// note++
		// note = note % (128 * 25)

		rl.EndDrawing()

		imgui.IMReset()
	}

	return
}

type Oscillator struct {
	phase     float32
	freq      float32
	gain      float32
	sampleIdx int32
}

func (osc *Oscillator) UpdateSignal(data []float32, frequency float32, sampleDuration float32) {
	// for i := 0; i < maxSamples; i++ {
	// 	signal[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
	// }

	for i := 0; i < maxSamples; i++ {
		data[i] = m32.Sin(osc.phase+2*rl.Pi*frequency*sampleDuration*float32(i)) * osc.gain
		osc.phase += sampleDuration / frequency * osc.freq
	}
}

func (osc *Oscillator) NextSample(frequency float32, sampleDuration float32) float32 {
	frame := m32.Sin(2*rl.Pi*frequency*sampleDuration*float32(osc.sampleIdx)) * osc.gain
	osc.sampleIdx++
	return frame
}

func (osc *Oscillator) Reset() {
	osc.sampleIdx = 0
}

// Name
// units dB
// var type float
// low -96
// high 0
// init -6
type Volume struct {
	volume float32
}

func (n *Volume) ProcessAudioFrame(input []float32, output []float32) {
	output[0] = input[0] * n.volume

	if len(input) == 1 && len(output) == 2 {
		output[1] = input[0] * n.volume
	}

	if len(input) == 2 && len(output) == 2 {
		output[1] = input[1] * n.volume
	}
}

func (osc *Oscillator) GetFrequency() float32 {
	return osc.freq
}
