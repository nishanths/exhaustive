package bar

type T2 string // want T2:"^A,AA,B,C,D,E,F,I$"

const (
    A  T2 = "."
    AA    = A
    B  T2 = "-"
    C  T2 = "+"
    D  T2 = "*"
    E  T2 = "&"
    F  T2 = "|"
    I  T2 = "<"
)
