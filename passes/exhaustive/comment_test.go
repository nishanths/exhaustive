// Partly based on code from package go/ast file
// commentmap_test.go in the standard Go library.
// See LICENSE file.

package exhaustive

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"golang.org/x/tools/go/ast/inspector"
)

var testcases = []struct {
	filename string
	nodeType ast.Node
	f        func(inspector.Cursor, *token.FileSet, []*ast.CommentGroup) []*ast.CommentGroup
	res      map[string]string
}{
	{"comment1.go", (*ast.SwitchStmt)(nil), switchComments, map[string]string{
		"  6:  2 *ast.SwitchStmt": "",
		" 16:  2 *ast.SwitchStmt": "a0\n",
		" 20:  2 *ast.SwitchStmt": "a1\n",
		" 30:  2 *ast.SwitchStmt": "a2\n\na2+continued\n",
		" 34:  2 *ast.SwitchStmt": "a3\n",
		" 39:  3 *ast.SwitchStmt": "a4\n",
		" 45:  2 *ast.SwitchStmt": "",
		" 48:  3 *ast.SwitchStmt": "a5\n",
		" 59:  2 *ast.SwitchStmt": "",
		" 63:  2 *ast.SwitchStmt": "a6\n",
		" 71:  2 *ast.SwitchStmt": "a7\n",
		" 76: 11 *ast.SwitchStmt": " a8\n",
		" 82: 22 *ast.SwitchStmt": "a9\na9+continued\n a9.1\n\t * a9.1+continued\n",
		" 86: 13 *ast.SwitchStmt": "a10\n a10.1\n",
		" 90:  2 *ast.SwitchStmt": "",
		" 96:  2 *ast.SwitchStmt": "",
		"103:  2 *ast.SwitchStmt": "",
		"112:  2 *ast.SwitchStmt": "a11\n",
		"112: 13 *ast.SwitchStmt": "",
	}},
	{"comment2.go", (*ast.CompositeLit)(nil), compositeLitComments, map[string]string{
		" 17:  7 *ast.CompositeLit": "",
		" 18:  7 *ast.CompositeLit": "",
		" 23:  7 *ast.CompositeLit": "a0\n",
		" 24:  7 *ast.CompositeLit": "a0\n",
		" 27: 10 *ast.CompositeLit": "",
		" 30: 10 *ast.CompositeLit": "a1\n",
		" 35:  7 *ast.CompositeLit": "",
		" 36:  7 *ast.CompositeLit": "",
		" 40:  7 *ast.CompositeLit": "a2\n",
		" 41:  7 *ast.CompositeLit": "",
		" 46:  6 *ast.CompositeLit": "",
		" 49:  6 *ast.CompositeLit": "a3\n",
		" 53: 10 *ast.CompositeLit": "",
		" 56: 10 *ast.CompositeLit": "a4\n",
		" 62: 33 *ast.CompositeLit": "",
		" 65:  6 *ast.CompositeLit": "a5\n",
		" 68:  7 *ast.CompositeLit": "a6\n",
		" 71: 13 *ast.CompositeLit": "a7\n",
		" 71: 23 *ast.CompositeLit": "a7\n",
		" 76:  9 *ast.CompositeLit": "",
		" 76: 19 *ast.CompositeLit": "",
		" 82: 10 *ast.CompositeLit": "",
	}},
	{"comment3.go", (*ast.CompositeLit)(nil), compositeLitComments, map[string]string{
		"  6:  6 *ast.CompositeLit": "",
		" 14:  6 *ast.CompositeLit": "b1\nb1.1\nb0\n",
		" 16:  6 *ast.CompositeLit": "b2\nb0\n",
		" 21:  6 *ast.CompositeLit": "b0\n",
		" 25:  6 *ast.CompositeLit": "b3\nb0\n",
		" 28:  5 *ast.CompositeLit": "b0\n",
		" 38:  6 *ast.CompositeLit": "a0\n",
		" 41:  6 *ast.CompositeLit": "a1\n",
		" 50:  6 *ast.CompositeLit": "a2\n\na2+continued\n",
		" 53:  6 *ast.CompositeLit": "a3\n",
		" 58:  8 *ast.CompositeLit": "a4\n",
		" 64:  6 *ast.CompositeLit": "",
		" 67:  8 *ast.CompositeLit": "a5\n",
		" 78:  6 *ast.CompositeLit": "",
		" 81:  6 *ast.CompositeLit": "a6\n",
		" 88:  6 *ast.CompositeLit": "a7\n",
		" 92: 15 *ast.CompositeLit": " a8\n",
		" 97: 26 *ast.CompositeLit": "a9\na9+continued\n a9.1\n\t * a9.1+continued\n",
		"100: 17 *ast.CompositeLit": "a10\n a10.1\n",
		"103:  6 *ast.CompositeLit": "",
		"108:  6 *ast.CompositeLit": "",
		"114:  6 *ast.CompositeLit": "",
		"122:  6 *ast.CompositeLit": "a11\n",
		"122: 21 *ast.CompositeLit": "",
	}},
}

func ctext(list []*ast.CommentGroup) string {
	var buf strings.Builder
	for _, g := range list {
		buf.WriteString(g.Text())
	}
	return buf.String()
}

func TestComment(t *testing.T) {
	for _, tt := range testcases {
		t.Run(tt.filename, func(t *testing.T) {
			src, err := os.ReadFile(filepath.Join("testdata", tt.filename))
			if err != nil {
				t.Fatal(err)
			}

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}

			in := inspector.New([]*ast.File{f})

			got := make(map[string]string)
			for c := range in.Root().Preorder(tt.nodeType) {
				key := fmt.Sprintf("%3d:%3d %T", fset.Position(c.Node().Pos()).Line, fset.Position(c.Node().Pos()).Column, c.Node())
				fc := tt.f(c, fset, f.Comments)
				got[key] = ctext(fc)
			}
			for k := range got {
				if got[k] != tt.res[k] {
					t.Errorf("%s: got %q; want %q", k, got[k], tt.res[k])
				}
			}
			for k := range tt.res {
				if got[k] != tt.res[k] {
					t.Errorf("%s: got %q; want %q", k, got[k], tt.res[k])
				}
			}

			// support code to update test:
			// set genMap to true to generate res map
			const genMap = false
			if genMap {
				var out []string
				for c := range in.Root().Preorder(tt.nodeType) {
					fc := tt.f(c, fset, f.Comments)
					out = append(out, fmt.Sprintf("\t\"%3d:%3d %T\":\t%q,", fset.Position(c.Node().Pos()).Line, fset.Position(c.Node().Pos()).Column, c.Node(), ctext(fc)))
				}
				slices.Sort(out)
				for _, s := range out {
					fmt.Println(s)
				}
			}
		})
	}
}
