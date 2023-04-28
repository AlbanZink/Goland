package bot

import (
	"fmt"                       //to print errors
	"golang-discord-bot/config" //importing our config package which we have created above
	"time"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var (
	BotId string
	goBot *discordgo.Session
)

func Start() {
	//creating new bot session
	goBot, err := discordgo.New("Bot " + config.Token)

	//Handling error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	//Handlinf error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Storing our id from u to BotId .
	BotId = u.ID

	// The first value is the created commands. You can save them in a variable, in case you want to clean them up at some point.
	_, err = goBot.ApplicationCommandBulkOverwrite(config.AppID, config.GuildID, []*discordgo.ApplicationCommand{
		{
			Name:        "time",
			Description: "Répond à l'intéraction par le temps actuel en utilisant la balise de temps.",
		},
	})
	if err != nil {
		fmt.Println("Error registering slash commands:", err)
		return
	}

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(messageHandler)
	goBot.AddHandler(timeCommandhandler)

	err = goBot.Open()
	//Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//If everything works fine we will be printing this.
}

// Definition of messageHandler function it takes two arguments first one is discordgo.Session which is s , second one is discordgo.MessageCreate which is m.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotId {
		return
	}
	//If we message ping to our bot in our discord it will return us pong .
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func timeCommandhandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "time":
		// Get the current timestamp in seconds
		currentTimestamp := time.Now().Unix()
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Replace EPOCH with the actual timestamp
					Content: fmt.Sprintf("<t:%d>", currentTimestamp),
				},
			},
		)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
