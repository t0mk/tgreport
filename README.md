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

You can download the latest release binary from the [GitHub Releases page](https://github.com/t0mk/tgreport/releases/latest).

Example (Linux):

```
wget https://github.com/t0mk/tgreport/releases/latest/download/tgreport-linux-amd64
chmod +x tgreport-linux-amd64
./tgreport-linux-amd64 test.yaml
```

Example (macOS):

```
curl -LO https://github.com/t0mk/tgreport/releases/latest/download/tgreport-darwin-amd64
chmod +x tgreport-darwin-amd64
./tgreport-darwin-amd64 test.yaml
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

## Running from Cron

You can automate tgreport using cron and provide the necessary environment variables for Telegram reporting.

### Example Cron Configuration

Edit your crontab with:
```
crontab -e
```

Add a line like the following (replace paths and values as needed):

```
0 8 * * * TG_TOKEN=your_telegram_bot_token TG_CHAT=your_telegram_chat_id /full/path/to/tgreport -t /full/path/to/test.yaml >> /tmp/tgreport.log 2>&1
```

- `0 8 * * *` runs the job every day at 8:00 AM.
- Replace `/full/path/to/tgreport` and `/full/path/to/test.yaml` with the actual paths on your system.
- Replace `your_telegram_bot_token` and `your_telegram_chat_id` with your actual Telegram bot token and chat ID.
- Output is appended to `/tmp/tgreport.log` for debugging.

**Note:** Cron runs with a minimal environment, so always use full paths and set required environment variables inline or at the top of your crontab.

## License
MIT 