package exhaustive

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"iter"
	"regexp"
	"slices"
	"strings"

	"github.com/nishanths/exhaustive/passes/enumerated"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "exhaustive",
	Doc:      "check that expression switch statements of enumerated types are exhaustive",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer, enumerated.Analyzer},
}

func init() {
	Analyzer.Flags.BoolVar(&fNeedEnforceDirective, "e", fNeedEnforceDirective, "check a switch statement only if it has '//exhaustive:enforce' comment")
	Analyzer.Flags.BoolVar(&fDefaultEx, "d", fDefaultEx, "including a default case makes a switch statement exhaustive")
	Analyzer.Flags.BoolVar(&fDefaultRequired, "defrequire", fDefaultRequired, "default case must always be present")
	Analyzer.Flags.StringVar(&fCheck, "check", fCheck, "specify elements in the syntax tree that the analysis should check")
	Analyzer.Flags.BoolVar(&fCheckGenerated, "g", fCheckGenerated, "additionally analyze generated files")
	Analyzer.Flags.Var(&fExcludeType, "typeignore", "switch statements in which the type name is matched by `regexp` are not checked")
	Analyzer.Flags.Var(&fIncludeType, "typeonly", "only switch statements in which the type name is matched by `regexp` are checked")
	Analyzer.Flags.Var(&fExcludeConst, "constignore", "constant names matched by `regexp` do not have to be included in case expressions")
}

var (
	fNeedEnforceDirective = false
	fDefaultEx            = false
	fDefaultRequired      = false
	fCheck                = string(exprswitch)
	fCheckGenerated       = false
	fIncludeType          = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fExcludeType          = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fExcludeConst         = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
)

func resetFlags() {
	fNeedEnforceDirective = false
	fDefaultEx = false
	fDefaultRequired = false
	fCheck = string(exprswitch)
	fCheckGenerated = false
	fIncludeType = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fExcludeType = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fExcludeConst = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
}

// repeatFlag defines a flag that may be repeated to specify
// multiple values.
// Note: The default behavior of package flag when a flag is
// repeated is to use the last value.
type repeatFlag[T any] struct {
	raw  []string
	vals []T
	set  func(string) (T, error)
}

func (r *repeatFlag[T]) String() string {
	if r == nil {
		return fmt.Sprintf("%v", []string(nil)) // formats as []
	}
	return fmt.Sprintf("%v", []string(r.raw))
}

func (r *repeatFlag[T]) Set(arg string) error {
	v, err := r.set(arg)
	if err != nil {
		return err
	}
	r.raw = append(r.raw, arg)
	r.vals = append(r.vals, v)
	return nil
}

type syntaxElement string

const (
	exprswitch syntaxElement = "switch"
	mapliteral syntaxElement = "mapliteral"
	typeswitch syntaxElement = "typeswitch"
)

func parseSyntaxElements(arg string) ([]syntaxElement, error) {
	if len(arg) == 0 {
		return nil, errors.New("empty string")
	}
	var ret []syntaxElement
	for _, v := range strings.Split(arg, ",") {
		switch v := syntaxElement(v); v {
		case exprswitch, mapliteral, typeswitch:
			if !slices.Contains(ret, v) {
				ret = append(ret, v)
			}
		default:
			return nil, fmt.Errorf("unknown element %q", v)
		}
	}
	return ret, nil
}

type options struct {
	defaultCaseExhaustive bool
	defaultCaseRequired   bool
	needEnforceDirective  bool
	checkGenerated        bool
	includeTypePatterns   []*regexp.Regexp
	excludeTypePatterns   []*regexp.Regexp
	excludeConstPatterns  []*regexp.Regexp
}

