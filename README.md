# Slack Discord Bridge for CakePHP Community

## Configuration

The `main.go` contains all the necessary configuration options. You can set the following environment variables:
- `SLACK_BOT_TOKEN`: Slack API token for the bot.
- `DISCORD_BOT_TOKEN`: Discord API token for the bot.

## How to build

Make sure to have Go version 1.24.2 or later installed.

```bash
go mod downloads
go build -o slack-discord-bridge
```

## How to run

Copy `.env.example` to `.env` file, adjust accordingly and execute `source .env` to load the environment 

or

make sure your environment variables are set correctly in the first place

```bash
./slack-discord-bridge
```
