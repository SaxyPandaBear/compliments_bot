package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Auth struct {
	Token string
}

func init() {
	// read the token from the json blob "discord_token.json"
	file, err := os.Open("discord_token.json")
	if err != nil {
		fmt.Println("error reading file", err)
		os.Exit(1)
	}
	bytes := make([]byte, 1000 ) // TODO: not sure how many bytes I need to read in total yet
	count, err := file.Read(bytes) // read bytes of file into bytes array
	if err != nil {
		fmt.Println("error reading file", err)
		os.Exit(1)
	}
	var tokenJson Auth
	err = json.Unmarshal(bytes[:count], &tokenJson)
	if err != nil {
		fmt.Println("error decoding json blob", err)
		os.Exit(1)
	}
	defer file.Close()
	Token = tokenJson.Token

	// read compliments before we create our bot
	if success := readCompliments(); !success {
		os.Exit(1)
	}
}

var Token string // client token from Discord

var Compliments []string // array of compliment templates that the bot uses

// reads a file into the compliments array
// expects a compliments.txt file in the same directory as this code
// returns the state of success
func readCompliments() (bool) {
	file, err := os.Open("compliments.txt")
	if err != nil {
		fmt.Println("error opening compliments file", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Compliments = append(Compliments, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error reading compliments file", err)
		return false
	}

	// if our compliments array is empty, then we can't continue because there's no compliments to use
	return len(Compliments) > 0
}

// gets a random compliment from the compliments array
// expects compliments array to be non-empty
func getRandCompliment() (string) {
	index := rand.Intn(len(Compliments)) // get a random number between [0, len(Compliments))
	return Compliments[index]
}

func main() {
	// instantiate the discord bot
	dg, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// event handler for messages in channels that the bot is authorized to access
	dg.AddHandler(onMessage)
	// event handler for when new users join a server
	dg.AddHandler(onNewMember)

	// open a web socket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// see https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go example
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	close(sc)
	fmt.Println("\nBot is now closing.")
	dg.Close()
}

// function that handles a new message event
// we want to listen only on messages that start with !waifu
const prefix string = "!waifu"
func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages that are sent by this bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, prefix){
		// we need to check a subcommand
		// if no target is specified, then we cannot do anything. print a usage message
		// !waifu compliment me => bot picks a random compliment directed towards the message author
		// !waifu compliment @user => bot picks a random compliment directed towards user.
		// - assumes that @user is a valid user in the server
		mParts := strings.Split(m.Content, " ")
		if len(mParts) < 3 {
			// if not enough args, return
			s.ChannelMessageSend(m.ChannelID, botUsage())
			return
		}
		if mParts[1] != "compliment" {
			// if we invoke the bot without the compliment command, return
			s.ChannelMessageSend(m.ChannelID, botUsage())
			return
		}
		var user string // user that the bot will mention
		if mParts[2] == "me" {
			user = m.Author.Mention()
		} else {
			user = mParts[2] // we assume this is a valid string
		}
		compliment := getRandCompliment()
		var message string = fmt.Sprintf("Master %s ... %s", user, compliment)
		s.ChannelMessageSend(m.ChannelID, message)
	}
}

// function that handles a new member event
// When a member joins the server (GuildMemberAdd event), greet the user.
// TODO: check the order of channels when a new user joins
func onNewMember(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		fmt.Println("error getting server", err)
		return
	}
	if len(guild.Channels) < 1 {
		fmt.Println("no channels found")
		return
	}
	var channel *discordgo.Channel = nil // we want to get the first channel whose Type is 1, a text chat
	// we want the first text channel because Discord defines this as the "Default" channel now.
	channels := guild.Channels
	for i := 0; i < len(channels); i++ {
		if channels[i].Type == discordgo.ChannelTypeGuildText {
			channel = channels[i]
			break
		}
	}
	if channel == nil {
		fmt.Println("no text channel found")
		return
	}
	s.ChannelMessageSend(channel.ID, fmt.Sprintf("Hello, Master %s", m.User.Mention()))
}

// function that returns a string that explains the bot's compliment usage
func botUsage() (string) {
	return "Type !waifu compliment me to have me compliment you.\n" +
		"Type !waifu compliment @someone to have me compliment them."
}
