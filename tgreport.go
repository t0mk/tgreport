package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// Check represents a single command check from the YAML config
type Check struct {
	Command        string `yaml:"command"`
	ExpectedOutput string `yaml:"expectedoutput"`
}

type Config struct {
	Checks []Check `yaml:"checks"`
}

func sendTelegramMessage(token, chat, message string, apiBase ...string) error {
	base := "https://api.telegram.org/bot%s/sendMessage"
	if len(apiBase) > 0 && apiBase[0] != "" {
		base = apiBase[0]
	}
	apiURL := fmt.Sprintf(base, token)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {chat},
		"text":    {message},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}
	return nil
}

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func reportFailedCheck(w io.Writer, cmd string, stdout, stderr, expected string) {
	if stderr != "" {
		fmt.Fprintf(w, "[FAIL] Command: %s\nStderr: %s\n", cmd, stderr)
	} else {
		fmt.Fprintf(w, "[FAIL] Command: %s\nGot: %s\nExpected: %s\n", cmd, stdout, expected)
	}
}

func runCommand(command string) (stdout string, stderr string, err error) {
	cmd := exec.Command("sh", "-c", command)
	stdoutBytes, err := cmd.Output()
	stderr = ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	stdout = string(stdoutBytes)
	return
}

func loadConfig(filename string) (Config, error) {
	var cfg Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func main() {
	useTelegram := flag.Bool("t", false, "Send failed check report to Telegram")
	configPath := flag.String("c", "", "Path to config YAML file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Usage: tgreport [-t] -c <config.yaml>")
		os.Exit(1)
	}
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	var tgToken, tgChat string
	if *useTelegram {
		tgToken = os.Getenv("TG_TOKEN")
		tgChat = os.Getenv("TG_CHAT")
		if tgToken == "" || tgChat == "" {
			fmt.Println("TG_TOKEN and TG_CHAT environment variables must be set for Telegram reporting.")
			os.Exit(1)
		}
	}

	for _, check := range cfg.Checks {
		stdout, stderr, err := runCommand(check.Command)
		cleanStdout := ansiRegexp.ReplaceAllString(stdout, "")
		trimmedStdout := strings.TrimSpace(cleanStdout)
		trimmedExpected := strings.TrimSpace(check.ExpectedOutput)
		if err != nil || stderr != "" || trimmedStdout != trimmedExpected {
			var report string
			if stderr != "" {
				report = fmt.Sprintf("[FAIL] Command: %s\nStderr: %s\n", check.Command, stderr)
			} else {
				report = fmt.Sprintf("[FAIL] Command: %s\nGot:\n%s\n---------\nExpected:\n%s\n", check.Command, trimmedStdout, trimmedExpected)
			}
			if *useTelegram {
				err := sendTelegramMessage(tgToken, tgChat, report)
				if err != nil {
					fmt.Printf("Failed to send Telegram message: %v\n", err)
				}
			} else {
				fmt.Print(report)
			}
		}
	}
}
