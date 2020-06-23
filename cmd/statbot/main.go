package main

import (
	"context"
	"flag"
	"log"

	"github.com/diamondburned/arikawa/state"
	"github.com/erebid/statbot"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	var (
		token       string
		databaseURL string
	)
	flag.StringVar(&token, "t", "", "discord bot token")
	flag.StringVar(&databaseURL, "d", "", "postgres database url")
	flag.Parse()

	state, err := state.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}

	dbpool, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalln("unable to connect to database:", err)
	}
	defer dbpool.Close()

	client, err := statbot.NewClient(state, dbpool)
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}
	if err := client.Open(); err != nil {
		log.Fatalln("Failed to connect:", err)
	}
	log.Println("connected")

	defer client.Close()

	select {}
}
