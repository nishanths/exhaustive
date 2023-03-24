package enum

// Integer, string, and float can be enum types.
// Bool cannot be enum type.
// Only basic types can be enum types.

// TODO: add test coverage for each type that can be an enum.

type UIntEnum uint // want UIntEnum:"^UIntA,UIntB$"

const (
	UIntA UIntEnum = 0
	UIntB UIntEnum = 1
)

type StringEnum string // want StringEnum:"^StringA,StringB,StringC$"

const (
	StringA StringEnum = "stringa"
	StringB StringEnum = "stringb"
	StringC StringEnum = "stringc"
)

type RuneEnum rune // want RuneEnum:"^RuneA$"

const (
	RuneA RuneEnum = 'a'
)

type ByteEnum byte // want ByteEnum:"^ByteA$"

const (
	ByteA = ByteEnum('a')
)

type Int32Enum int32 // want Int32Enum:"^Int32A,Int32B$"

const (
	Int32A Int32Enum = 0
	Int32B Int32Enum = 1
)

type Float64Enum float64 // want Float64Enum:"^Float64A,Float64B$"

const (
	Float64A Float64Enum = iota
	Float64B
)

type BoolNotEnum bool

const (
	BoolNotEnumA BoolNotEnum = true
	BoolNotEnumB BoolNotEnum = false
)

type NonBasicType S

type S struct{ F int }

var (
	SA NonBasicType = NonBasicType{F: 1}
	SB NonBasicType = NonBasicType{F: 2}
)

const There AcrossBlocksDeclsFiles = 2

type _ int // blank identifier type
