package exhaustive

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"regexp"

	"golang.org/x/tools/go/analysis"
)

// nodeVisitor is like the visitor function used by Inspector.WithStack,
// except that it returns an additional value: a short description of
// the result of this node visit.
//
// The result is typically useful in debugging or in unit tests to check
// that the nodeVisitor function took the expected code path.
type nodeVisitor func(n ast.Node, push bool, stack []ast.Node) (proceed bool, result string)

// toVisitor converts the nodeVisitor to a function suitable for use
// with Inspector.WithStack.
func toVisitor(v nodeVisitor) func(ast.Node, bool, []ast.Node) bool {
	return func(node ast.Node, push bool, stack []ast.Node) bool {
		proceed, _ := v(node, push, stack)
		return proceed
	}
}

// Result values returned by node visitors.
const (
	resultEmptyMapLiteral = "empty map literal"
	resultNotMapLiteral   = "not map literal"
	resultKeyNilPkg       = "nil map key package"
	resultKeyNotEnum      = "not all map key type terms are known enum types"

	resultNoSwitchTag = "no switch tag"
	resultTagNotValue = "switch tag not value type"
	resultTagNilPkg   = "nil switch tag package"
	resultTagNotEnum  = "not all switch tag terms are known enum types"

	resultNotPush              = "not push"
	resultGeneratedFile        = "generated file"
	resultIgnoreComment        = "has ignore comment"
	resultNoEnforceComment     = "has no enforce comment"
	resultEnumMembersAccounted = "required enum members accounted for"
	resultDefaultCaseSuffices  = "default case satisfies exhaustiveness"
	resultReportedDiagnostic   = "reported diagnostic"
)

// switchChecker returns a node visitor that checks exhaustiveness of
// enum switch statements for the supplied pass, and reports
// diagnostics. The node visitor expects only *ast.SwitchStmt nodes.
func switchChecker(pass *analysis.Pass, cfg switchConfig, generated boolCache, comments commentCache) nodeVisitor {
	return func(n ast.Node, push bool, stack []ast.Node) (bool, string) {
		if !push {
			// The proceed return value should not matter; it is ignored by
			// inspector package for pop calls.
			// Nevertheless, return true to be on the safe side for the future.
			return true, resultNotPush
		}

		file := stack[0].(*ast.File)

		if !cfg.checkGenerated && generated.get(file) {
			// Don't check this file.
			// Return false because the children nodes of node `n` don't have to be checked.
			return false, resultGeneratedFile
		}

		sw := n.(*ast.SwitchStmt)

		switchComments := comments.get(pass.Fset, file)[sw]
		if !cfg.explicit && hasComment(switchComments, ignoreComment) {
			// Skip checking of this switch statement due to ignore
			// comment. Still return true because there may be nested
			// switch statements that are not to be ignored.
			return true, resultIgnoreComment
		}
		if cfg.explicit && !hasComment(switchComments, enforceComment) {
			// Skip checking of this switch statement due to missing
			// enforce comment.
			return true, resultNoEnforceComment
		}

		if sw.Tag == nil {
			return true, resultNoSwitchTag
		}

		t := pass.TypesInfo.Types[sw.Tag]
		if !t.IsValue() {
			return true, resultTagNotValue
		}

		es, all := composingEnumTypes(pass, t.Type)
		if !all {
			log.Printf("%#v", es)
			return true, resultTagNotEnum // TODO(nishanths) could be other reasons e.g. nil pkg, make this more generic
		}

		var checkl checklist
		checkl.ignore(cfg.ignoreEnumMembers)

		for _, e := range es {
			checkl.add(e.et, e.em, pass.Pkg == e.et.Pkg())
		}

		def := analyzeSwitchClauses(sw, pass.TypesInfo, checkl.found)
		if len(checkl.remaining()) == 0 {
			// All enum members accounted for.
			// Nothing to report.
			return true, resultEnumMembersAccounted
		}
		if def && cfg.defaultSignifiesExhaustive {
			// Though enum members are not accounted for, the
			// existence of the default case signifies
			// exhaustiveness.  So don't report.
			return true, resultDefaultCaseSuffices
		}
		pass.Report(makeSwitchDiagnostic(sw, toTypes(es), checkl.remaining()))
		return true, resultReportedDiagnostic
	}
}

func toTypes(es []typeAndMembers) []enumType {
	out := make([]enumType, len(es))
	for i := range es {
		out[i] = es[i].et
	}
	return out
}

// switchConfig is configuration for switchChecker.
type switchConfig struct {
	explicit                   bool
	defaultSignifiesExhaustive bool
	checkGenerated             bool
	ignoreEnumMembers          *regexp.Regexp // can be nil
}

func isDefaultCase(c *ast.CaseClause) bool {
	return c.List == nil // see doc comment on List field
}

// analyzeSwitchClauses analyzes the clauses in the supplied switch
// statement. The info param typically is pass.TypesInfo. The each
// function is called for each enum member name found in the switch
// statement. The hasDefaultCase return value indicates whether the
// switch statement has a default clause.
func analyzeSwitchClauses(sw *ast.SwitchStmt, info *types.Info, each func(val constantValue)) (hasDefaultCase bool) {
	for _, stmt := range sw.Body.List {
		caseCl := stmt.(*ast.CaseClause)
		if isDefaultCase(caseCl) {
			hasDefaultCase = true
			continue
		}
		for _, expr := range caseCl.List {
			if val, ok := exprConstVal(expr, info); ok {
				each(val)
			}
		}
	}
	return hasDefaultCase
}

func makeSwitchDiagnostic(sw *ast.SwitchStmt, enumTypes []enumType, missing map[member]struct{}) analysis.Diagnostic {
	return analysis.Diagnostic{
		Pos: sw.Pos(),
		End: sw.End(),
		Message: fmt.Sprintf(
			"missing cases in switch of type %s: %s",
			diagnosticEnumTypes(enumTypes),
			diagnosticGroups(groupMissing(missing, enumTypes)),
		),
	}
}
