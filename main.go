package main

import (
	"fmt"
	"math"

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
	rl.InitWindow(800, 450, "rl audio")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetMasterVolume(0.25)

	rl.SetAudioStreamBufferSizeDefault(maxSamples)
	stream := rl.LoadAudioStream(sampleRate, 32, 1)
	defer rl.UnloadAudioStream(stream)

	osc1buf := make([]float32, maxSamples)
	osc2buf := make([]float32, maxSamples)
	master := make([]float32, maxSamples)

	for i := 0; i < maxSamples; i++ {
		osc1buf[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
		osc2buf[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
		master[i] = float32(math.Sin(float64((2*rl.Pi*float32(i))/2) * rl.Deg2rad))
	}

	rl.PlayAudioStream(stream)

	totalSamples := int32(maxSamples)
	samplesLeft := int32(totalSamples)

	position := rl.NewVector2(0, 0)

	osc1 := Oscillator{gain: 0.5}
	osc2 := Oscillator{gain: 0.5}

	rl.SetTargetFPS(165)
	frequency := float32(200.0)
	sampleDuration := float32(1.0) / sampleRate

	coef1 := float32(1.0)
	coef2 := float32(1.0)
	minfreq := float32(20)
	maxfreq := float32(20000)
	maxgain := float32(.5)
	mingain := float32(.0)

	for !rl.WindowShouldClose() {
		if rl.IsAudioStreamProcessed(stream) {
			numSamples := int32(0)
			if samplesLeft >= maxSamplesPerUpdate {
				numSamples = maxSamplesPerUpdate
			} else {
				numSamples = samplesLeft
			}
			// fmt.Println("numSamples", numSamples, samplesLeft, stream.Buffer)
			osc1.UpdateSignal(osc1buf, frequency, sampleDuration)
			osc2.UpdateSignal(osc2buf, frequency, sampleDuration)
			for i := 0; i < maxSamplesPerUpdate; i++ {
				master[i] = osc1buf[i] + osc2buf[i]
			}

			rl.UpdateAudioStream(stream, master[totalSamples-samplesLeft:], numSamples)

			samplesLeft -= numSamples

			// Reset samples feeding (loop audio)
			if samplesLeft <= 0 {
				samplesLeft = totalSamples
			}
		}

		startSample := int(sampleRate*rl.GetTime()) % maxSamples
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("sine wave should be playing!", 240, 140, 20, rl.LightGray)
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
			position.Y = 650 + 50*master[(startSample+i)%maxSamples]
			rl.DrawPixelV(position, rl.Green)
		}

		height2 := float32(rl.GetScreenHeight() / 2)
		width2 := float32(rl.GetScreenWidth() / 2)
		for i := 1; i < int(minX); i++ {
			// newPos := rl.Vector2{
			// 	X: m32.Cos(float32(i))*height2*master[(startSample+i)%maxSamples] + width2,
			// 	Y: m32.Sin(float32(i))*height2*master[(startSample+i)%maxSamples] + height2,
			// }

			// oldPos := rl.Vector2{
			// 	X: m32.Cos(float32(i-1))*height2*master[(startSample+i+maxSamples-1)%maxSamples] + width2,
			// 	Y: m32.Sin(float32(i-1))*height2*master[(startSample+i+maxSamples-1)%maxSamples] + height2,
			// }

			// rl.DrawLine(int32(oldPos.X), int32(oldPos.Y), int32(newPos.X), int32(newPos.Y), rl.NewColor(200, 122, 255, uint8(i/minX*255)))
			position.X = m32.Cos(float32(i))*height2*master[(startSample+i)%maxSamples] + width2
			position.Y = m32.Sin(float32(i))*height2*master[(startSample+i)%maxSamples] + height2
			rl.DrawPixelV(position, rl.Purple)
		}

		NumericSlicer("coef1", &coef1, 1.0, nil, nil, "", 1)
		NumericSlicer("freq1", &osc1.freq, coef1, &minfreq, &maxfreq, "hz", 0.01)
		NumericSlicer("gain1", &osc1.gain, 1.0, &mingain, &maxgain, "", 0.01)
		IMColumn()

		NumericSlicer("coef2", &coef2, 1.0, nil, nil, "", 1)
		NumericSlicer("freq2", &osc2.freq, coef2, &minfreq, &maxfreq, "hz", 0.01)
		NumericSlicer("gain2", &osc2.gain, 1.0, &mingain, &maxgain, "", 0.01)

		rl.EndDrawing()

		IMReset()
	}

	return
}

type Oscillator struct {
	phase float32
	freq  float32
	gain  float32
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

var gap int32 = 20
var x int32 = gap
var y int32 = gap
var margin int32 = 10

var column int32 = 100
var idx int32 = 0
var active int32 = -1
var nullActive int32 = -1

func IMReset() {
	idx = 0
	x = gap
	y = gap
}

func IMColumn() {
	x += gap + column
	y = gap
}

func NumericSlicer(label string, value *float32, coef float32, min *float32, max *float32, units string, step float32) {
	mouse := rl.GetMousePosition()
	width := int32(100) + 2*margin
	height := int32(40) + 2*margin

	hot := x <= int32(mouse.X) && int32(mouse.X) <= x+width &&
		y <= int32(mouse.Y) && int32(mouse.Y) <= y+height
	lmb := rl.IsMouseButtonDown(rl.MouseLeftButton)
	isActive := active == idx
	noActive := active == nullActive

	if (hot && noActive) || isActive {
		rl.DrawRectangleLines(x, y, width, height, rl.NewColor(230, 230, 230, 255))
		*value += m32.Pow(rl.GetMouseWheelMove()*step, coef)

		if noActive {
			if lmb {
				active = idx
				*value += m32.Pow(rl.GetMouseDelta().Y*step, coef)
				// osc2.freq = 77 * rl.GetMousePosition().X
			}
		} else if isActive && !lmb {
			active = nullActive
		}
	}

	if isActive && lmb {
		// *value = rl.GetMousePosition().Y
		*value += m32.Pow(rl.GetMouseDelta().Y*step, coef)
	}

	if min != nil && *value <= *min {
		*value = *min
	}

	if max != nil && *value >= *max {
		*value = *max
	}

	x += margin
	y += margin

	// rl.DrawRectangle()
	rl.DrawText(label, x, y, 20, rl.LightGray)
	y += 20
	rl.DrawText(fmt.Sprintf("%.2f%s", *value, units), x, y, 20, rl.LightGray)
	y += 20

	x -= margin
	y += margin
	idx++
}
