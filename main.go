package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
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

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

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
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	pandas := []string{"p9", "panda"}

	for _, panda := range pandas {
		if strings.Contains(strings.ToLower(m.Content), panda) {
			s.ChannelMessageSend(m.ChannelID, "-1 xp for "+m.Author.Username)
		}
	}

	if strings.HasPrefix(m.Content, ".") {

		mm := strings.Fields(m.Content)

		log.Println(mm)

		for _, m := range mm {

			switch true {
			case strings.ContainsAny(m, "."):
				log.Println("command:", m)
			case strings.ContainsAny(m, "@"):
				log.Println("user:", m)
			case len(m) != 0:
				log.Println("reason:", mm[2:len(mm)])
			default:
				log.Println(".warn [@username] reason")
			}
		}
	}

	/*

		javascript version:

		if (message.author.id === discord.user.id || !message.member) return false;
		if (message.content && message.content.startsWith(".")){
			var text = message.content;
			var command = text.substring(1,text.indexOf(" "));
			var args = text.substring(text.indexOf(" ")+1);

			if(args != "" && commands.hasOwnProperty(command) && typeof commands[command] == "function"){
				if(args != "" && args != null)
					commands[command](message,args);
			}
			message.delete();
		}

	*/
}
