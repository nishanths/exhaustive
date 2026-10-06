package gostd

import (
	"go/types"
	"net"
)

func f1() bool {
	var t types.Type

	switch t.(type) {
	case *types.Array,
		*types.Basic,
		*types.Chan,
		*types.Interface,
		*types.Map,
		*types.Pointer,
		*types.Signature,
		*types.Slice,
		*types.Struct,
		*types.Tuple,
		*types.TypeParam,
		*types.Union:
		return false
	case *types.Alias:
		return true
	case *types.Named:
		return true
	default:
		return false
	}
}

func f2() bool {
	var t types.Type

	switch t.(type) { // want "^type switch not exhaustive: missing cases: \\*types.Basic, \\*types.Pointer$"
	case *types.Array,
		*types.Chan,
		*types.Interface,
		*types.Map,
		*types.Signature,
		*types.Slice,
		*types.Struct,
		*types.Tuple,
		*types.TypeParam,
		*types.Union:
		return false
	case *types.Alias:
		return true
	case *types.Named:
		return true
	default:
		return false
	}
}

func f3() {
	var addr net.Addr

	switch addr.(type) {
	case *net.IPAddr:
	case *net.IPNet:
	case *net.TCPAddr:
	case *net.UDPAddr:
	case *net.UnixAddr:
	}

	switch addr.(type) { // want "^type switch not exhaustive: missing cases: \\*net.IPAddr, \\*net.IPNet, \\*net.TCPAddr, \\*net.UDPAddr, \\*net.UnixAddr$"
	}

	switch addr.(type) { // want "^type switch not exhaustive: missing cases: \\*net.IPNet$"
	case *net.IPAddr:
	case *net.TCPAddr:
	case *net.UDPAddr:
	case *net.UnixAddr:
	}
}
