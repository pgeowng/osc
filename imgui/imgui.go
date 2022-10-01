package imgui

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