package exhaustive

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

func checkMapLiterals(pass *analysis.Pass, inspect *inspector.Inspector, comments map[*ast.File]ast.CommentMap) {
	// TODO
}
