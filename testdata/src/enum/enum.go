package enum

// Var members (as opposed const members) cannot be enum members.

type VarMember int

var (
	VarMemberA VarMember = 1
	VarMemberB VarMember = 2
)

// Mixed var and const declarations (only const are members)

type VarConstMixed int // want VarConstMixed:"^VCMixedB$"

var (
	VCMixedA VarConstMixed = 0
)

const (
	VCMixedB VarConstMixed = 1
)

// Basic iota test

type IotaEnum uint8 // want IotaEnum:"^IotaA,IotaB$"

const (
	IotaA IotaEnum = iota << 1
	IotaB
)

// Memberless types cannot be enums.

type MemberlessEnum int

// Only the identifier name matters, not the value.
// So the enum type here has two members, not one.

type RepeatedValue int // want RepeatedValue:"^RepeatedValueA,RepeatedValueB$"

const (
	RepeatedValueA RepeatedValue = 1
	RepeatedValueB RepeatedValue = 1
)

// Enum members can live across blocks, declaration types (const vs. var), and
// files.

type AcrossBlocksDeclsFiles int // want AcrossBlocksDeclsFiles:"^Here,Separate,There$"

const (
	Here AcrossBlocksDeclsFiles = 0
)

const Separate AcrossBlocksDeclsFiles = 1

// Basic test for enum type with all unexported members.

type UnexportedMembers int // want UnexportedMembers:"^unexportedMembersA,unexportedMembersB$"

const (
	unexportedMembersA UnexportedMembers = 1
	unexportedMembersB UnexportedMembers = 2
)

type ParenVal int // want ParenVal:"^ParenVal0,ParenVal1$"

const (
	ParenVal0 ParenVal = 0
	ParenVal1 ParenVal = (1)
)

type EnumRHS Int32Enum // want EnumRHS:"^EnumRHS_A,EnumRHS_B$"

const (
	EnumRHS_A EnumRHS = iota
	EnumRHS_B
)

type WithMethod int // want WithMethod:"^WithMethodA,WithMethodB$"

const (
	WithMethodA WithMethod = 1
	WithMethodB WithMethod = 2
)

func (WithMethod) String() string { return "whatever" }
