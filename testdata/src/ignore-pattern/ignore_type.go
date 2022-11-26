package ignorepattern

import (
	"reflect"
	"time"
)

type label string // want label:"^home,work,other$"

const (
	home  label = "home"
	work  label = "work"
	other label = "other"
)

func _i() {
	var v label
	switch v {
	case work:
	}

	_ = map[label]struct{}{
		home: {},
	}
}

func _j() {
	var v time.Duration
	switch v {
	case time.Nanosecond:
	case 5 * time.Second:
	}

	_ = map[time.Duration]struct{}{
		time.Hour: {},
	}
}

func _k() {
	var v reflect.Kind
	switch v {
	case reflect.Invalid:
	case reflect.Bool:
	case reflect.Uintptr:
	case reflect.Interface:
	default:
	}

	_ = map[reflect.Kind]struct{}{
		reflect.Uint:    {},
		reflect.Uint8:   {},
		reflect.Uint16:  {},
		reflect.Uint32:  {},
		reflect.Uint64:  {},
		reflect.Uintptr: {},
	}
}
