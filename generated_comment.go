package exhaustive

import (
	"go/ast"
	"strings"
)

// Adapted from https://gotools.org/dmitri.shuralyov.com/go/generated

func isGeneratedFile(file *ast.File) bool {
	for _, c := range file.Comments {
		for _, cc := range c.List {
			s := cc.Text // "\n" already removed (see doc comment)
			if len(s) >= 1 && s[len(s)-1] == '\r' {
				s = s[:len(s)-1] // Trim "\r".
			}
			if isGeneratedFileComment(s) {
				return true
			}
		}
	}

	return false
}

func isGeneratedFileComment(s string) bool {
	return strings.HasPrefix(s, genCommentPrefix) &&
		strings.HasSuffix(s, genCommentSuffix)
}

const (
	genCommentPrefix = "// Code generated "
	genCommentSuffix = " DO NOT EDIT."
)
