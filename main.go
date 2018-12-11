package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the onMessage func as a callback for MessageCreate events.
	dg.AddHandler(onMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	var commands = []string{"warn", "kick", "ban"}

	if strings.HasPrefix(m.Content, ".") {

		// Gives back command substrings
		mm := strings.Fields(m.Content)

		// verify user
		if strings.ContainsAny(mm[1], "@") {
			log.Println("user: ", mm[1])
		} else {
			// handle error
		}

		// verify reason
		if len(mm[2:]) != 0 {
			log.Println("reason: ", strings.Join(mm[2:], " "))
		} else {
			// handle error
		}

		// verify command
		for _, command := range commands {
			if strings.Contains(mm[0], command) {

				switch command {
				case commands[0]:

					author := m.Author.Username
					id := m.ChannelID
					color := 112244
					command := mm[0]
					user := mm[1]
					reason := strings.Join(mm[2:], " ")

					sendMessage(s, id, author, user)
					sendPrivateMessage() // good for now
					Log(s, m, color, command, user, reason)
				case commands[1]:
				case commands[2]:
				}
			}
		}
	}

	/*

		command := []string{".command [@username] reason", "Valid commands are:", "warn", "kick", "ban"}
		warning := strings.Join(command, "\n")

		s.ChannelMessageSend(m.ChannelID, warning)

	*/
}

// Log logs a new command
func Log(s *discordgo.Session, m *discordgo.MessageCreate, color int, command, user, reason string) {

	embed := &discordgo.MessageEmbed{
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Command",
				Value: command,
			},
			&discordgo.MessageEmbedField{
				Name:  "User",
				Value: user,
			},
			&discordgo.MessageEmbedField{
				Name:  "Reason",
				Value: reason,
			},
			&discordgo.MessageEmbedField{
				Name:  "Text Channel",
				Value: "<#" + m.ChannelID + ">",
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: m.Author.AvatarURL(""),
			Text:    m.Author.String(),
		},
	}

	st, _ := s.Channel(m.ChannelID)
	xc, _ := s.GuildChannels(st.GuildID)
	for _, v := range xc {
		if strings.Contains(v.Name, "moderation-log") == true {
			s.ChannelMessageSendEmbed(v.ID, embed)
		} else {
			// make channel
		}
	}
}

func sendMessage(s *discordgo.Session, id, author, user string) {
	s.ChannelMessageSend(id, author+" has issued a warning to user "+user)
}

func sendPrivateMessage() {}
