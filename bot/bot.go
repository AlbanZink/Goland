package bot

import (
	"encoding/json"
	"fmt"                       //to print errors
	"golang-discord-bot/config" //importing our config package which we have created above
	"os"
	"time"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var (
	BotId string
	goBot *discordgo.Session
)

type User struct {
	Name string `json:"name"`
}

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
		{
			Name:        "writemyname",
			Description: "Ecris ton pseudo discord dans un fichier data.json",
		},
		{
			Name:        "saymyname",
			Description: "Lire le data.json et lister le nom",
		},
		{
			Name:        "hello",
			Description: "Creates a message with a button",
		},
	})
	if err != nil {
		fmt.Println("Error registering slash commands:", err)
		return
	}

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(messageHandler)
	goBot.AddHandler(timeCommandhandler)
	goBot.AddHandler(writeMyNameCommandHandler)
	goBot.AddHandler(sayMyNameCommandHandler)
	goBot.AddHandler(helloCommandHandler)
	goBot.AddHandler(buttonClickHandler)

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
	// Only process ApplicationCommand interactions
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()
	switch data.Name {
	case "time":
		// Get the current timestamp in seconds
		currentTimestamp := time.Now().Unix()
		nameUser := i.Member.User.Username
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// Replace EPOCH with the actual timestamp
				Content: fmt.Sprintf("%s, the current time is: <t:%d>", nameUser, currentTimestamp),
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func writeMyNameCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only process ApplicationCommand interactions
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	switch data.Name {
	case "writemyname":
		userName := i.Member.User.Username
		user := User{Name: userName}
		users := []User{user}
		userJSON, err := json.Marshal(users)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		file, err := os.Create("userName.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		_, err = file.Write(userJSON)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Send a response back to the user
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "User data written to userName.json file",
			},
		})
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}

func sayMyNameCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "saymyname":
			file, err := os.Open("userName.json")
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()
			// Decode the JSON data
			var users []User
			err = json.NewDecoder(file).Decode(&users)
			if err != nil {
				fmt.Println("Error decoding JSON:", err)
				return
			}
			// Print the users
			for _, user := range users {
				userName := user.Name
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: userName,
					},
				})
				if err != nil {
					fmt.Println("Error sending response:", err)
					return
				}
			}
		}
	}
}

func helloCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "hello":
			// Create a new button component
			button := &discordgo.Button{
				Label:    "Click me!",
				Style:    discordgo.PrimaryButton,
				CustomID: "hello_button",
			}
			buttonRow := &discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{button},
			}
			// Create the interaction response with the button component
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please click the button!",
					Components: []discordgo.MessageComponent{
						buttonRow,
					},
				},
			}
			// Respond to the interaction with the message containing the button
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				fmt.Println("Error responding to interaction:", err)
				return
			}
		}
	}
}

func buttonClickHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only process MessageComponent interactions
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	// Check if the clicked button has the custom ID "hello_button"
	if i.MessageComponentData().CustomID == "hello_button" {
		// Respond to the button click with "Hello World!"
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello World!",
			},
		})
		if err != nil {
			fmt.Println("Error responding to button click:", err)
			return
		}
	}
}
