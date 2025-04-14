package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/slack-go/slack"
)

var (
    // Slack Channel ID <=> Discord Channel ID (last part of the copy link URL)
    channelMap = map[string]string{
        "C053DPNGT": "525598064984195074", // support <=> support
        "C053DPNH5": "1361405944872829230", // random <=> slackbot test
    }
)

func main() {
	slackToken := os.Getenv("SLACK_BOT_TOKEN")
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")

	if slackToken == "" || discordToken == "" {
		log.Fatal("SLACK_BOT_TOKEN or DISCORD_BOT_TOKEN is not set")
	}

	// Init Slack client
	slackApi := slack.New(slackToken)
	rtm := slackApi.NewRTM()
	go rtm.ManageConnection()

	// Init Discord client
	discordApi, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}
	err = discordApi.Open()
	if err != nil {
		log.Fatalf("Error opening Discord connection: %v", err)
	}
	defer discordApi.Close()

	// Slack → Discord
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				discordChanID, ok := channelMap[ev.Channel]
				if !ok || ev.SubType != "" {
					continue
				}

				user, err := slackApi.GetUserInfo(ev.User)
				if err != nil {
					log.Printf("Slack: Failed to get user info: %v", err)
					continue
				}

                log.Printf("Slack: %s: %s", user.Name, ev.Text)

				text := fmt.Sprintf("**%s**: %s", user.Name, ev.Text)
				_, err = discordApi.ChannelMessageSend(discordChanID, text)
				if err != nil {
					log.Printf("Error sending message to Discord: %v", err)
				}
			}
		}
	}()

	// Discord → Slack
	discordApi.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return // Ignore bot messages to prevent loops
		}

		for slackChanID, discordChanID := range channelMap {
			if m.ChannelID == discordChanID {
				text := fmt.Sprintf("%s: %s", m.Author.Username, m.Content)

                log.Printf("Discord: %s", text)

				_, _, err := slackApi.PostMessage(slackChanID, slack.MsgOptionText(text, false))
				if err != nil {
					log.Printf("Error sending message to Slack: %v", err)
				}
				break
			}
		}
	})

	// Keep this program running until interrupted
	select {}
}
