package x

import bar "map/y"

var (
	a = map[bar.Grain]int{
		bar.Wheat: 0,
		bar.Rice:  1,
	}
)

var b = map[Stationery]bool{
	Pen:       false,
	Pencil:    true,
	Paper:     true,
	Eraser:    true,
	sharpener: true,
}
