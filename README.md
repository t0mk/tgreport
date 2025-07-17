# tgreport

`tgreport` is a Go program for running command checks defined in a YAML file. It compares the output of each command to an expected value, and reports failures either to the console or via Telegram.

## Features
- Load checks from a YAML config file
- Run shell commands and compare their output to expected results
- Ignore ANSI color codes and whitespace differences in output
- Report failures to stdout or send them to a Telegram chat
- Includes tests for config loading, color code/whitespace handling, and Telegram integration

## Usage

### Basic
```
go build -o tgreport tgreport.go
./tgreport test.yaml
```

### With Telegram Reporting
Set the following environment variables:
- `TG_TOKEN`: Your Telegram bot token
- `TG_CHAT`: Your Telegram chat ID

Then run:
```
./tgreport -t test.yaml
```

## Download

You can download the latest release binary from the [GitHub Releases page](https://github.com/t0mk/tgreport/releases/tag/v0.1.0).

Example:

```
wget https://github.com/t0mk/tgreport/releases/download/v0.1.0/tgreport
chmod +x tgreport
./tgreport test.yaml
```

## YAML Config Example
```yaml
checks:
  - command: "echo 'expected file contents, everything ok'"
    expectedoutput: "expected file contents, everything ok\n"
  - command: "cat /nonexistentfile"
    expectedoutput: "this will not match"
  - command: |
      printf 'line1\nline2\nline3\n'
    expectedoutput: |
      line1
      line2
      line3
```
- Use `|` for multiline expected output.
- Whitespace and color codes are ignored in the comparison.

## Testing
Run all tests with:
```
go test -v
```

## License
MIT 