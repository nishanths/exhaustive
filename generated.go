package exhaustive

import (
	"bytes"
	"go/ast"
	"strings"
)

// Adapated from https://gotools.org/dmitri.shuralyov.com/go/generated

func isGeneratedFile(file *ast.File) bool {
	buf := bytes.NewBufferString("") // shared buffer, reset each loop

	for _, c := range file.Comments {
		buf.Reset()
		for _, cc := range c.List {
			s := cc.Text // "\n" already removed (see doc comment)
			if !isSlashSlashStyleComment(s) {
				continue
			}
			if len(s) >= 1 && s[len(s)-1] == '\r' {
				s = s[:len(s)-1] // Trim "\r".
			}
			if containsGeneratedComment(s) {
				return true
			}
		}
	}

	return false
}

func isSlashSlashStyleComment(s string) bool {
	return strings.HasPrefix(s, "//")
}

func containsGeneratedComment(s string) bool {
	return len(s) >= len(genCommentPrefix)+len(genCommentSuffix) &&
		strings.HasPrefix(s, genCommentPrefix) &&
		strings.HasSuffix(s, genCommentSuffix)
}

var (
	genCommentPrefix = "// Code generated "
	genCommentSuffix = " DO NOT EDIT."
)
