package imgui

import (
	"fmt"

	m32 "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var col = rl.LightGray

var notes = [...]string{
	0:  "C-",
	1:  "C#",
	2:  "D-",
	3:  "D#",
	4:  "E-",
	5:  "F-",
	6:  "F#",
	7:  "G-",
	8:  "G#",
	9:  "A-",
	10: "A#",
	11: "B-",
}

func IsHover(x int32, y int32, w int32, h int32) bool {
	mouse := rl.GetMousePosition()
	return x <= int32(mouse.X) && int32(mouse.X) <= x+w &&
		y <= int32(mouse.Y) && int32(mouse.Y) <= y+h
}

func Note(note *int32) {
	height := int32(20)

	startX := x
	{
		width := int32(32)
		display := "--"
		if *note > 0 && *note < 128 {
			display = notes[*note%12]
		}
		rl.DrawText(display, x, y, 20, col)

		hover := IsHover(x, y, width, height)
		wheel := rl.GetMouseWheelMove()

		if hover && m32.Abs(wheel) > 0.1 {
			*note += int32(wheel)
		}

		if hover {
			rl.DrawRectangleLines(x, y, width, height, col)
		}

		x += width
	}

	{
		width := int32(20)
		display := "-"
		if *note > 0 && *note < 128 {
			display = fmt.Sprint(*note / 12)
		}
		rl.DrawText(display, x, y, 20, col)

		hover := IsHover(x, y, width, height)
		wheel := rl.GetMouseWheelMove()

		if hover && m32.Abs(wheel) > 0.1 {
			*note += int32(wheel) * 12
		}

		if hover {
			rl.DrawRectangleLines(x, y, width, height, col)
		}
		x += 20
	}

	// rl.DrawRectangleLines(x, y, width, height, rl.NewColor(230, 230, 230, 255))

	x = startX
	y += 5 + 20

}

func NoteToFreq(note int32) float32 {
	return 440.0 * m32.Pow(2, float32(note-69)/12)
}

func Button(label string) bool {
	rl.DrawText(label, x, y, 20, rl.LightGray)

	mouse := rl.GetMousePosition()
	width := int32(3 * 20)
	height := int32(20)
	hover := x <= int32(mouse.X) && int32(mouse.X) <= x+width &&
		y <= int32(mouse.Y) && int32(mouse.Y) <= y+height

	if hover && active == nullActive || active == idx {
		rl.DrawRectangleLines(x, y, width, height, rl.NewColor(230, 230, 230, 255))
	}

	if hover && rl.IsMouseButtonDown(rl.MouseLeftButton) && active == nullActive {
		active = idx
		return true
	}

	if !rl.IsMouseButtonDown(rl.MouseLeftButton) && active == idx {
		active = nullActive
	}

	idx++
	return false
}

func Label(text string) {
	rl.DrawText(text, x, y, 20, col)
	y += 20
}
