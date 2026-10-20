package exhaustive

import (
	"go/ast"
	"go/token"
	"sort"

	"golang.org/x/tools/go/ast/inspector"
)

func switchComments(c inspector.Cursor, fset *token.FileSet, groups []*ast.CommentGroup) []*ast.CommentGroup {
	return nodeComments(c, fset, groups)
}

func compositeLitComments(c inspector.Cursor, fset *token.FileSet, groups []*ast.CommentGroup) []*ast.CommentGroup {
	var ret []*ast.CommentGroup
	seen := make(map[*ast.CommentGroup]struct{})
	add := func(gs ...*ast.CommentGroup) {
		for _, g := range gs {
			if _, ok := seen[g]; ok {
				continue
			}
			seen[g] = struct{}{}
			ret = append(ret, g)
		}
	}

loop:
	for c := range c.Enclosing((*ast.AssignStmt)(nil), (*ast.ValueSpec)(nil), (*ast.GenDecl)(nil)) {
		switch n := c.Node().(type) {
		case *ast.AssignStmt:
			add(nodeComments(c, fset, groups)...)
			break loop
		case *ast.ValueSpec:
			add(n.Doc)
			add(n.Comment)
			add(nodeComments(c, fset, groups)...)
			// note: do not break the loop here; also need to
			// see the parent *ast.GenDecl.
		case *ast.GenDecl:
			add(n.Doc)
			add(nodeComments(c, fset, groups)...)
			break loop
		}
	}

	return ret
}

// nodeComments returns the comments for the node specified by c.
// The node must be a ast.Decl, ast.Spec, or ast.Stmt node.
// The node must be from a syntactically correct and full Go source file.
// The input comment groups must be in lexical order.
func nodeComments(c inspector.Cursor, fset *token.FileSet, groups []*ast.CommentGroup) []*ast.CommentGroup {
	// The comment group must start after the line of the
	// previous node.
	var afterLine int
	if v, ok := c.PrevSibling(); ok {
		afterLine = fset.Position(v.Node().End()).Line
	} else {
		afterLine = fset.Position(c.Parent().Node().Pos()).Line
	}
	idx := sort.Search(len(groups), func(i int) bool { return fset.Position(groups[i].Pos()).Line > afterLine })
	groups = groups[idx:]

	// The comment group must end before the target node.
	idx = sort.Search(len(groups), func(i int) bool { return groups[i].End() > c.Node().Pos() })
	groups = groups[:idx]

	// The comment group must end either on the line
	// immediately before the target node or on the same line
	// as the target node.
	// In other words, there must be no empty lines between
	// the comment group and the target node.
	cLine := fset.Position(c.Node().Pos()).Line
	idx = sort.Search(len(groups), func(i int) bool { return fset.Position(groups[i].End()).Line >= cLine-1 })
	groups = groups[idx:]

	return groups
}
