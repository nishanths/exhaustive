package iotavalue

type FileMode uint32 // want FileMode:"^permittedtype$" FileMode:"^enumerated$" FileMode:"^elements:ModeDir = 2147483648, ModeAppend = 1073741824, ModeExclusive = 536870912, ModeTemporary = 268435456, ModeSymlink = 134217728, ModeDevice = 67108864, ModeNamedPipe = 33554432, ModeSocket = 16777216, ModeSetuid = 8388608, ModeSetgid = 4194304, ModeCharDevice = 2097152, ModeSticky = 1048576, ModeIrregular = 524288, ModeType = 2401763328$"

const (
	ModeDir        FileMode = 1 << (32 - 1 - iota) // want ModeDir:"^enumerated$"          ModeDir:"^elementof:FileMode$"
	ModeAppend                                     // want ModeAppend:"^enumerated$"       ModeAppend:"^elementof:FileMode$"
	ModeExclusive                                  // want ModeExclusive:"^enumerated$"    ModeExclusive:"^elementof:FileMode$"
	ModeTemporary                                  // want ModeTemporary:"^enumerated$"    ModeTemporary:"^elementof:FileMode$"
	ModeSymlink                                    // want ModeSymlink:"^enumerated$"      ModeSymlink:"^elementof:FileMode$"
	ModeDevice                                     // want ModeDevice:"^enumerated$"       ModeDevice:"^elementof:FileMode$"
	ModeNamedPipe                                  // want ModeNamedPipe:"^enumerated$"    ModeNamedPipe:"^elementof:FileMode$"
	ModeSocket                                     // want ModeSocket:"^enumerated$"       ModeSocket:"^elementof:FileMode$"
	ModeSetuid                                     // want ModeSetuid:"^enumerated$"       ModeSetuid:"^elementof:FileMode$"
	ModeSetgid                                     // want ModeSetgid:"^enumerated$"       ModeSetgid:"^elementof:FileMode$"
	ModeCharDevice                                 // want ModeCharDevice:"^enumerated$"   ModeCharDevice:"^elementof:FileMode$"
	ModeSticky                                     // want ModeSticky:"^enumerated$"       ModeSticky:"^elementof:FileMode$"
	ModeIrregular                                  // want ModeIrregular:"^enumerated$"    ModeIrregular:"^elementof:FileMode$"

	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice | ModeCharDevice | ModeIrregular // want ModeType:"^enumerated$" ModeType:"^elementof:FileMode$"

	ModePerm FileMode = 0777
)
