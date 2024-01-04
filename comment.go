package exhaustive

import (
	"go/ast"
	"go/token"
	"strings"
)

const (
	exhaustiveComment                 = "//exhaustive:"
	ignoreComment                     = "ignore"
	enforceComment                    = "enforce"
	ignoreDefaultCaseRequiredComment  = "ignore-default-case-required"
	enforceDefaultCaseRequiredComment = "enforce-default-case-required"
)

type directive int64

const (
	ignoreDirective = 1 << iota
	enforceDirective
	ignoreDefaultCaseRequiredDirective
	enforceDefaultCaseRequiredDirective
)

type directiveSet int64

func parseDirectives(commentGroups []*ast.CommentGroup) directiveSet {
	var out directiveSet
	for _, commentGroup := range commentGroups {
		for _, comment := range commentGroup.List {
			commentLine := comment.Text
			if !strings.HasPrefix(commentLine, exhaustiveComment) {
				continue
			}
			directive := commentLine[len(exhaustiveComment):]
			if whiteSpaceIndex := strings.IndexAny(directive, " \t"); whiteSpaceIndex != -1 {
				directive = directive[:whiteSpaceIndex]
			}
			switch directive {
			case ignoreComment:
				out |= ignoreDirective
			case enforceComment:
				out |= enforceDirective
			case ignoreDefaultCaseRequiredComment:
				out |= ignoreDefaultCaseRequiredDirective
			case enforceDefaultCaseRequiredComment:
				out |= enforceDefaultCaseRequiredDirective
			}
		}
	}
	return out
}

func (d directiveSet) has(directive directive) bool {
	return int64(d)&int64(directive) != 0
}

func fileCommentMap(fset *token.FileSet, file *ast.File) ast.CommentMap {
	return ast.NewCommentMap(fset, file, file.Comments)
}
