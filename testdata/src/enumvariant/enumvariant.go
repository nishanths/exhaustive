// want package:"^AcrossBlocksDeclsFiles:Here,Separate,There; ByteEnum:ByteA; FloatEnum:FloatEnumA,FloatEnumB; Int32Enum:Int32A,Int32B; IotaEnum:IotaA,ItoaB; RepeatedValue:RepeatedValueA,RepeatedValueB; RuneEnum:RuneA; StringEnum:StringA,StringB,StringC; UIntEnum:UIntA,UIntB; UnexportedMembers:unexportedMembersA,unexportedMembersB; VarMembers:VarMemberA$"

package enumvariants

// Var members (as opposed const members) can be enum members too.

type VarMembers int

var (
	VarMemberA VarMembers
)

// Basic iota test

type IotaEnum uint8

const (
	IotaA IotaEnum = iota << 1
	ItoaB
)

// Memberless types cannot be enums.

type MemberlessEnum int

// Only the identifier name matters, not the value.
// So the enum type here has two members, not one.

type RepeatedValue int

const (
	RepeatedValueA RepeatedValue = 1
	RepeatedValueB RepeatedValue = 1
)

// Enum members can live across blocks, declaration types (const vs. var), and
// files.

type AcrossBlocksDeclsFiles int

const (
	Here AcrossBlocksDeclsFiles = 0
)

var Separate AcrossBlocksDeclsFiles = 1

// Basic test for enum type with all unexported members.

type UnexportedMembers int

const (
	unexportedMembersA UnexportedMembers = 1
	unexportedMembersB UnexportedMembers = 2
)

// Only top-level values and types form enums.

type NonTopLevel uint

func _nonTopLevel() {
	const (
		A NonTopLevel = 0
		B NonTopLevel = 1
	)
}
