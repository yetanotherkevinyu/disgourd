package main

import (
	"disgourd/internal/config"
	"disgourd/internal/handlers/discord"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

var session *discordgo.Session

func init() {
	var err error
	conf := config.LoadConfig()
	session, err = discordgo.New("Bot " + conf.DiscordClientSecret)
	if err != nil {
		log.Fatalf("Invalid Discord session: %v", err)
	}
}

func waitForInterrupt() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

func main() {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is ready!") })
	commands := []*discordgo.ApplicationCommand{&discordhandlers.AudioUpload}
	ic_map := discordhandlers.GetInteractionCreateCommandTree(commands)
	session.AddHandler(discordhandlers.VoiceStateUpdateHandler)
	discordhandlers.AddSlashCommandHandlers(session, ic_map)

	err := session.Open()
	if err != nil {
		log.Fatalf("Cannot open the Discord session: %v", err)
	}
	defer session.Close()

	discordhandlers.RegisterCommands(session, commands)

	waitForInterrupt()
}
