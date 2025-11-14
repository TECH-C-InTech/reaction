package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"reaction/internal/entities"
	"reaction/internal/gateways"
	"reaction/internal/handlers"
	"reaction/internal/usecases"

	"github.com/bwmarrin/discordgo"
)

func main() {
	cfg, err := entities.LoadConfig()
	if err != nil {
		log.Fatal("設定の読み込みに失敗:", err)
	}

	// Discord Bot 接続
	dg, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		log.Fatal("Discordセッションの作成に失敗:", err)
	}

	// Gateway を作成（外部APIとの橋渡し）
	discordGateway := gateways.NewDiscordGateway(dg)

	// UseCase を作成
	transferUseCase := usecases.NewTransferMessageUseCase(cfg)

	// Handler を作成し、UseCase と Gateway を注入
	discordHandler := handlers.NewDiscordHandler(transferUseCase, discordGateway, cfg)

	dg.AddHandler(discordHandler.HandleReactionAdd)
	dg.AddHandler(discordHandler.HandleReactionRemove)

	// メッセージリアクションのインテントを有効にする
	dg.Identify.Intents = discordgo.IntentGuildMessageReactions

	err = dg.Open()
	if err != nil {
		log.Fatal("Discordへの接続に失敗:", err)
	}

	fmt.Println("Reaction Bot が起動しました!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// シャットダウン処理
	fmt.Println("Reaction Bot を終了しています...")
	err = dg.Close()
	if err != nil {
		log.Printf("Discord接続のクローズに失敗: %v", err)
	}
}
