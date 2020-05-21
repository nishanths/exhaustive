package enumvariants

type UIntEnum uint

const (
	UIntA UIntEnum = 0
	UIntB UIntEnum = 1
)

type StringEnum string

const (
	StringA StringEnum = "stringa"
	StringB StringEnum = "stringb"
	StringC StringEnum = "stringc"
)

type RuneEnum rune

const (
	RuneA RuneEnum = 'a'
)

type ByteEnum byte

const (
	ByteA ByteEnum = 'a'
)

type Int32Enum int32

const (
	Int32A Int32Enum = 0
	Int32B Int32Enum = 1
)

type BoolNotEnum bool

const (
	BoolNotEnumA BoolNotEnum = true
	BoolNotEnumB BoolNotEnum = false
)

const There AcrossBlocksDeclsFiles = 2
