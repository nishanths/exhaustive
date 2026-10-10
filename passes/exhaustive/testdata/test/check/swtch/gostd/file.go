package gostd

import (
	"html/template"
)

func f1() {
	e := new(template.Error)

	switch e.ErrorCode { // want "^missing cases in expression switch: template.ErrBadHTML, template.ErrNoSuchTemplate, template.ErrPredefinedEscaper$"
	case template.OK:
	case template.ErrAmbigContext:
	case template.ErrBranchEnd:
	case template.ErrEndContext:
	case template.ErrOutputContext:
	case template.ErrPartialCharset:
	case template.ErrPartialEscape:
	case template.ErrRangeLoopReentry:
	case template.ErrSlashAmbig:
	case template.ErrJSTemplate:
	}
}
