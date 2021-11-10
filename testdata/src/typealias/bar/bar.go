package bar

type T2 string // want T2:"^A,B,C,D,E,F$"

const (
    A T2 = "."
    B T2 = "-"
    C T2 = "+"
    D T2 = "*"
    E T2 = "&"
    F T2 = "|"
)
