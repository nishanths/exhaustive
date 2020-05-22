// want package:"Stationery:Pencil,Pen,Paper,Eraser,sharpener"

package x

import bar "map/y"

var (
	// some other comment
	//exhaustive:ignore f
	// some other comment
	f = map[bar.Grain]int{
		bar.Wheat: 0,
	}
)

var e = map[Stationery]bool{
	Pen:       false,
	Pencil:    true,
	sharpener: true,
} //exhaustive:ignore e more comment

var (
	g = map[bar.Grain]int{
		bar.Wheat: 0,
	} //exhaustive:ignore g
)

//exhaustive:ignore i
// map `i` should be ignored, because this is GenDecl of one spec without "()".
var i = map[Stationery]bool{
	Pen:       false,
	Pencil:    true,
	sharpener: true,
}

//exhaustive:ignore GenDecl h -- map `h` should not be ignored, because GenDecl comment isn't associated with `h`
var (
	h = map[Stationery]bool{ // want "missing keys in map h of key type Stationery: Eraser, Paper"
		Pen:       false,
		Pencil:    true,
		sharpener: true,
	}
)
