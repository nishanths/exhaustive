package exhaustive

import (
	"go/ast"
)

func isGeneratedFile(file *ast.File) bool {
	return ast.IsGenerated(file)
}
