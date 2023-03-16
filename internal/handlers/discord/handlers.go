package discordhandlers

import (
	"bytes"
	"disgourd/internal/audio_manager"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
)

var AudioUpload = discordgo.ApplicationCommand{
	Name:        "audio-upload",
	Description: "Upload audio to Fungible Duck",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user-to-map",
			Description: "User to associate with audio file",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionAttachment,
			Name:        "audio-file",
			Description: "Audio file to associate with user",
			Required:    true,
		},
	},
}

func GetInteractionCreateCommandTree(commands []*discordgo.ApplicationCommand) map[string]func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	icMap := make(map[string]func(s *discordgo.Session, ic *discordgo.InteractionCreate))
	for _, command := range commands {
		icMap[command.Name] = respond
	}
	return icMap
}

func AddSlashCommandHandlers(s *discordgo.Session, commandHandlers map[string]func(s *discordgo.Session, ic *discordgo.InteractionCreate)) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func VoiceStateUpdateHandler(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	if m.BeforeUpdate == nil && m.ChannelID != "" {
		log.Printf("%s has joined the voice channel: %s", m.UserID, m.ChannelID)
	} else if m.BeforeUpdate != nil && m.ChannelID == "" {
		log.Printf("%s has left the voice channel: %s", m.UserID, m.ChannelID)
	} else if m.BeforeUpdate != nil && m.ChannelID != m.BeforeUpdate.ChannelID {
		log.Printf("%s has switched from %s to %s", m.UserID, m.BeforeUpdate.ChannelID, m.ChannelID)
	}
}

func RegisterCommands(s *discordgo.Session, commands []*discordgo.ApplicationCommand) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func RemoveCommands(s *discordgo.Session, commands []*discordgo.ApplicationCommand) {
	for _, v := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//userID := i.ApplicationCommandData().Options[0].Value.(string)
	audioFileKey := i.ApplicationCommandData().Options[1].Value.(string)
	// Check if the session variable is nil
	//attachmentURL := i.ApplicationCommandData().Resolved.Attachments["0"].URL
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[audioFileKey].URL

	resp, err := http.Get(attachmentURL)
	if err != nil {
		log.Fatalf("Failed to download attachment: %v", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		log.Fatalf("Failed to read attachment: %v", err)
	}

	client := audiomanager.InitMinioConnection()
	audiomanager.CreateBucket(client, "test-bucket", "us-west-1")
	audiomanager.UploadFile(client, "test-bucket", "new-object", &buf)
	//err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//	Type: discordgo.InteractionResponseChannelMessageWithSource,
	//	Data: &discordgo.InteractionResponseData{
	//		Content: "Audio uploaded successfully!",
	//	},
	//})
	//if err != nil {
	//	log.Fatalf("Failed to respond to Discord interaction: %v", err)
	//}
}
