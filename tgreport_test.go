package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestReportFailedCheck(t *testing.T) {
	check := Check{
		Command:        "echo hello",
		ExpectedOutput: "hello\n",
	}

	stdout, stderr, err := runCommand(check.Command)

	var buf bytes.Buffer
	if err != nil || stderr != "" || stdout != check.ExpectedOutput {
		reportFailedCheck(&buf, check.Command, stdout, stderr, check.ExpectedOutput)
		if buf.Len() == 0 {
			t.Error("Expected reportFailedCheck to produce output, but got none")
		}
	} else {
		reportFailedCheck(&buf, check.Command, stdout+"extra", stderr, check.ExpectedOutput)
		if buf.Len() == 0 {
			t.Error("Expected reportFailedCheck to produce output for failing check, but got none")
		}
	}
}

func TestLoadConfig(t *testing.T) {
	cfg, err := loadConfig("test.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expected := Config{
		Checks: []Check{
			{Command: "echo 'expected file contents, everything ok'", ExpectedOutput: "expected file contents, everything ok\n"},
			{Command: "cat /nonexistentfile", ExpectedOutput: "this will not match"},
			{Command: "printf 'line1\\nline2\\nline3\\n'\n", ExpectedOutput: "line1\nline2\nline3 "},
		},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Loaded config does not match expected.\nGot: %#v\nExpected: %#v", cfg, expected)
	}
}

func TestColorCodeRemoval(t *testing.T) {
	// Simulate output with ANSI color codes
	colorOutput := "\x1b[31mhello world\x1b[0m"
	expected := "hello world"
	clean := ansiRegexp.ReplaceAllString(colorOutput, "")
	if clean != expected {
		t.Errorf("Color codes not removed properly. Got: %q, Expected: %q", clean, expected)
	}
}

func TestWhitespaceTrimming(t *testing.T) {
	// Simulate output and expected with extra whitespace
	output := "  hello world  \n"
	expected := "hello world"
	trimmedOutput := strings.TrimSpace(output)
	trimmedExpected := strings.TrimSpace(expected)
	if trimmedOutput != trimmedExpected {
		t.Errorf("Whitespace not trimmed properly. Got: %q, Expected: %q", trimmedOutput, trimmedExpected)
	}
}

func TestSendTelegramMessage(t *testing.T) {
	// Start a local HTTP server to mock Telegram API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botTOKEN/sendMessage" {
			t.Errorf("Unexpected URL path: %s", r.URL.Path)
		}
		err := r.ParseForm()
		if err != nil {
			t.Errorf("Failed to parse form: %v", err)
		}
		if r.FormValue("chat_id") != "CHAT" {
			t.Errorf("chat_id: got %q, want %q", r.FormValue("chat_id"), "CHAT")
		}
		if r.FormValue("text") != "test message" {
			t.Errorf("text: got %q, want %q", r.FormValue("text"), "test message")
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	apiBase := ts.URL + "/bot%s/sendMessage"
	if err := sendTelegramMessage("TOKEN", "CHAT", "test message", apiBase); err != nil {
		t.Errorf("sendTelegramMessage failed: %v", err)
	}
}
