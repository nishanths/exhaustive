package bar

type Phylum int // want Phylum:"^Chordata,Echinodermata,Mollusca,platyhelminthes$"

const (
	Chordata Phylum = iota
	Echinodermata
	Mollusca
	platyhelminthes
)

type IntWrapper int
