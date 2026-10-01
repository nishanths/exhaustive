package exhaustive

import (
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"iter"
	"regexp"
	"slices"
	"strings"

	"github.com/nishanths/exhaustive/passes/enumerated"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "exhaustive",
	Doc:      "check that expression switch statements of enumerated types are exhaustive",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer, enumerated.Analyzer},
}

func init() {
	Analyzer.Flags.BoolVar(&fNeedEnforceDirective, "e", fNeedEnforceDirective, "check switch statement only if it has '//exhaustive:enforce' comment")
	Analyzer.Flags.BoolVar(&fDefaultEx, "d", fDefaultEx, "presence of default case makes switch statement exhaustive regardless of other cases")
	Analyzer.Flags.BoolVar(&fDefaultRequired, "defrequire", fDefaultRequired, "default case must always be present")
	Analyzer.Flags.StringVar(&fCheck, "check", fCheck, "specify the syntax tree elements that the analysis should check")
	Analyzer.Flags.BoolVar(&fCheckGenerated, "g", fCheckGenerated, "analyze generated files, too")
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

type synElement string

const (
	exprswitch synElement = "switch"
	mapliteral synElement = "mapliteral"
)

func parseSyntaxElements(arg string) ([]synElement, error) {
	if len(arg) == 0 {
		return nil, errors.New("empty string")
	}
	var ret []synElement
	for _, v := range strings.Split(arg, ",") {
		switch v := synElement(v); v {
		case exprswitch, mapliteral:
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

	opts := options{
		defaultCaseExhaustive: fDefaultEx,
		defaultCaseRequired:   fDefaultRequired,
		needEnforceDirective:  fNeedEnforceDirective,
		checkGenerated:        fCheckGenerated,
		includeTypePatterns:   fIncludeType.vals,
		excludeTypePatterns:   fExcludeType.vals,
		excludeConstPatterns:  fExcludeConst.vals,
	}
	for _, v := range elems {
		switch v {
		case exprswitch:
			checkSwitch(pass, &opts)
		case mapliteral:
			checkMapLiteral(pass, &opts)
		}
	}
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

func printNode(fset *token.FileSet, n ast.Node) string {
	var buf strings.Builder
	if err := printer.Fprint(&buf, fset, n); err != nil {
		// this should not happen for a valid program?
		// package x/tools/go/internal/astutil func Format
		// ignores the error in a similar situation.
		return "<?>"
	}
	return buf.String()
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
