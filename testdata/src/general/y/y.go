package bar

type Phylum uint8 // want Phylum:"^Chordata,Echinodermata,Mollusca,platyhelminthes$"

const (
	Chordata Phylum = iota
	Echinodermata
	Mollusca
	platyhelminthes
)

type IntWrapper int

type Uppercase int // want Uppercase:"^ReallyExported$"
const ReallyExported Uppercase = 1

func f() {
	type AliasForUppercase = Uppercase
	const NotReallyExported AliasForUppercase = 2 // not exported, and in fact not even a member of enum type Uppercase
}
