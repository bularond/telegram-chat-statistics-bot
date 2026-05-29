# Telegram Chat Statistics Bot

A Telegram bot written in Go that analyzes exported chat histories and replies with statistics about the conversation.

## How it works

1. Export a chat history from the Telegram desktop app as a **machine-readable JSON** file.
2. Send the exported file to the bot.
3. The bot parses the file and replies with statistics for the chat.

## Statistics

For each uploaded chat the bot reports:

- Total number of messages, words and characters
- Top 10 most active users (with their message, word and character counts)
- Top 10 most active days

## Usage

- `/start` or `/help` — show instructions on how to export and send a chat file.
- Send a chat export JSON file (up to **20 MB**) to receive the statistics.

## Configuration

The bot reads its token from the `TG_BOT_TOKEN` environment variable:

```bash
export TG_BOT_TOKEN="<your-telegram-bot-token>"
```

## Running

```bash
go run ./src
```

Or build a binary:

```bash
go build -o telegram-stats-bot ./src
./telegram-stats-bot
```

## Requirements

- Go 1.17+
- [go-telegram-bot-api/telegram-bot-api v5](https://github.com/go-telegram-bot-api/telegram-bot-api)

## Project structure

```
src/
  main.go              # Entry point: starts the worker and the bot
  bot/                 # Telegram bot: handles updates, commands and file uploads
  analytic/            # JSON parsing and chat statistics computation
```