func run(pass *analysis.Pass) (any, error) {
	elems, err := parseSyntaxElements(fCheck)
	if err != nil {
		return nil, fmt.Errorf("invalid value for flag -check: %s", err)
	}

	opts := &options{
		defaultCaseExhaustive: fDefaultEx,
		defaultCaseRequired:   fDefaultRequired,
		needEnforceDirective:  fNeedEnforceDirective,
		checkGenerated:        fCheckGenerated,
		includeTypePatterns:   fIncludeType.vals,
		excludeTypePatterns:   fExcludeType.vals,
		excludeConstPatterns:  fExcludeConst.vals,
	}

	var nodes []ast.Node
	for _, v := range elems {
		switch v {
		case exprswitch:
			nodes = append(nodes, (*ast.SwitchStmt)(nil))
		case mapliteral:
			nodes = append(nodes, (*ast.CompositeLit)(nil))
		case typeswitch:
			nodes = append(nodes, (*ast.TypeSwitchStmt)(nil))
		}
	}

	in := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	commentsByFile := make(map[*ast.File]ast.CommentMap)

	prepare := func(c inspector.Cursor) (ast.CommentMap, bool) {
		file := head(c.Enclosing((*ast.File)(nil))).Node().(*ast.File)
		if !opts.checkGenerated && ast.IsGenerated(file) {
			return nil, false
		}
		comments, ok := commentsByFile[file]
		if !ok {
			comments = ast.NewCommentMap(pass.Fset, file, file.Comments)
			commentsByFile[file] = comments
		}
		return comments, true
	}

	in.Root().Inspect(nodes, func(c inspector.Cursor) (descend bool) {
		switch n := c.Node().(type) {
		case *ast.SwitchStmt:
			comments, ok := prepare(c)
			if !ok {
				return true
			}
			checkSwitch(pass, n, comments, opts)
		case *ast.CompositeLit:
			comments, ok := prepare(c)
			if !ok {
				return true
			}
			checkMapLiteral(pass, c, n, comments, opts)
		case *ast.TypeSwitchStmt:
			comments, ok := prepare(c)
			if !ok {
				return true
			}
			checkTypeSwitch(pass, n, comments, opts)
		default:
			panic(fmt.Sprintf("internal error: unexpected type %T", n))
		}
		return true
	})

	return nil, nil
}

type directive int

const (
	dirIgnore directive = iota
	dirEnforce
	dirDefrequire
)

// A comment directive is matched by the regular
// expression: "//[a-z0-9]+:[a-z0-9]".
// Note that a line of comment text can have at most one
// directive, which must start at the very beginning.
//
// References:
// package go/ast func isDirective
// package x/tools@v0.50.0/internal/astutil type Directive and related functions
//
// The parse routines in x/tools/internal/astutil allow a
// separate regular comment after the optional args at the end
// of the comment directive.
//
//	//tool:name [arg ...] // regular comment
//
// On the other hand the parse routines in package go/ast
// in the standard library consider the args to be everything
// after the directive name.
func parseDirectives(groups []*ast.CommentGroup) (map[directive]bool, error) {
	var (
		errConflict = errors.New("conflicting directives")
		result      = make(map[directive]bool)
	)
	for _, g := range groups {
		if g != nil { // Package x/tools/internal/astutil employs a similar guard.
			for _, c := range g.List {
				if d, isDirective := ast.ParseDirective(c.Pos(), c.Text); isDirective && d.Tool == "exhaustive" {
					if len(d.Args) != 0 {
						return nil, errors.New("args not allowed")
					}
					switch d.Name {
					case "ignore":
						result[dirIgnore] = true
					case "enforce":
						result[dirEnforce] = true
					case "defrequire=0", "ignore-default-case-required":
						// Note: The latter name is supported but deprecated.
						if v, ok := result[dirDefrequire]; ok && v != false {
							return nil, errConflict
						}
						result[dirDefrequire] = false
					case "defrequire=1", "enforce-default-case-required":
						// Ditto note.
						if v, ok := result[dirDefrequire]; ok && v != true {
							return nil, errConflict
						}
						result[dirDefrequire] = true
					default:
						return nil, fmt.Errorf("invalid name %q", d.Name)
					}
				}
			}
		}
	}
	if result[dirIgnore] && result[dirEnforce] {
		return nil, errConflict
	}
	return result, nil
}

func isPkgLevel(obj types.Object) bool {
	return obj.Pkg() != nil && obj.Pkg().Scope().Lookup(obj.Name()) == obj
}

func fullname(obj types.Object) string {
	if !isPkgLevel(obj) {
		panic("not a package-level object")
	}
	var b strings.Builder
	b.WriteString(obj.Pkg().Path())
	b.WriteString(".")
	b.WriteString(obj.Name())
	return b.String()
}

func matchAny(rs []*regexp.Regexp, s string) bool {
	for _, r := range rs {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

func head[T any](seq iter.Seq[T]) T {
	for t := range seq {
		return t
	}
	panic("empty sequence")
}

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}
