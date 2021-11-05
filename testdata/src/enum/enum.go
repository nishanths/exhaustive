// want package:"^AcrossBlocksDeclsFiles:Here,Separate,There; ByteEnum:ByteA; Float64Enum:Float64A,Float64B; Int32Enum:Int32A,Int32B; IotaEnum:IotaA,IotaB; ParenVal:ParenVal0,ParenVal1; RepeatedValue:RepeatedValueA,RepeatedValueB; RuneEnum:RuneA; StringEnum:StringA,StringB,StringC; UIntEnum:UIntA,UIntB; UnexportedMembers:unexportedMembersA,unexportedMembersB; VarConstMixed:VCMixedB$"

package enum

// Var members (as opposed const members) cannot be enum members.

type VarMember int

var (
	VarMemberA VarMember = 1
	VarMemberB VarMember = 2
)

// Mixed var and const declarations (only const are members)

type VarConstMixed int

var (
	VCMixedA VarConstMixed = 0
)

const (
	VCMixedB VarConstMixed = 1
)

// Basic iota test

type IotaEnum uint8

const (
	IotaA IotaEnum = iota << 1
	IotaB
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

const Separate AcrossBlocksDeclsFiles = 1

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

type ParenVal int

const (
	ParenVal0 ParenVal = 0
	ParenVal1 ParenVal = (1)
)
