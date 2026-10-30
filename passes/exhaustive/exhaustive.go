package exhaustive

import (
	"errors"
	"fmt"
	"go/ast"
	"regexp"
	"slices"
	"strings"

	"github.com/nishanths/exhaustive/internal/passes/finder"
	"github.com/nishanths/exhaustive/passes/enumerated"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// controls whether reported diagnostics use short messages
// to simplify regexp matching in tests;
// must be set in tests only.
var testingDiagnosticShort = false

// Analyzer implements the exhaustive static analysis.
//
// The set of analyzers specified in the Requires field may not
// all be actually required.
// The set of analyzers that are actually required depends on the
// runtime value of the -check flag.
// Programs that want to be efficient should edit the Requires
// field before the analysis begins. See also: func [Requires].
var Analyzer = &analysis.Analyzer{
	Name: "exhaustive",
	Doc:  "check that expression switches and type switches are exhaustive",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,    // always required
		enumerated.Analyzer, // required if -check argument specifies either 'switch' or 'mapliteral'
		finder.Analyzer,     // required if -check argument specifies 'typeswitch'
	},
}

// Requires returns the set of analyzers required by [Analyzer]
// for a given value of the -check flag.
// The result is a subset of Analyzer.Requires.
func Requires(check string) ([]*analysis.Analyzer, error) {
	require := map[*analysis.Analyzer]bool{
		inspect.Analyzer:    true,
		enumerated.Analyzer: false,
		finder.Analyzer:     false,
	}

	elems, err := parseSyntaxElements(check)
	if err != nil {
		return nil, fmt.Errorf("invalid value for -check: %s", err)
	}
	for _, elem := range elems {
		switch elem {
		case exprswitch, mapliteral:
			require[enumerated.Analyzer] = true
		case typeswitch:
			require[finder.Analyzer] = true
		}
	}

	var res []*analysis.Analyzer // order should not matter
	for k, v := range require {
		if v {
			res = append(res, k)
		}
	}
	return res, nil
}

func init() {
	Analyzer.Flags.BoolVar(&fRequireNilCase, "casenil", fRequireNilCase, "type switches must include a nil case to be exhaustive")
	Analyzer.Flags.StringVar(&fCheck, "check", fCheck, "comma-separated list of elements to analyze; available values: switch, typeswitch, mapliteral")
	Analyzer.Flags.Var(&fExcludeConst, "constignore", "constant names matched by `regexp` do not have to be listed in the cases of an expression switch")
	Analyzer.Flags.BoolVar(&fDefaultEx, "d", fDefaultEx, "default case in a switch statement satisifies all missing cases")
	Analyzer.Flags.BoolVar(&fRequireDefaultCase, "defrequire", fRequireDefaultCase, "check that analyzed expression switches always include a default case")
	Analyzer.Flags.BoolVar(&fCheckEnforceOnly, "e", fCheckEnforceOnly, "analyze a switch statement only if it has an '//exhaustive:enforce' comment directive")
	Analyzer.Flags.BoolVar(&fCheckGenerated, "g", fCheckGenerated, "additionally analyze generated files")
	Analyzer.Flags.Var(&fExcludeType, "typeignore", "do not analyze a switch statement if the switch expression type name is matched by `regexp`")
	Analyzer.Flags.Var(&fIncludeType, "typeonly", "analyze a switch statement only if the switch expression type name is matched by `regexp`")
}

var (
	fRequireNilCase     = false
	fCheck              = string(exprswitch)
	fExcludeConst       = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fDefaultEx          = false
	fRequireDefaultCase = false
	fCheckEnforceOnly   = false
	fCheckGenerated     = false
	fExcludeType        = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fIncludeType        = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
)

func resetFlags() {
	fRequireNilCase = false
	fCheck = string(exprswitch)
	fExcludeConst = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fDefaultEx = false
	fRequireDefaultCase = false
	fCheckEnforceOnly = false
	fCheckGenerated = false
	fExcludeType = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
	fIncludeType = repeatFlag[*regexp.Regexp]{set: regexp.Compile}
}

// repeatFlag defines a flag that may be repeated to specify
// multiple values.
// Note: The default behavior of package flag when a flag is
// repeated is to use the last value.
type repeatFlag[T any] struct {
	raw    []string
	values []T
	set    func(string) (T, error)
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
	r.values = append(r.values, v)
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
	requireDefaultCase    bool
	requireNilCase        bool
	checkEnforceOnly      bool
	checkGenerated        bool
	includeTypePatterns   []*regexp.Regexp
	excludeTypePatterns   []*regexp.Regexp
	excludeConstPatterns  []*regexp.Regexp
}

func run(pass *analysis.Pass) (any, error) {
	elems, err := parseSyntaxElements(fCheck)
	if err != nil {
		return nil, fmt.Errorf("invalid value for -check: %s", err)
	}

	opts := &options{
		defaultCaseExhaustive: fDefaultEx,
		requireDefaultCase:    fRequireDefaultCase,
		requireNilCase:        fRequireNilCase,
		checkEnforceOnly:      fCheckEnforceOnly,
		checkGenerated:        fCheckGenerated,
		includeTypePatterns:   fIncludeType.values,
		excludeTypePatterns:   fExcludeType.values,
		excludeConstPatterns:  fExcludeConst.values,
	}

	var nodes []ast.Node
	for _, elem := range elems {
		switch elem {
		case exprswitch:
			nodes = append(nodes, (*ast.SwitchStmt)(nil))
		case mapliteral:
			nodes = append(nodes, (*ast.CompositeLit)(nil))
		case typeswitch:
			nodes = append(nodes, (*ast.TypeSwitchStmt)(nil))
		}
	}

	in := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	typeSwitchCh := new(typeSwitchState)

	in.Root().Inspect([]ast.Node{(*ast.File)(nil)}, func(c inspector.Cursor) (descend bool) {
		file := c.Node().(*ast.File)
		if !opts.checkGenerated && ast.IsGenerated(file) {
			return false
		}

		c.Inspect(nodes, func(c inspector.Cursor) (descend bool) {
			switch n := c.Node().(type) {
			case *ast.SwitchStmt:
				comments := switchComments(c, pass.Fset, file.Comments)
				checkExprSwitch(pass, n, comments, opts)
			case *ast.CompositeLit:
				comments := compositeLitComments(c, pass.Fset, file.Comments)
				checkMapLit(pass, n, comments, opts)
			case *ast.TypeSwitchStmt:
				comments := switchComments(c, pass.Fset, file.Comments)
				checkTypeSwitch(pass, n, comments, typeSwitchCh, opts)
			default:
				panic(fmt.Sprintf("internal error: unexpected node of type %T", n))
			}
			return true
		})

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
						if v, ok := result[dirDefrequire]; ok && v {
							return nil, errConflict
						}
						result[dirDefrequire] = false
					case "defrequire=1", "enforce-default-case-required":
						// Ditto note.
						if v, ok := result[dirDefrequire]; ok && !v {
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

func matchAny(rs []*regexp.Regexp, s string) bool {
	for _, r := range rs {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}
