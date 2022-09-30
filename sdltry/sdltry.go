package main

// typedef unsigned char Uint8;
// void MyAudioCallback(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

//export MyAudioCallback
func MyAudioCallback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {

	stream
	volume := 0.2
	frequency := 200.0

	for sid := 0; sid

	return
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() (err error) {
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}

	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800,
		600,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		return
	}

	rect := sdl.Rect{0, 0, 200, 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("quit")
				running = false
				break
			}
		}
	}

	want := sdl.AudioSpec{
		Freq:     44100,
		Format:   sdl.AUDIO_F32,
		Channels: 2,
		Silence:  0,
		Samples:  512,
		Size:     0,
		// Callback: &[0]byte{},
		Callback: sdl.AudioCallback(C.MyAudioCallback),
		UserData: nil,
	}

	have := sdl.AudioSpec{}

	sdl.OpenAudioDevice("", false, &want, &have, sdl.AUDIO_ALLOW_FORMAT_CHANGE)

	return
}
