package calc

import "github.com/nishanths/exhaustive/testdata/xxx/token"

func f(t token.Token) {
	switch t {
	case token.Add:
	case token.Subtract:
	case token.Multiply:
	default:
	}
}

var m = map[token.Token]string{
	token.Add:      "add",
	token.Subtract: "subtract",
	token.Multiply: "multiply",
}

// Testing instructions
//
// 	% go build ./cmd/exhaustive
// 	% ./exhaustive ./testdata/xxx/calc
//  %
