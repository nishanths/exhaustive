// want package:"AcrossBlocksDeclsFiles:Here,Separate,There; ByteEnum:ByteA; Int32Enum:Int32A,Int32B; IotaEnum:IotaA,ItoaB; RepeatedValue:RepeatedValueA,RepeatedValueB; RuneEnum:RuneA; StringEnum:StringA,StringB,StringC; UIntEnum:UIntA,UIntB; UnexportedMembers:unexportedMembersA,unexportedMembersB; VarMembers:VarMemberA"

package enumvariants

type VarMembers int

var (
	VarMemberA VarMembers
)

type IotaEnum uint8

const (
	IotaA IotaEnum = iota << 1
	ItoaB
)

type MemberlessEnum int

type RepeatedValue int

const (
	RepeatedValueA RepeatedValue = 1
	RepeatedValueB RepeatedValue = 1
)

type AcrossBlocksDeclsFiles int

const (
	Here AcrossBlocksDeclsFiles = 0
)

type UnexportedMembers int

const (
	unexportedMembersA UnexportedMembers = 1
	unexportedMembersB UnexportedMembers = 2
)

var Separate AcrossBlocksDeclsFiles = 1

type NonBasicType S

type S struct{ F int }

var (
	SA NonBasicType = NonBasicType{F: 1}
	SB NonBasicType = NonBasicType{F: 2}
)

type NonTopLevel uint

func _nonTopLevel() {
	const (
		A NonTopLevel = 0
		B NonTopLevel = 1
	)
}
