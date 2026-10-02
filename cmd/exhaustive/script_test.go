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

var exhaustiveCmd = flag.String("cmd", "exhaustive-scripttest", "path to executable to use in script tests")

func TestScript(t *testing.T) {
	ctx := t.Context()

	exhaustiveAbs, err := filepath.Abs(*exhaustiveCmd)
	if err != nil {
		t.Fatal(err)
	}
	cmds := scripttest.DefaultCmds()
	cmds["exhaustive"] = script.Program(exhaustiveAbs, func(cmd *exec.Cmd) error { return cmd.Process.Signal(os.Interrupt) }, 100*time.Millisecond)

	env := os.Environ()
	engine := &script.Engine{
		Conds: scripttest.DefaultConds(),
		Cmds:  cmds,
		Quiet: !testing.Verbose(),
	}
	scripttest.Test(t, ctx, engine, env, filepath.Join("testdata", "script", "*.txt"))
}
