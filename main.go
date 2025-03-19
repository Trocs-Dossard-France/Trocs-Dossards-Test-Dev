package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Charger les variables d'environnement
func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("❌ Erreur lors du chargement du fichier .env")
	}
}

// Token du bot
var Token = os.Getenv("DISCORD_TOKEN")

// ID du rôle à attribuer automatiquement
var roleID = os.Getenv("DISCORD_ROLE_ID")

// ID du salon où envoyer le message de bienvenue
var welcomeChannelID = os.Getenv("DISCORD_WELCOME_CHANNEL_ID")

func main() {
	if Token == "" || roleID == "" || welcomeChannelID == "" {
		fmt.Println("❌ Erreur : Une ou plusieurs variables d'environnement sont manquantes.")
		return
	}

	// Initialisation du bot Discord
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("❌ Erreur lors de l'initialisation du bot :", err)
		return
	}

	// Ajouter le handler pour l'accueil des membres
	dg.AddHandler(memberJoin)

	// Connexion au serveur
	err = dg.Open()
	if err != nil {
		fmt.Println("❌ Impossible de se connecter à Discord :", err)
		return
	}

	fmt.Println("✅ Bot connecté avec succès !")

	// Gestion de la fermeture proprement
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("🛑 Déconnexion du bot...")
	dg.Close()
}

// Fonction exécutée quand un membre rejoint le serveur
func memberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Envoyer un message privé de bienvenue
	dmChannel, err := s.UserChannelCreate(m.User.ID)
	if err == nil {
		s.ChannelMessageSend(dmChannel.ID, fmt.Sprintf("👋 Bienvenue sur le serveur, %s !\n\nLis les règles et amuse-toi bien !", m.User.Username))
	}

	// Envoyer un message de bienvenue dans le salon public
	s.ChannelMessageSend(welcomeChannelID, fmt.Sprintf("👋 Bienvenue à %s sur le serveur ! 🎉", m.User.Mention()))

	// Attribuer un rôle automatiquement
	err = s.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID)
	if err != nil {
		fmt.Println("❌ Erreur lors de l'attribution du rôle :", err)
	}
}
