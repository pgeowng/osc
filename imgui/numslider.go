package imgui

import (
	"fmt"

	m32 "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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