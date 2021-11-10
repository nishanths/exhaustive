package exhaustive

import (
	"go/ast"
	"regexp"
)

// Spec
// http://golang.org/s/generatedcode
//
// To convey to humans and machine tools that code is generated, generated
// source should have a line that matches the following regular expression (in
// Go syntax):
//
//   ^// Code generated .* DO NOT EDIT\.$
//
// This line must appear before the first non-comment, non-blank
// text in the file.

func isGeneratedFile(file *ast.File) bool {
	// NOTE: file.Comments includes file.Doc as well, so no need
	// to separately check file.Doc.

	for _, c := range file.Comments {
		for _, cc := range c.List {
			// This check is intended to handle "must appear before the
			// first non-comment, non-blank text in the file".
			// TODO: Is this check fully correct?
			if c.Pos() > file.Package {
				return false
			}
			// According to the docs:
			// '\r' has been removed.
			// '\n' has been removed for //-style comments, which is what we care about.
			// Also manually verified.
			if isGeneratedFileComment(cc.Text) {
				return true
			}
		}
	}

	return false
}

func isGeneratedFileComment(s string) bool {
	return generatedCodeRx.MatchString(s)
}

var generatedCodeRx = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)
