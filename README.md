# tgreport

A tool to run command checks from a YAML config and optionally report failures to Telegram.

## Usage

```
tgreport [-t] -c <config.yaml>
```

- `-t`: Send failed check report to Telegram (requires `TG_TOKEN` and `TG_CHAT` environment variables)
- `-c <config.yaml>`: Path to the YAML config file (required)

## Example

```
tgreport -c test.yaml
```

or with Telegram reporting:

```
tgreport -t -c test.yaml
```

## Download

Go to the [GitHub Releases page](https://github.com/t0mk/tgreport/releases/latest) for all available binaries.

### Example download commands

#### Linux x86_64
```sh
wget https://github.com/t0mk/tgreport/releases/latest/download/tgreport-linux-amd64 -O tgreport
chmod +x tgreport
./tgreport -c test.yaml
```

#### macOS (Intel)
```sh
curl -Lo tgreport https://github.com/t0mk/tgreport/releases/latest/download/tgreport-darwin-amd64
chmod +x tgreport
./tgreport -c test.yaml
```

For other platforms, see the [Releases page](https://github.com/t0mk/tgreport/releases/latest).

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
0 8 * * * TG_TOKEN=your_telegram_bot_token TG_CHAT=your_telegram_chat_id /full/path/to/tgreport -t -c /full/path/to/test.yaml >> /tmp/tgreport.log 2>&1
```

- `0 8 * * *` runs the job every day at 8:00 AM.
- Replace `/full/path/to/tgreport` and `/full/path/to/test.yaml` with the actual paths on your system.
- Replace `your_telegram_bot_token` and `your_telegram_chat_id` with your actual Telegram bot token and chat ID.
- Output is appended to `/tmp/tgreport.log` for debugging.

**Note:** Cron runs with a minimal environment, so always use full paths and set required environment variables inline or at the top of your crontab.

## License
MIT 