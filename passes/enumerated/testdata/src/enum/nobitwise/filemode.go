package nobitwise

type FileMode uint32 // want FileMode:"^permittedtype$" FileMode:"^enumerated$" FileMode:"^elements:ModePerm = 511$"

const (
	ModeDir FileMode = 1 << (32 - 1 - iota)
	ModeAppend
	ModeExclusive
	ModeTemporary
	ModeSymlink
	ModeDevice
	ModeNamedPipe
	ModeSocket
	ModeSetuid
	ModeSetgid
	ModeCharDevice
	ModeSticky
	ModeIrregular

	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice | ModeCharDevice | ModeIrregular

	ModePerm FileMode = 0777 // want ModePerm:"^enumerated$"    ModePerm:"^elementof:FileMode$"
)
