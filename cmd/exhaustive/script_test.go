package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

// Run one of these commands in the directory of this package to
// create the executable at the default path expected by the
// script tests.
//
//	go build -o exhaustive-scripttest
//	go build -o exhaustive-scripttest.exe  (Windows)
var cmdPath = flag.String("cmd", "exhaustive-scripttest", "path to executable to use in script tests")

func TestScript(t *testing.T) {
	ctx := t.Context()

	cmdPathAbs, err := filepath.Abs(*cmdPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(cmdPathAbs); err != nil {
		t.Fatal(err)
	}

	cmds := scripttest.DefaultCmds()
	cmds["exhaustive"] = script.Program(cmdPathAbs, func(cmd *exec.Cmd) error { return cmd.Process.Signal(os.Interrupt) }, 100*time.Millisecond)

	env := os.Environ()
	engine := &script.Engine{
		Conds: scripttest.DefaultConds(),
		Cmds:  cmds,
		Quiet: !testing.Verbose(),
	}
	scripttest.Test(t, ctx, engine, env, filepath.Join("testdata", "script", "*.txt"))
}
