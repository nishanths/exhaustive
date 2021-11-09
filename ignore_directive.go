package exhaustive

import (
	"go/ast"
	"strings"
)

// ignoreDirective is used to exclude checking of specific switch statements.
// See package comment for details.
const ignoreDirective = "//exhaustive:ignore"

func containsIgnoreDirective(groups []*ast.CommentGroup) bool {
	for _, group := range groups {
		if containsIgnoreDirectiveComments(group.List) {
			return true
		}
	}
	return false
}

func containsIgnoreDirectiveComments(comments []*ast.Comment) bool {
	for _, c := range comments {
		if strings.HasPrefix(c.Text, ignoreDirective) {
			return true
		}
	}
	return false
}
