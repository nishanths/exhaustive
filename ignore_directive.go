package exhaustive

import (
	"go/ast"
	"strings"
)

// IgnoreDirectivePrefix is used to exclude checking of specific switch statements.
// See package comment for details.
const IgnoreDirectivePrefix = "//exhaustive:ignore"

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
		if strings.HasPrefix(c.Text, IgnoreDirectivePrefix) {
			return true
		}
	}
	return false
}
