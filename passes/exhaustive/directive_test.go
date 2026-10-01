package exhaustive

import (
	"go/ast"
	"reflect"
	"testing"
)

func TestParseDirective(t *testing.T) {
	type testcase struct {
		caption    string
		input      []*ast.CommentGroup
		directives map[directive]bool
		err        string
	}

	empty := make(map[directive]bool)
	testcases := []testcase{
		{"none1", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//foo:x"}}},
		}, empty, ""},
		{"none2", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//foo:enforce arg1 arg2 // xyz"}}},
		}, empty, ""},
		{"space after comment marker", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "// exhaustive:enforce"}}},
		}, empty, ""},
		{"commented out", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "// //exhaustive:enforce"}}},
		}, empty, ""},
		{"args", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enforce arg1 arg2"}}},
		}, nil, "args not allowed"},
		{"comment after directive", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:ignore // comment"}}},
		}, nil, "args not allowed"},
		{"invalid name", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enfocre"}}},
		}, nil, "invalid name \"enfocre\""},
		{"invalid name", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enforce-foo"}}},
		}, nil, "invalid name \"enforce-foo\""},
		{"conflict1", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enforce"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:ignore"}}},
		}, nil, "conflicting directives"},
		{"conflict2", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:ignore"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:enforce"}}},
		}, nil, "conflicting directives"},
		{"conflict3", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:defrequire=0"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:defrequire=1"}}},
		}, nil, "conflicting directives"},
		{"conflictalt1", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:ignore-default-case-required"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:enforce-default-case-required"}}},
		}, nil, "conflicting directives"},
		{"conflictalt2", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:defrequire=0"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:enforce-default-case-required"}}},
		}, nil, "conflicting directives"},
		{"typical", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:ignore"}}},
			{List: []*ast.Comment{{Text: "//bar:y"}}},
			{List: []*ast.Comment{{Text: "// comment"}}},
		}, map[directive]bool{dirIgnore: true}, ""},
		{"single1", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:ignore"}}},
		}, map[directive]bool{dirIgnore: true}, ""},
		{"single2", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enforce"}}},
		}, map[directive]bool{dirEnforce: true}, ""},
		{"single3", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:defrequire=0"}}},
		}, map[directive]bool{dirDefrequire: false}, ""},
		{"single4", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:defrequire=1"}}},
		}, map[directive]bool{dirDefrequire: true}, ""},
		{"singlealt1", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:ignore-default-case-required"}}},
		}, map[directive]bool{dirDefrequire: false}, ""},
		{"singlealt2", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enforce-default-case-required"}}},
		}, map[directive]bool{dirDefrequire: true}, ""},
		{"multi1", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enforce"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:defrequire=0"}}},
		}, map[directive]bool{dirEnforce: true, dirDefrequire: false}, ""},
		{"multi2", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:enforce"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:defrequire=1"}}},
		}, map[directive]bool{dirEnforce: true, dirDefrequire: true}, ""},
		{"multialt1", []*ast.CommentGroup{
			{List: []*ast.Comment{{Text: "//exhaustive:ignore"}}},
			{List: []*ast.Comment{{Text: "//exhaustive:ignore-default-case-required"}}},
		}, map[directive]bool{dirIgnore: true, dirDefrequire: false}, ""},
		{"misc1", []*ast.CommentGroup{
			// The target comment directive is
			// between comment groups, and it has
			// surrounding line comments.
			// There are multiple directives.
			{List: []*ast.Comment{{Text: "// hello, world"}}},
			{List: []*ast.Comment{{Text: "// comment"}, {Text: "//foo:x"}, {Text: "//exhaustive:enforce"}, {Text: "// comment"}}},
			{List: []*ast.Comment{{Text: "//bar:y"}, {Text: "//exhaustive:defrequire=1"}, {Text: "// comment"}}},
		}, map[directive]bool{dirEnforce: true, dirDefrequire: true}, ""},
	}

	for _, tt := range testcases {
		t.Run(tt.caption, func(t *testing.T) {
			got, err := parseDirectives(tt.input)
			switch {
			case err == nil && tt.err != "":
				t.Errorf("got nil error, want %q", tt.err)
			case err != nil && tt.err == "":
				t.Errorf("got error %q, want nil", err)
			case err != nil && err.Error() != tt.err:
				t.Errorf("got error %q, want %q", err, tt.err)
			case !reflect.DeepEqual(got, tt.directives):
				t.Errorf("got: %v, want: %v", got, tt.directives)
			}
		})
	}
}
