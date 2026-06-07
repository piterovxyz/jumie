package main

import (
	"bytes"
	"io"
	"jumie/internal/ai"
	"os"
	"testing"
	"time"
)

func TestTypewriterPrint(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	typewriterPrint("test", 1*time.Millisecond)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	if buf.String() != "test\n" {
		t.Errorf("expected test\\n, got %q", buf.String())
	}
}

func TestDoEmptyPlan(t *testing.T) {
	plan := &ai.Plan{}
	if do(plan) {
		t.Errorf("expected false for empty plan")
	}
}

func TestDoYes(t *testing.T) {
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write([]byte("y\n"))
	w.Close()

	plan := &ai.Plan{Steps: []ai.Step{{Command: "echo", Description: "test"}}}
	res := do(plan)
	os.Stdin = oldStdin

	if !res {
		t.Errorf("expected true")
	}
}

func TestSpinners(t *testing.T) {
	stop, update := startSpinner()
	update("test tip")
	time.Sleep(50 * time.Millisecond)
	stop()

	stopPrompt := startPromptSpinner()
	time.Sleep(50 * time.Millisecond)
	stopPrompt()

	stopLogin := startLoginPrompt()
	time.Sleep(50 * time.Millisecond)
	stopLogin()
}
